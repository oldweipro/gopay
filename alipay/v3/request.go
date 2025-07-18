package alipay

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-pay/crypto/aes"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/xhttp"
	"github.com/go-pay/util"
)

var defaultRequestIdFunc = &requestIdFunc{}

type requestIdFunc struct{}

func (d *requestIdFunc) RequestId() string {
	return fmt.Sprintf("%s-%d", util.RandomString(21), time.Now().Unix())
}

// DoAliPayAPISelfV3 支付宝接口自行实现方法
func (a *ClientV3) DoAliPayAPISelfV3(ctx context.Context, method, path string, bm gopay.BodyMap, aliRsp any) (res *http.Response, err error) {
	var (
		bs            []byte
		authorization string
	)
	aat := bm.GetString(HeaderAppAuthToken)
	switch method {
	case MethodGet:
		bm.Remove(HeaderAppAuthToken)
		uri := path + "?" + bm.EncodeURLParams()
		authorization, err = a.authorization(MethodGet, uri, nil, aat)
		if err != nil {
			return nil, err
		}
		res, bs, err = a.doGet(ctx, uri, authorization, aat)
		if err != nil {
			return nil, err
		}
	case MethodPost:
		authorization, err = a.authorization(MethodPost, path, bm, aat)
		if err != nil {
			return nil, err
		}
		res, bs, err = a.doPost(ctx, bm, path, authorization, aat)
		if err != nil {
			return nil, err
		}
	case MethodPatch:
		authorization, err = a.authorization(MethodPatch, path, bm, aat)
		if err != nil {
			return nil, err
		}
		res, bs, err = a.doPatch(ctx, bm, path, authorization, aat)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("method:%s not support", method)
	}
	if err = json.Unmarshal(bs, aliRsp); err != nil {
		return nil, err
	}
	return res, nil
}

func (a *ClientV3) doPost(ctx context.Context, bm gopay.BodyMap, uri, authorization, aat string) (res *http.Response, bs []byte, err error) {
	var url = v3BaseUrlCh + uri
	if !a.IsProd {
		url = v3SandboxBaseUrl + uri
	}
	if a.proxyHost != "" {
		url = a.proxyHost + uri
	}
	req := a.hc.Req() // default json
	req.Header.Add(HeaderAuthorization, authorization)
	req.Header.Add(HeaderRequestID, a.requestIdFunc.RequestId())
	req.Header.Add(HeaderSdkVersion, "gopay/"+gopay.Version)
	if aat != gopay.NULL {
		req.Header.Add(HeaderAppAuthToken, aat)
	} else if a.AppAuthToken != "" {
		req.Header.Add(HeaderAppAuthToken, a.AppAuthToken)
	}
	req.Header.Add("Accept", "application/json")
	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Url: %s", url)
		a.logger.Debugf("Alipay_V3_Req_Body: %s", bm.JsonBody())
		a.logger.Debugf("Alipay_V3_Req_Headers: %#v", req.Header)
	}
	res, bs, err = req.Post(url).SendBodyMap(bm).EndBytesForAlipayV3(ctx)
	if err != nil {
		return nil, nil, err
	}

	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Response: %d >> %s", res.StatusCode, string(bs))
		a.logger.Debugf("Alipay_V3_Rsp_Headers: %#v", res.Header)
	}
	return res, bs, nil
}

func (a *ClientV3) doPatch(ctx context.Context, bm gopay.BodyMap, uri, authorization, aat string) (res *http.Response, bs []byte, err error) {
	var url = v3BaseUrlCh + uri
	if !a.IsProd {
		url = v3SandboxBaseUrl + uri
	}
	if a.proxyHost != "" {
		url = a.proxyHost + uri
	}
	req := a.hc.Req() // default json
	req.Header.Add(HeaderAuthorization, authorization)
	req.Header.Add(HeaderRequestID, a.requestIdFunc.RequestId())
	req.Header.Add(HeaderSdkVersion, "gopay/"+gopay.Version)
	if aat != gopay.NULL {
		req.Header.Add(HeaderAppAuthToken, aat)
	} else if a.AppAuthToken != "" {
		req.Header.Add(HeaderAppAuthToken, a.AppAuthToken)
	}
	req.Header.Add("Accept", "application/json")
	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Url: %s", url)
		a.logger.Debugf("Alipay_V3_Req_Body: %s", bm.JsonBody())
		a.logger.Debugf("Alipay_V3_Req_Headers: %#v", req.Header)
	}
	res, bs, err = req.Patch(url).SendBodyMap(bm).EndBytesForAlipayV3(ctx)
	if err != nil {
		return nil, nil, err
	}

	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Response: %d >> %s", res.StatusCode, string(bs))
		a.logger.Debugf("Alipay_V3_Rsp_Headers: %#v", res.Header)
	}
	return res, bs, nil
}

