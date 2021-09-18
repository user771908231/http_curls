package main

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"testing"
)

type Msg struct {
	Code int    //
	Data []byte //
}

func TestHeader(t *testing.T) {
	/**
	/Users/js43/go/src/lottery_data_supplement/libs/httpcurl/http_curl -u https://api.zgjdgj.com/auto/history -m POST -h {"origin":"https://pk.happipk.com","referer":"https://pk.happipk.com"} -v {"fc_type":"yfliuhecai_sk","token":"2640787a28ba68cb4fda856b103e3090","page":"1","pagenum":"1440"} -e  -c utf-8 -proxy  -t 500
	 */
	var encod = false
	var outInfo bytes.Buffer
	a := struct{
		Origin	string `json:"origin"`
		Referer string `json:"referer"`
		UserAgent string `json:"user-agent"`
		Host    string `json:"host"`
	}{
		Origin:"https://pk.happipk.com/",
		Referer: "https://pk.happipk.com/",
		Host:"https://pk.happipk.com/",
		UserAgent:"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
	}
	v := struct {
		FcType 		string `json:"fc_type"`
		Token		string `json:"token"`
		Gameplay	string `json:"gameplay"`
		Pankou		string `json:"pankou"`
	}{
		FcType: "jsliuhecai_sk",
		Token:  "1d7cad1456e1f308d3444cac3f874036",
		Gameplay: "165",
		Pankou: "A",
	}
	data, err := json.Marshal(a)
	if err != nil {
		fmt.Println("json.marshal failed, err:", err)
		return
	}
	t.Log(string(data))
	data_v, err := json.Marshal(v)
	if err != nil {
		fmt.Println("json.marshal failed, err:", err)
		return
	}
	log.Println(string(data_v))

	cmd := exec.Command("/Users/js43/go/src/http_curl/http_curl",
		"-u","https://api.zgjdgj.com/auto/history",
		"-m","POST",
		"-v",string(data_v),
		"-h",string(data),
		"-e","")
	cmd.Stdout = &outInfo
	// 执行
	err = cmd.Run()
	if err != nil {
		t.Log("aaa",err)
	}
	res_string_data := string(outInfo.Bytes())
	if strings.HasPrefix(res_string_data, "DATA|") {
		res_string_data = strings.TrimLeft(res_string_data, "DATA|")
	} else {
		t.Log("err",res_string_data)
	}

	decoded, err := base64.StdEncoding.DecodeString(res_string_data)

	log.Println(string(decoded))
	//res_data,err := base64.StdEncoding.DecodeString(string(outInfo.Bytes()))
	if err != nil {
		t.Log("ddd",err)
	}
	if encod {
		in := bytes.NewReader(decoded)
		dic := []byte(`了道你我的去来在是他便都那也得又把却这人<br/>,.!”`)
		var out bytes.Buffer
		r, err := zlib.NewReaderDict(in, dic)
		_, err = io.Copy(&out, r)
		if err != nil {
			t.Log(err)
		}
		t.Log("www",string(out.Bytes()))
		return
	}


	buffer := bytes.NewBuffer(decoded)
	err = Unzip("/Users/js43/go/src/http_curl/zip", buffer)
	t.Log(err)
}

func TestF(t *testing.T) {
	var sss = "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"
	decoded, err := base64.StdEncoding.DecodeString(sss)
	if err != nil {
		t.Log(err)
	}
	t.Log(string(decoded))

}

func TestStringJsonToMap(t *testing.T) {
	s := "{\"fc_type\":\"ffc_sk\",\"token\":\"cd9c715db273c95b3b432ca1f9d79ade\"}"
	res,err  := StringJsonToMap(s)
	if err != nil {
		t.Log(err)
	}
	t.Log(res)
}

func TestOnec(t *testing.T) {
		/**
			  $data = Helper::encrypt(
		            "fc_type=" . 'ffc_sk' .      //彩种类型，必传
		            "/\\\\/periods=" . '' .           //默认为空,期数
		            "/\\\\/number=" . 10 .            //请求条数，必传
		            "/\\\\/agent_name=" . 'zzz' .          //必传,如:zzz,aok
		            "/\\\\/day=" . '2019-07-01/2019-07-10' .     //必传,需要查询多天的格式：2019-08-24/2019-08-26，一天格式:2019-08-26
		            "/\\\\/line_id=" . 'ccc'  .         //线路id，必传
		              "/\\\\/page=" . '1'           //页数
		            ,'APDVGLJK'
		        );
		 */

	params := make([]string, 0)
	params = append(params, "fc_type=lfc_sk")
	params = append(params, "periods=")
	params = append(params, "number=10")
	params = append(params, "agent_name=zzz")
	params = append(params, "day=2021-07-15")
	params = append(params, "line_id=ccc")
	params = append(params, "page=1")

	//var resp interface{}

	// mcrypt  加密
	//q_params, err := model.EncryptNew(strings.Join(params, "/\\\\/"), []byte(m.DesKey))
	//if err != nil {
	//	golog.Error(siteId, "", "VideoGameService", "List", "err:", err)
	//	return resp, err
	//}
	res,err := DesEncrypt([]byte(strings.Join(params, "/\\\\/")),[]byte("APDVGLJK"))
	t.Log(res,err)
	md5Key := "aspfchjklmwertiup"
	w := md5.New()
	io.WriteString(w, res+md5Key)   //将str写入到w中
	md5str2 := fmt.Sprintf("%x", w.Sum(nil))
	s := fmt.Sprintf("data=%s&key=%s",res,md5str2)
	encoded := base64.StdEncoding.EncodeToString([]byte(s))

	fmt.Println(encoded,"-->>>>")
	u := fmt.Sprintf("https://ourbackend.lltest.me/systemapi/lottery/autolist?paramsStr=%s",encoded)

	t.Log(u)
}

func TestTew(t *testing.T) {
	s := make([]int,2)
	s = append(s,1)
	s = append(s,2)
	t.Log(s)
}
