package weixinpay

import "encoding/xml"

//统一下单
type UnifiedOrder struct {
	XMLName xml.Name `xml:"xml"`

	//微信分配的公众账号ID（企业号corpid即为此appId）
	Appid string `xml:"appid"`

	//微信支付分配的商户号
	MchId string `xml:"mch_id"`

	//终端设备号(门店号或收银设备ID)，注意：PC网页或公众号内支付请传"WEB"
	DeviceInfo *string `xml:"device_info"`

	//随机字符串，不长于32位
	NonceStr string `xml:"nonce_str"`

	//签名
	Sign *string `xml:"sign"`

	//签名类型
	SignType *string `xml:"sign_type"`

	//商品简单描述，该字段须严格按照规范传递
	Body string `xml:"body"`

	//商品详情
	Detail *string `xml:"detail"`

	//附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
	Attach *string `xml:"attach"`

	//商户系统内部的订单号,32个字符内、可包含字母
	OutTradeNo string `xml:"out_trade_no"`

	//货币类型
	//符合ISO 4217标准的三位字母代码，默认人民币：CNY
	FeeType *string `xml:"fee_type"`

	//订单总金额，单位为分
	TotalFee int64 `xml:"total_fee"`

	//必须传正确的用户端IP
	SpbillCreateIp string `xml:"spbill_create_ip"`

	//订单生成时间，格式为yyyyMMddHHmmss
	//如2009年12月25日9点10分10秒表示为20091225091010
	TimeStart *string `xml:"time_start"`

	//订单失效时间，格式为yyyyMMddHHmmss
	//如2009年12月27日9点10分10秒表示为20091227091010
	TimeExpire *string `xml:"time_expire"`

	//商品标记，代金券或立减优惠功能的参数
	GoodsTag *string `xml:"goods_tag"`

	//接收微信支付异步通知回调地址，通知url必须为直接可访问的url，不能携带参数。
	NotifyUrl string `xml:"notify_url"`

	//交易类型
	TradeType string `xml:"trade_type"`

	//商品ID
	//trade_type=NATIVE，此参数必传。此id为二维码中包含的商品ID，商户自行定义。
	ProductId *string `xml:"product_id"`

	//指定支付方式
	//no_credit--指定不能使用信用卡支付
	LimitPay *string `xml:"limit_pay"`

	//用户标识
	//trade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识
	Openid *string `xml:"openid"`

	//场景信息
	//该字段用于上报支付的场景信息,针对H5支付有以下三种场景
	//IOS移动应用
	//{"h5_info": {"type":"IOS","app_name": "王者荣耀","bundle_id": "com.tencent.wzryIOS"}}
	//安卓移动应用
	//{"h5_info": {"type":"Android","app_name": "王者荣耀","package_name": "com.tencent.tmgp.sgame"}}
	//WAP网站应用
	//{"h5_info": {"type":"Wap","wap_url": "https://pay.qq.com","wap_name": "腾讯充值"}}
	SceneInfo *string `xml:"scene_info"`
}

// PlaceOrderResponse represent place order reponse message from weixin pay.
// For field explanation refer to: http://pay.weixin.qq.com/wiki/doc/api/app.php?chapter=9_1
type PlaceOrderResponse struct {
	XMLName     xml.Name `xml:"xml"`
	ReturnCode  string   `xml:"return_code"`
	ReturnMsg   string   `xml:"return_msg"`
	AppId       string   `xml:"appid"`
	MchId       string   `xml:"mch_id"`
	NonceStr    string   `xml:"nonce_str"`
	Sign        string   `xml:"sign"`
	ResultCode  string   `xml:"result_code"`
	ErrCode     string   `xml:"err_code"`
	ErrCodeDesc string   `xml:"err_code_des"`
	DeviceInfo  string   `xml:"device_info"`
	TradeType   string   `xml:"trade_type"`
	PrepayId    string   `xml:"prepay_id"`
	CodeUrl     string   `xml:"code_url"`
	MWebURL     string   `xml:"mweb_url"`
}

// Parse the reponse message from weixin pay to struct of PlaceOrderResult
func ParsePlaceOrderResponse(data []byte) (*PlaceOrderResponse, error) {
	var resp PlaceOrderResponse
	err := xml.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (p *PlaceOrderResponse) IsSuccess() bool {
	return p.ReturnCode == "SUCCESS"
}

func (p *PlaceOrderResponse) Error() *Error {
	if !p.IsSuccess() {
		return GetError(p.ErrCode)
	}
	return nil
}

type CloseOrderResponse struct {
	XMLName     xml.Name `xml:"xml"`
	ReturnCode  string   `xml:"return_code"`
	ReturnMsg   string   `xml:"return_msg"`
	AppId       string   `xml:"appid"`
	MchId       string   `xml:"mch_id"`
	NonceStr    string   `xml:"nonce_str"`
	Sign        string   `xml:"sign"`
	ResultCode  string   `xml:"result_code"`
	ErrCode     string   `xml:"err_code"`
	ErrCodeDesc string   `xml:"err_code_des"`
}

func ParseCloseOrderResponse(data []byte) (*CloseOrderResponse, error) {
	var resp CloseOrderResponse
	err := xml.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
