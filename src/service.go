package account_service

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	xutil "github.com/donetkit/wechat/util"
	"github.com/gin-gonic/gin"
)

const (
	showAcInfoURL            = "/open/ac/water/showAcInfo"
	acSetURL				 = "/open/ac/water/acSet"
	elecFeeDailyURL          = "/open/ac/water/elecFeeSum"
)

var AccountService = new(accountService)

// Account Service
type accountService struct {
	*gin.Context
	wg	*sync.WaitGroup
}

func (a *accountService) GetQueryInfo(scheme *Scheme) (res []*AcStatus, err error) {
	params := url.Values{}
	params.Set("appKey", scheme.AppKey)
	params.Set("timestamp", strconv.FormatInt(MillUnix(), 10))
	params.Set("sign", strings.ToUpper(getSign(scheme, nil)))

	uri := scheme.RequestUrl + showAcInfoURL + "?" + params.Encode()
	response, err := xutil.HTTPGet(uri)
	if err != nil {
		return
	}

	var resp AcInfoResponse
	// fmt.Println(string(response))
	err = json.Unmarshal(response, &resp)
	if err != nil {
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

func (a *accountService) SetQueryParam(scheme *Scheme, outParams *AcSetParams) (alert bool, err error) {
	var out = make(map[string]interface{})
	out["account"] 	  = outParams.Account
	out["action"] 	  = outParams.Action
	out["onOff"] 	  = outParams.OnOff
	out["temp"] 	  = outParams.Temp
	out["workMode"]   = outParams.WorkMode
	out["speed"] 	  = outParams.Speed
	out["selectedAc"] = outParams.SelectedAc

	params := url.Values{}
	params.Set("appKey", scheme.AppKey)
	params.Set("timestamp", strconv.FormatInt(MillUnix(), 10))
	params.Set("sign", strings.ToUpper(getSign(scheme, out)))
	for k, v := range out {
		params.Set(k, fmt.Sprintf("%v", v))
	}

	uri := scheme.RequestUrl + acSetURL + "?" + params.Encode()
	response, err := xutil.HTTPPost(uri, params.Encode())
	if err != nil {
		return
	}

	var resp AcSetResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		return
	}

	// fmt.Println(string(response), resp)
	if resp.Code > 0 {
		return true, nil
	}
	return
}

func (a *accountService) GetDailyInfo(scheme *Scheme, outParams *DailySumParams) (sum *DailySum, err error) {
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

	uri := scheme.RequestUrl + elecFeeDailyURL + "?" + params.Encode()
	response, err := xutil.HTTPGet(uri)
	if err != nil {
		return
	}

	var resp DailySumResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		return
	}

	// fmt.Println(string(response), resp.Data)
	if resp.Code > 0 {
		sum, err = resp.Data, nil
	}
	return
}

func (a *accountService) LoadQuery(scheme *Scheme) (status []*AcStatus) {
	status, err := a.GetQueryInfo(scheme)
	if err != nil {
		return
	}

	return status
}
