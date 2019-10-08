package logic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Cluas/gim/internal/logic/conf"
	"github.com/Cluas/gim/pkg/log"
)

const (
	split = "@"
)

func ParseNetwork(str string) (network, addr string, err error) {
	if idx := strings.Index(str, split); idx == -1 {
		err = fmt.Errorf("addr: \"%s\" error, must be network@tcp:port or network@unixsocket", str)
		return
	} else {
		network = str[:idx]
		addr = str[idx+1:]
		return
	}
}
func InitHTTP() (err error) {
	var network, addr string

	for i := 0; i < len(conf.Conf.HTTP.Address); i++ {

		httpServeMux := http.NewServeMux()
		httpServeMux.HandleFunc("/api/v1/push", Push)

		if network, addr, err = ParseNetwork(conf.Conf.HTTP.Address[i]); err != nil {
			log.Errorf("ParseNetwork() error(%v)", err)
			return
		}

		log.Infof("start http listen:\"%s\"", conf.Conf.HTTP.Address[i])

		go httpListen(httpServeMux, network, addr)
		select {}

	}
	return
}
func Push(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", 405)
	}

	var (
		auth         = r.URL.Query().Get("auth")
		acceptUserId = r.URL.Query().Get("user_id")
		err          error
		bodyBytes    []byte
		body         string
		formUserInfo *Router
		res          = map[string]interface{}{"code": 1, "msg": "success"}
		sendData     *Send
	)

	if formUserInfo, err = getRouter(auth); err != nil {
		log.Errorf("get router error : %s", err)
		return
	}

	log.Infof("push round userId %s", formUserInfo.UserID)

	if bodyBytes, err = ioutil.ReadAll(r.Body); err != nil {
		res["code"] = 599
		res["msg"] = "define.NETWORK_ERR_MSG"
		log.Errorf("get router error : %s", err)
		return
	}

	serverID := RedisCli.Get(getKey(acceptUserId)).Val()
	serverID = "1"
	sid, err := strconv.ParseInt(serverID, 10, 8)
	if err != nil {
		res["code"] = -11
		res["msg"] = "define.SEND_ERR_MSG"
		log.Errorf("router err %v", err)
		return
	}

	defer retPWrite(w, r, res, &body, time.Now())

	err = json.Unmarshal(bodyBytes, &sendData)
	if err != nil {
		log.Errorf("sendData err: %s", err)
	}
	sendData.FormUserName = formUserInfo.Username
	sendData.FormUserId = formUserInfo.UserID
	sendData.Op = int32(2)
	if bodyBytes, err = json.Marshal(sendData); err != nil {
		log.Errorf("redis Publish err: %s", err)
	}
	body = string(bodyBytes)

	if err := RedisPublishCh(int8(sid), acceptUserId, bodyBytes); err != nil {
		log.Errorf("redis Publish err: %s", err)

	}
	log.Info("send message")

}
func retPWrite(w http.ResponseWriter, r *http.Request, res map[string]interface{}, body *string, start time.Time) {
	data, err := json.Marshal(res)
	if err != nil {
		log.Errorf("json.Marshal(\"%v\") error(%v)", res, err)
		return
	}
	dataStr := string(data)
	log.Infof("dataStr %s", dataStr)

	w.Header().Set("Access-Control-Allow-Origin", "*")             // 允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") // header的类型
	w.Header().Set("content-type", "application/json")
	// 返回数据格式是json
	if _, err := w.Write([]byte(dataStr)); err != nil {
		log.Errorf("w.Write(\"%s\") error(%v)", dataStr, err)
	}

	log.Infof("req: \"%s\", post: \"%s\", res:\"%s\", ip:\"%s\", time:\"%fs\"", r.URL.String(), *body, dataStr, r.RemoteAddr, time.Now().Sub(start).Seconds())
}

func httpListen(mux *http.ServeMux, network, addr string) {

	httpServer := &http.Server{Handler: mux, ReadTimeout: conf.Conf.HTTP.ReadTimeout, WriteTimeout: conf.Conf.HTTP.WriteTimeout}
	httpServer.SetKeepAlivesEnabled(true)

	l, err := net.Listen(network, addr)
	if err != nil {
		log.Errorf("net.Listen(\"%s\", \"%s\") error(%v)", network, addr, err)
		panic(err)
	}
	if err := httpServer.Serve(l); err != nil {
		log.Errorf("server.Serve() error(%v)", err)
		panic(err)
	}
}

type Send struct {
	Code         int32  `json:"code"`
	Msg          string `json:"msg"`
	FormUserId   string `json:"fuid"`
	FormUserName string `json:"fname"`
	Op           int32  `json:"op"`
}
