package weixinpay

import "testing"

func getMerchant() *Merchant {
	return NewMerchant("", "", "", "")
}

func TestPlaceOrder(t *testing.T) {
	merchant := getMerchant()
	resp, err := merchant.PlaceOrder(UnifiedOrder{
		Body:           "测试商品",
		OutTradeNo:     "wx_test123",
		TotalFee:       1200,
		SpbillCreateIp: "127.0.0.1",
		NotifyUrl:      "http://www.52177.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ReturnCode != "SUCCESS" || resp.ResultCode != "SUCCESS" {
		t.Fatal("place order fail")
	}
	t.Log(resp)
}

func TestPlaceOrderApp(t *testing.T) {
	merchant := getMerchant()
	resp, err := merchant.PlaceOrderApp("wx_test456", "", "测试商品", "测试商品", "127.0.0.1", "http://www.52177.com", 1200, "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.ReturnCode != "SUCCESS" || resp.ResultCode != "SUCCESS" {
		t.Fatal("place order fail", resp)
	}
	t.Log(resp)
}

func TestPlaceOrderH5(t *testing.T) {
	merchant := getMerchant()
	resp, err := merchant.PlaceOrderH5("wx_test789", "", "测试商品", "测试商品", "127.0.0.1", "http://www.52177.com", 1200, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.ReturnCode != "SUCCESS" || resp.ResultCode != "SUCCESS" {
		t.Fatal("place order fail", resp)
	}
	t.Log(resp)
}

func TestPlaceOrderJSAPI(t *testing.T) {
	merchant := getMerchant()
	resp, err := merchant.PlaceOrderJSAPI("wx_test789", "", "测试商品", "测试商品", "127.0.0.1", "http://www.52177.com", 1200, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.ReturnCode != "SUCCESS" || resp.ResultCode != "SUCCESS" {
		t.Fatal("place order fail", resp)
	}
	t.Log(resp)
}
