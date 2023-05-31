package consul

import (
	"bytes"
	"cmgo/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	gServerIp  string
	gServerIdc string
)

type Consul struct {
	Token       string `toml:"token"`
	ServiceName string `toml:"service_name"`
	Port        string `toml:"port"`
}

const (
	DefaultIdc             = "bx"
	consulAgentAddr        = "http://127.0.0.1:8500"
	consulClientRegister   = consulAgentAddr + "/v1/agent/service/register?replace-existing-checks=true"
	consulClientDeRegister = consulAgentAddr + "/v1/agent/service/deregister/"
	idcUrl                 = ""
)

// ConsulRequest is JSON defination of request to consul.
type ConsulRequest struct {
	ID    string       `json:"ID"`
	Name  string       `json:"Name"`
	Addr  string       `json:"Address"`
	Port  int          `json:"Port"`
	Tags  []string     `json:"Tags"`
	Meta  MetaRequest  `json:"Meta"`
	Check CheckRequest `json:"Check"`
}

// MetaRequest is JSON defination of Meta of ConsulRequest.
type MetaRequest struct {
	Idc string `json:"_idc"`
}

// CheckRequest is JSON defination of Check of ConsulRequest.
type CheckRequest struct {
	DeregisterAfter string `json:"DeregisterCriticalServiceAfter"`
	Interval        string `json:"Interval"`
	Timeout         string `json:"Timeout"`
	TCP             string `json:"TCP"`
}

func RegisterConsul(consulCfg Consul) (err error) {
	gServerIp = utils.GetLocalIp()
	gServerIdc = GetServerIDC()
	portNumber, _ := strconv.Atoi(consulCfg.Port)
	cr := ConsulRequest{
		ID:   getServiceId(consulCfg.ServiceName, gServerIp),
		Name: consulCfg.ServiceName,
		Addr: gServerIp,
		Port: portNumber,
		Tags: []string{consulCfg.ServiceName, gServerIp},
		Meta: MetaRequest{
			Idc: gServerIdc,
		},
		Check: CheckRequest{
			DeregisterAfter: "5m",
			Interval:        "15s",
			Timeout:         "2s",
			TCP:             gServerIp + ":" + consulCfg.Port,
		},
	}
	payload, err := json.Marshal(cr)
	if err != nil {
		return
	}
	req, err := http.NewRequest("PUT", consulClientRegister, bytes.NewReader(payload))
	if err != nil {
		return
	}
	req.Header.Add("X-Consul-Token", consulCfg.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	return
}

func DeregisterConsul(consulCfg Consul) {
	req, _ := http.NewRequest("PUT", consulClientDeRegister+getServiceId(consulCfg.ServiceName, gServerIp), strings.NewReader(""))
	req.Header.Add("X-Consul-Token", consulCfg.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("DeregisterConsul err:%s", err.Error())
		return
	}
	defer res.Body.Close()
}

func getServiceId(consulServiceName, ip string) string {
	return fmt.Sprintf("%s:%s", consulServiceName, ip)
}

func GetServerIDC() (name string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("get idc fatal error:%v", err)
		}
	}()
	name = DefaultIdc
	type result struct {
		IDC       string `json:"idc"`
		IpAddress string `json:"ipaddress"`
	}
	ip := utils.GetLocalIp()
	if ip == "" {
		return
	}
	requestUrl := fmt.Sprintf(idcUrl, ip)
	c := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := c.Get(requestUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var res result
	err = json.Unmarshal(body, &res)
	if err == nil && res.IDC != "" {
		name = res.IDC
	}
	return
}
