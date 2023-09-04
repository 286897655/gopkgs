package httptool

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func SimplePost(url string,buffer []byte)([]byte,error){
	resp, err := http.Post(url, "application/json", bytes.NewReader(buffer))
	if err != nil{
		return nil,err
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func JsonPost(postUrl string, para interface{}) ([]byte, error) {
	bs, err := json.Marshal(para)
	if err != nil {
		return nil, err
	}
	requestDo, err := http.NewRequest("POST", postUrl, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	requestDo.Header.Set("Content-type", "application/json")
	// do it
	client := &http.Client{
		// 通过TransPort设置具体项超时
		// Transport: &http.Transport{
		// 	Dial: func(network, addr string) (net.Conn, error) {
		// 		conn, err := net.DialTimeout(network, addr, time.Second*3) // 连接3秒超时
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 		err = conn.SetDeadline(time.Now().Add(time.Second * 3)) //设置发送和接收数据超时时间
		// 		return conn, err
		// 	},
		// 	ResponseHeaderTimeout: time.Second * 3, //response header timeout 2 second
		// },
		// 设置简单版超时3秒
		Timeout: time.Second * 3,
	}
	resp, err := client.Do(requestDo)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func 