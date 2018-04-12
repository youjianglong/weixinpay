package weixinpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var (
	SUCCESS = "SUCCESS"
	FAIL    = "FAIL"
)

type Merchant struct {
	AppId     string // 微信公众账号或开放平台APP的唯一标识
	AppKey    string // API密钥
	AppSecret string // API接口密码
	MchId     string // 微信支付商户号
}

func NewMerchant(appid, appkey, mchid, appsecret string) *Merchant {
	return &Merchant{
		AppId:     appid,
		AppKey:    appkey,
		MchId:     mchid,
		AppSecret: appsecret,
	}
}

func (m *Merchant) IsValid() bool {
	return m.AppId != "" && m.MchId != "" && m.AppKey != ""
}

// sign and return xml
func (m *Merchant) Sign(params Params, signType ...string) Params {
	sType := "md5"
	if len(signType) > 0 && signType[0] != "" {
		sType = signType[0]
	}
	return append(params, Param{"sign", Sign(params, m.AppKey, sType)})
}

var (
	NATIVE = "NATIVE"
	JSAPI  = "JSAPI"
	APP    = "APP"
	WAP    = "WAP"
	MWEB   = "MWEB"
)

//统一下单
func (m *Merchant) PlaceOrder(order UnifiedOrder) (*PlaceOrderResponse, error) {
	if order.NonceStr == "" {
		order.NonceStr = NewNonceString()
	}

	if order.TradeType == "" {
		order.TradeType = MWEB
	}

	if order.Appid == "" {
		order.Appid = m.AppId
	}

	if order.MchId == "" {
		order.MchId = m.MchId
	}

	params := Params{}

	sign := ""

	valueOf := reflect.ValueOf(order)
	typeOf := valueOf.Type()
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}
		tagv := typeOf.Field(i).Tag.Get("xml")
		if tagv == "" || tagv == "xml" || tagv == "-" {
			continue
		}
		val := fmt.Sprintf("%v", field.Interface())
		if tagv == "sign" {
			sign = val
			continue
		}
		params = append(params, Param{tagv, val})
	}

	var postData string

	if sign == "" {
		signType := "md5"
		if order.SignType != nil {
			signType = *order.SignType
		}
		params = m.Sign(params, signType)
	}
	postData = params.ToXmlString()

	data, err := doHttpPost(PlaceOrderUrl, []byte(postData))
	if err != nil {
		return nil, err
	}

	resp, err := ParsePlaceOrderResponse(data)
	if err != nil {
		return nil, err
	}

	if resp.IsSuccess() {
		ok, err := Verify(resp, m.AppKey, resp.Sign)
		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, errors.New("signature error")
		}
	}
	return resp, nil
}

// 统一下单 https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_1 APP
func (m *Merchant) PlaceOrderApp(orderId, product_id, goodsname, desc, clientIp, notifyUrl string, amount int64, attach string) (*PlaceOrderResponse, error) {
	order := UnifiedOrder{
		Body:           goodsname,
		Detail:         &desc,
		OutTradeNo:     orderId,
		TotalFee:       amount,
		SpbillCreateIp: clientIp,
		NotifyUrl:      notifyUrl,
		TradeType:      APP,
	}
	if product_id != "" {
		order.ProductId = &product_id
	}
	if attach != "" {
		order.Attach = &attach
	}
	return m.PlaceOrder(order)
}

// 统一下单 https://pay.weixin.qq.com/wiki/doc/api/H5.php?chapter=9_20 H5
func (m *Merchant) PlaceOrderH5(orderId, product_id, goodsname, desc, clientIp, notifyUrl string, amount int64, scene map[string]string, attach string) (*PlaceOrderResponse, error) {
	order := UnifiedOrder{
		Body:           goodsname,
		Detail:         &desc,
		OutTradeNo:     orderId,
		TotalFee:       amount,
		SpbillCreateIp: clientIp,
		NotifyUrl:      notifyUrl,
		TradeType:      MWEB,
	}
	if scene != nil {
		sceneInfo := map[string]interface{}{
			"h5_info": scene,
		}
		p, _ := json.Marshal(sceneInfo)
		scenejson := string(p)
		order.SceneInfo = &scenejson
	}
	if product_id != "" {
		order.ProductId = &product_id
	}
	if attach != "" {
		order.Attach = &attach
	}
	return m.PlaceOrder(order)
}