func (a *ClientV3) doPut(ctx context.Context, bm gopay.BodyMap, uri, authorization, aat string) (res *http.Response, bs []byte, err error) {
	var url = v3BaseUrlCh + uri
	if !a.IsProd {
		url = v3SandboxBaseUrl + uri
	}
	if a.proxyHost != "" {
		url = a.proxyHost + uri
	}
	req := a.hc.Req() // default json
	req.Header.Add(HeaderAuthorization, authorization)
	req.Header.Add(HeaderRequestID, a.requestIdFunc.RequestId())
	req.Header.Add(HeaderSdkVersion, "gopay/"+gopay.Version)
	if aat != gopay.NULL {
		req.Header.Add(HeaderAppAuthToken, aat)
	} else if a.AppAuthToken != "" {
		req.Header.Add(HeaderAppAuthToken, a.AppAuthToken)
	}
	req.Header.Add("Accept", "application/json")
	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Url: %s", url)
		a.logger.Debugf("Alipay_V3_Req_Body: %s", bm.JsonBody())
		a.logger.Debugf("Alipay_V3_Req_Headers: %#v", req.Header)
	}
	res, bs, err = req.Put(url).SendBodyMap(bm).EndBytesForAlipayV3(ctx)
	if err != nil {
		return nil, nil, err
	}

	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Response: %d > %s", res.StatusCode, string(bs))
		a.logger.Debugf("Alipay_V3_Rsp_Headers: %#v", res.Header)
	}
	return res, bs, nil
}

func (a *ClientV3) doGet(ctx context.Context, uri, authorization, aat string) (res *http.Response, bs []byte, err error) {
	var url = v3BaseUrlCh + uri
	if !a.IsProd {
		url = v3SandboxBaseUrl + uri
	}
	if a.proxyHost != "" {
		url = a.proxyHost + uri
	}
	req := a.hc.Req() // default json
	req.Header.Add(HeaderAuthorization, authorization)
	req.Header.Add(HeaderRequestID, a.requestIdFunc.RequestId())
	req.Header.Add(HeaderSdkVersion, "gopay/"+gopay.Version)
	if aat != gopay.NULL {
		req.Header.Add(HeaderAppAuthToken, aat)
	} else if a.AppAuthToken != "" {
		req.Header.Add(HeaderAppAuthToken, a.AppAuthToken)
	}
	req.Header.Add("Accept", "application/json")
	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Url: %s", url)
		a.logger.Debugf("Alipay_V3_Req_Headers: %#v", req.Header)
	}
	res, bs, err = req.Get(url).EndBytesForAlipayV3(ctx)
	if err != nil {
		return nil, nil, err
	}

	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Response: %d > %s", res.StatusCode, string(bs))
		a.logger.Debugf("Alipay_V3_Rsp_Headers: %#v", res.Header)
	}
	return res, bs, nil
}

func (a *ClientV3) doProdPostFile(ctx context.Context, bm gopay.BodyMap, uri, authorization, aat string) (res *http.Response, bs []byte, err error) {
	var url = v3BaseUrlCh + uri
	if a.proxyHost != "" {
		url = a.proxyHost + uri
	}
	req := a.hc.Req(xhttp.TypeMultipartFormData)
	req.Header.Add(HeaderAuthorization, authorization)
	req.Header.Add(HeaderRequestID, a.requestIdFunc.RequestId())
	req.Header.Add(HeaderSdkVersion, "gopay/"+gopay.Version)
	if aat != gopay.NULL {
		req.Header.Add(HeaderAppAuthToken, aat)
	} else if a.AppAuthToken != "" {
		req.Header.Add(HeaderAppAuthToken, a.AppAuthToken)
	}
	req.Header.Add("Accept", "application/json")
	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Url: %s", url)
		a.logger.Debugf("Alipay_V3_Req_Body: %s", bm.JsonBody())
		a.logger.Debugf("Alipay_V3_Req_Headers: %#v", req.Header)
	}
	res, bs, err = req.Post(url).SendMultipartBodyMap(bm).EndBytesForAlipayV3(ctx)
	if err != nil {
		return nil, nil, err
	}
	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Response: %d >> %s", res.StatusCode, string(bs))
		a.logger.Debugf("Alipay_V3_Rsp_Headers: %#v", res.Header)
	}
	return res, bs, nil
}

func (a *ClientV3) doDelete(ctx context.Context, bm gopay.BodyMap, uri, authorization, aat string) (res *http.Response, bs []byte, err error) {
	var url = v3BaseUrlCh + uri
	if !a.IsProd {
		url = v3SandboxBaseUrl + uri
	}
	if a.proxyHost != "" {
		url = a.proxyHost + uri
	}
	req := a.hc.Req() // default json
	req.Header.Add(HeaderAuthorization, authorization)
	req.Header.Add(HeaderRequestID, a.requestIdFunc.RequestId())
	req.Header.Add(HeaderSdkVersion, "gopay/"+gopay.Version)
	if aat != gopay.NULL {
		req.Header.Add(HeaderAppAuthToken, aat)
	} else if a.AppAuthToken != "" {
		req.Header.Add(HeaderAppAuthToken, a.AppAuthToken)
	}
	req.Header.Add("Accept", "application/json")
	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Url: %s", url)
		a.logger.Debugf("Alipay_V3_Req_Body: %s", bm.JsonBody())
		a.logger.Debugf("Alipay_V3_Req_Headers: %#v", req.Header)
	}
	res, bs, err = req.Delete(url).SendBodyMap(bm).EndBytesForAlipayV3(ctx)
	if err != nil {
		return nil, nil, err
	}

	if a.DebugSwitch == gopay.DebugOn {
		a.logger.Debugf("Alipay_V3_Response: %d >> %s", res.StatusCode, string(bs))
		a.logger.Debugf("Alipay_V3_Rsp_Headers: %#v", res.Header)
	}
	return res, bs, nil
}

func (a *ClientV3) encryptBizContent(originData string) (string, error) {
	encryptData, err := aes.CBCEncrypt([]byte(originData), []byte(a.aesKey), a.ivKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptData), nil
}
