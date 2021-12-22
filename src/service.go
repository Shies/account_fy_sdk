package account_service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	showAcInfoUri		= "/open/ac/water/showAcInfo"
	acSetUri			= "/open/ac/water/acSet"
	elecFeeSumUri		= "/open/ac/water/elecFeeSum"
)

var AccountService = NewService()
type ValueSet map[string]interface{}

// Account Service
type accountService struct {
	context.Context
	queryStatus map[string]*AcStatus
	wg			*sync.WaitGroup
}

func NewService() *accountService {
	s := &accountService{
		Context:     context.TODO(),
		queryStatus: make(map[string]*AcStatus),
		wg:          new(sync.WaitGroup),
	}

	return s
}

func (a *accountService) GetQueryInfo(scheme *Scheme) (res []*AcStatus, err error) {
	params := url.Values{}
	params.Set("appKey", scheme.AppKey)
	params.Set("timestamp", strconv.FormatInt(MillUnix(), 10))
	params.Set("sign", strings.ToUpper(getSign(scheme, nil)))

	url := scheme.RequestUrl + showAcInfoUri + "?" + params.Encode()
	response, err := HttpGet(url, nil)
	//log.Printf("%s,%s", url, response)
	if err != nil {
		return
	}

	var resp AcInfoResponse
	// fmt.Println(string(response))
	err = json.Unmarshal([]byte(response), &resp)
	if err != nil {
		return
	}

	if resp.Code <= 0 {
		err = errors.New(resp.Msg)
		return
	}

	for _, v := range resp.Data {
		if v.Account == scheme.Account {
			for _, status := range v.AcStatus {
				// status.ID = Xmd5(Account + "/" + status.Address)
				res = append(res, status)
			}
		}
	}
	return
}

func (a *accountService) SetQueryParam(scheme *Scheme, outParams *AcSetParams) (alert string, err error) {
	var out = make(ValueSet)
	out["account"] 	  = outParams.Account
	out["action"] 	  = outParams.Action
	out["onOff"] 	  = outParams.OnOff
	out["temp"] 	  = outParams.Temp
	out["workMode"]   = outParams.WorkMode
	out["speed"] 	  = outParams.Speed
	out["selectedAc"] = outParams.SelectedAc

	// 请求地址参数，方便调试
	dataSet := url.Values{}
	dataSet.Set("appKey", scheme.AppKey)
	dataSet.Set("timestamp", strconv.FormatInt(MillUnix(), 10))
	dataSet.Set("sign", strings.ToUpper(getSign(scheme, out)))
	for k, v := range out {
		dataSet.Set(k, fmt.Sprintf("%v", v))
	}

	// JSON请求参数
	out["appKey"] = scheme.AppKey
	out["timestamp"] = strconv.FormatInt(MillUnix(), 10)
	out["sign"] = strings.ToUpper(getSign(scheme, out))

	url := scheme.RequestUrl + acSetUri + "?" + dataSet.Encode()
	response, err := HttpPost(url, out, nil)
	//log.Printf("%s,%s", url, response)
	if err != nil {
		return
	}

	var resp AcSetResponse
	err = json.Unmarshal([]byte(response), &resp)
	if err != nil {
		return
	}

	// fmt.Println(string(response), resp)
	if resp.Code <= 0 {
		err = errors.New(resp.Msg)
		return
	}
	alert = string(response)
	return
}

func (a *accountService) GetElecFeeSum(scheme *Scheme, outParams *ElecSumParams) (sum *ElecSum, err error) {
	var out = make(map[string]interface{})
	out["account"] 	  = outParams.Account
	out["address"]	  = outParams.Address
	out["fromDate"]   = outParams.FromDate
	out["toDate"]	  = outParams.ToDate

	params := url.Values{}
	params.Set("appKey", scheme.AppKey)
	params.Set("timestamp", strconv.FormatInt(MillUnix(), 10))
	params.Set("sign", strings.ToUpper(getSign(scheme, out)))
	for k, v := range out {
		params.Set(k, fmt.Sprintf("%v", v))
	}

	url := scheme.RequestUrl + elecFeeSumUri + "?" + params.Encode()
	response, err := HttpGet(url, nil)
	//log.Printf("%s,%s", url, response)
	if err != nil {
		return
	}

	var resp ElecSumResponse
	err = json.Unmarshal([]byte(response), &resp)
	if err != nil {
		return
	}

	// fmt.Println(string(response), resp.Data)
	if resp.Code <= 0 {
		err = errors.New(resp.Msg)
		return
	}
	sum = resp.Data
	return
}

func (a *accountService) LoadQuery(scheme *Scheme) (status []*AcStatus) {
	queryStatus, err := a.GetQueryInfo(scheme)
	if err != nil {
		return
	}

	for _, v := range queryStatus {
		if _, ok := a.queryStatus[v.Address]; !ok {
			status = append(status, v)
		}
	}

	return status
}

func HttpGet(requestUrl string, header url.Values) (content string, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   time.Second * 5, //默认5秒超时时间
		Transport: tr,
	}

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return
	}

	if header != nil {
		for k, _ := range header {
			req.Header.Set(k, fmt.Sprintf("%v", header.Get(k)))
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	content = string(bytes)
	return
}

func HttpPost(requestUrl string, data ValueSet, header url.Values) (content string, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   time.Second * 5, //默认5秒超时时间
		Transport: tr,
	}

	str, _ := json.Marshal(data)
	//log.Printf("HttpPost Json Param: %s", str)
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(string(str)))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type","application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(data)))
	if header != nil {
		for k, _ := range header {
			req.Header.Set(k, fmt.Sprintf("%v", header.Get(k)))
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	content = string(bytes)
	return
}
