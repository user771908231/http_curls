package main

import (
	"bytes"
	"crypto/des"
	"encoding/base64"
	"fmt"
)

var _curlPath string
type RequestParameter struct {
	Url 	string
	Header	string
	Method	string
	Charset	string
	Encod	string

}



func HttpCutl(Parameter RequestParameter) {

}

func DesEncrypt(src, key []byte) (string, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	//src = ZeroPadding(src, bs)
	src = PKCS5Padding(src, bs)
	if len(src)%bs != 0 {

		return "", fmt.Errorf("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}

	return Base64Encode_(out), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func Base64Encode_(str []byte) string {
	return base64.StdEncoding.EncodeToString(str)
}