// 统一下单 https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_1 JSAPI
func (m *Merchant) PlaceOrderJSAPI(orderId, product_id, goodsname, desc, clientIp, notifyUrl string, amount int64, openID string, attach string) (*PlaceOrderResponse, error) {
	order := UnifiedOrder{
		Body:           goodsname,
		Detail:         &desc,
		NotifyUrl:      notifyUrl,
		OutTradeNo:     orderId,
		TotalFee:       amount,
		SpbillCreateIp: clientIp,
		TradeType:      JSAPI,
		Openid:         &openID,
	}

	if attach != "" {
		order.Attach = &attach
	}

	if product_id != "" {
		order.ProductId = &product_id
	}

	return m.PlaceOrder(order)
}

func (m *Merchant) CloseOrder(orderId string) (*CloseOrderResponse, error) {
	var params = Params{
		{"appid", m.AppId},
		{"mch_id", m.MchId},
		{"out_trade_no", orderId},
		{"nonce_str", NewNonceString()},
	}

	data, err := doHttpPost(CloseOrderUrl, []byte(m.Sign(params).ToXmlString()))
	if err != nil {
		return nil, err
	}

	return ParseCloseOrderResponse(data)
}

// 根据微信支付订单号查询订单 https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_2
func (m *Merchant) QueryOrderByTransId(transId string) (*PayResult, error) {
	var params = Params{
		{"appid", m.AppId},
		{"mch_id", m.MchId},
		{"nonce_str", NewNonceString()},
		{"transaction_id", transId},
	}

	data, err := doHttpPost(QueryOrderUrl, []byte(m.Sign(params).ToXmlString()))
	if err != nil {
		return nil, err
	}

	return ParsePayResult(data)
}

// 根据商户订单号查询订单 https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_2
func (m *Merchant) QueryOrderByOrderId(orderId string) ([]byte, *PayResult, error) {
	var params = Params{
		{"appid", m.AppId},
		{"mch_id", m.MchId},
		{"nonce_str", NewNonceString()},
		{"out_trade_no", orderId},
	}

	data, err := doHttpPost(QueryOrderUrl, []byte(m.Sign(params).ToXmlString()))
	if err != nil {
		return nil, nil, err
	}

	res, err := ParsePayResult(data)
	if err != nil {
		return nil, nil, err
	}
	return data, res, nil
}

// 生成二维码链接
// weixin：//wxpay/bizpayurl?sign=XXXXX&appid=XXXXX&mch_id=XXXXX&product_id=XXXXXX&time_stamp=XXXXXX&nonce_str=XXXXX
func (m *Merchant) GenQRLink(productId string) string {
	var params = Params{
		{"appid", m.AppId},
		{"mch_id", m.MchId},
		{"nonce_str", NewNonceString()},
		{"product_id", productId},
		{"time_stamp", NewTimestampString()},
	}

	sign := Sign(params, m.AppKey)
	params = append(params, Param{"sign", sign})
	return fmt.Sprintf("weixin://wxpay/bizpayurl?%s", params.ToQueryString())
}

func (m *Merchant) NewScanResponse(returnCode, returnMsg, prepayId, resultCode, errCodeDes string) *ScanResponse {
	return &ScanResponse{
		Params: Params{
			{"return_code", returnCode},
			{"return_msg", returnMsg},
			{"appid", m.AppId},
			{"mch_id", m.MchId},
			{"nonce_str", NewNonceString()},
			{"prepay_id", prepayId},
			{"result_code", resultCode},
			{"err_code_des", errCodeDes},
		},
		AppKey: m.AppKey,
	}
}
