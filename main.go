package main

import (
	"archive/zip"
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/axgle/mahonia"
	"http_curl/httpcli"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)
/**
/Users/js43/go/src/lottery_data_supplement/libs/httpcurl/http_curl -u https://api.zgjdgj.com/auto/history -m POST -h {"origin":"https://pk.happipk.com","referer":"https://pk.happipk.com"} -v {"fc_type":"yfliuhecai_sk","token":"2640787a28ba68cb4fda856b103e3090","page":"1","pagenum":"1440"} -e  -c utf-8 -proxy  -t 500
*/
var cur_path *string = flag.String("u","https://api.zgjdgj.com/auto/history","url path")
var header	 *string = flag.String("h","{\"origin\":\"https://pk.happipk.com\",\"referer\":\"https://pk.happipk.com\"}","http header")
var method	 *string = flag.String("m","POST","request method")
var charset	 *string = flag.String("c","utf-8","charset")
var encod	 *string   = flag.String("e","","encod")
var proxy 	 *string = flag.String("proxy","","proxy")
var timeout  *int64  = flag.Int64("t", 500, "timeout ")
var ip       *string = flag.String("ip", "", "ip ")
var value	 *string = flag.String("v","{\"fc_type\":\"yfliuhecai_sk\",\"token\":\"8ce859229b42898bfb993a5450b11366\",\"page\":\"1\",\"pagenum\":\"1440\"}","http value")
const RETRY = 3

var (
	headers = http.Header{}
	values  = url.Values{}
)

func main()  {
	flag.Parse()
	u, err := url.Parse(*cur_path)
	if err != nil {
		fmt.Println("ERR|url.Pars" + err.Error())
		return
	}
	if *header != "" {
		h_map,err := StringJsonToMap(*header)
		if err != nil {
			fmt.Println("ERR| headers" + err.Error())
			return
		}
		for key, val := range h_map {
			headers.Set(key,val)
		}
	}

	if *value != "" {
		h_map,err := StringJsonToMap(*value)
		if err != nil {
			fmt.Println("ERR| values" + err.Error())
			return
		}
		for key, val := range h_map {
			values.Add(key,val)
		}
	}
	//log.Println(values)
	//fmt.Println(*header)
	//log.Println(u)
	//log.Println("cur_path",*cur_path)
	//log.Println("header",*header)
	//log.Println("method",*method)
	//log.Println("charset",*charset)
	//log.Println("encod",*encod)
	//log.Println(*encod)
	res,err := NewHtml(u.String(),*method,*charset)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("DATA|" + string(res))
}

func NewHtml(rawurl,method ,charset string) ([]byte, error) {
	var res []byte
	var body io.Reader
	var resp *http.Response
	var err error

	c := httpcli.NewClient(
		&httpcli.ClientConfig{
			Timeout:   time.Duration(*timeout) * time.Second,
			Dial:      time.Duration(*timeout) * time.Second,
			KeepAlive: time.Duration(*timeout) * time.Second,
			ProxyURL:  *proxy,
		})
	//重试日志
	ips := *ip


	for i := 0; i < RETRY; i++ {
		if strings.ToTitle(method) == "GET"{
			res, resp, err = c.Get(context.TODO(), ips, rawurl, values,headers)
		}else{
			res, resp, err = c.Post(context.TODO(), ips, rawurl, values,headers)
		}

		if err == nil {
			break
		}
		if resp == nil {
			break
		}
		if resp.StatusCode == 301 || resp.StatusCode == 302 {
			locat, err := resp.Location()
			if err != nil {
				break
			}
			ips = ""
			rawurl = locat.String()
			continue
		}
		if resp != nil && (resp.StatusCode == 301 || resp.StatusCode == 302) {
			break
		}
		//休眠10ms，防止采集速度过快被屏蔽   //
		time.Sleep(time.Duration(100) * time.Millisecond)
		if err != nil {
			return nil, fmt.Errorf("ERR|charset : %s is error", err)
		}
	}

	if *encod != "" {
		body = bytes.NewReader(res)
		// 编码格式转码
		enc := mahonia.NewDecoder(charset)
		if enc == nil {
			return nil, fmt.Errorf("ERR|charset : %s is error", charset)
		}
		body = enc.NewReader(body)
		res, err = ioutil.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("ERR|ReadAll  enc : %s is error", charset)
		}

		////压缩
		var in bytes.Buffer
		dic := []byte(`了道你我的去来在是他便都那也得又把却这人<br/>,.!”`)
		w, err := zlib.NewWriterLevelDict(&in, zlib.BestCompression, dic)
		if err != nil {
			return nil, fmt.Errorf("ReadAll  enc : %s is error", charset)
		}
		w.Write(res)
		w.Close()
		res = in.Bytes()
		//fmt.Println(in.String())
	}
	if charset == "img" {

		//log.Println(res)
		//content_type := resp.Header.Get("content-type")

		//encoded := base64.StdEncoding.EncodeToString(res)

		return []byte(res), nil
	}
	if charset == "zip" {
		return res, nil
	}
	log.Println(string(res))
	encoded := base64.StdEncoding.EncodeToString(res)

	return []byte(encoded), nil
}

func StringJsonToMap(str string) (map[string]string,error) {
	if len(str) <= 0 { return nil,fmt.Errorf("string not find")}
	s := make(map[string]string,0)
	err := json.Unmarshal([]byte(str),&s)
	if err != nil {
		return nil,err
	}
	return s, err
}

func Unzip(basePath string, r io.Reader) error {
	//log.Println(r)
	/* 创建属于解压的缓存目录 */
	var dir = path.Join("/Users/js43/go/src/http_curl/temp", "zip")
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	log.Println(dir)
	/* 创建解压缓存文件 */
	f, e := ioutil.TempFile(dir, "zip")
	if nil != e { return e }
	log.Println(f.Name())
	//defer func() {
	//	f.Close()
	//	os.RemoveAll(f.Name())
	//}()

	_, e = io.Copy(f, r)
	if nil != e { return e }

	return unzip(basePath, f)
}

func unzip(basePath string, f *os.File) error {
	var reader *zip.Reader
	var stat, _ = f.Stat()
	log.Println(f.Name())
	reader, e := zip.NewReader(f, stat.Size())
	if nil != e { return e }

	err := os.MkdirAll(basePath, 0777) // 确保解压目录存在
	if err != nil {
		return err
	}
	for _, info := range reader.File {
		var fp = toLinux(path.Join(basePath, info.Name))
		log.Println(fp)
		if info.FileInfo().IsDir() {
			if e := os.MkdirAll(fp, info.FileInfo().Mode()); nil != e { return e }
			continue
		}

		readcloser, e := info.Open()
		if nil != e { return e }

		b, e := ioutil.ReadAll(readcloser)
		if nil != e { return e }
		readcloser.Close()
		log.Println(string(b))
		if e := ioutil.WriteFile(fp, b, info.FileInfo().Mode()); nil != e { return e }
	}
	return nil
}


func toLinux(basePath string) string {
	return strings.ReplaceAll(basePath, "\\", "/")
}

func Zip(fp string, w io.ReadWriter) error {
	archive := zip.NewWriter(w)
	defer archive.Close()

	linuxFilePath := toLinux(fp)
	filepath.Walk(linuxFilePath, func(path string, info os.FileInfo, err error) error {

		var linuxPath = toLinux(path)
		if linuxPath == linuxFilePath { return nil }

		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(linuxPath, linuxFilePath + "/")

		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}
		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(linuxPath)
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	})

	return nil
}
