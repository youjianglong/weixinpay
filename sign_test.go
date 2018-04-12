package weixinpay

import (
	"testing"
)

func TestSign(t *testing.T) {
	var (
		params = Params{
			{"appid", "wxd930ea5d5a258f4f"},
			{"mch_id", "10000100"},
			{"device_info", "1000"},
			{"body", "test"},
			{"nonce_str", "ibuaiVcKdpRxkhJA"},
		}

		key         = "192006250b4c09247ec02edce69f6a2d"
		correctSign = "9A0A8659F005D6984697E2CA0A9CF3B7"
		correctXml  = `<xml><appid>wxd930ea5d5a258f4f</appid><body>test</body><device_info>1000</device_info><mch_id>10000100</mch_id><nonce_str>ibuaiVcKdpRxkhJA</nonce_str><sign>9A0A8659F005D6984697E2CA0A9CF3B7</sign></xml>`
	)

	sign := Sign(params, key)
	if sign != correctSign {
		t.Fatal("sign error:", sign, "is not equal to", correctSign)
		return
	}

	params = append(params, Param{"sign", sign})
	if params.ToXmlString() != correctXml {
		t.Fatal("sign error:", params.ToXmlString(), "is not equal to", correctXml)
		return
	}
}

func TestSign2(t *testing.T) {
	var (
		params = Params{
			{"appid", "wxd930ea5d5a258f4f"},
			{"mch_id", "10000100"},
			{"device_info", "1000"},
			{"body", "test"},
			{"nonce_str", "ibuaiVcKdpRxkhJA"},
		}

		key         = "192006250b4c09247ec02edce69f6a2d"
		correctSign = "6A9AE1657590FD6257D693A078E1C3E4BB6BA4DC30B23E0EE2496E54170DACD6"
		correctXml  = `<xml><appid>wxd930ea5d5a258f4f</appid><body>test</body><device_info>1000</device_info><mch_id>10000100</mch_id><nonce_str>ibuaiVcKdpRxkhJA</nonce_str><sign>6A9AE1657590FD6257D693A078E1C3E4BB6BA4DC30B23E0EE2496E54170DACD6</sign></xml>`
	)

	sign := Sign(params, key, "HMAC-SHA256")
	if sign != correctSign {
		t.Fatal("sign error:", sign, "is not equal to", correctSign)
		return
	}

	params = append(params, Param{"sign", sign})
	if params.ToXmlString() != correctXml {
		t.Fatal("sign error:", params.ToXmlString(), "is not equal to", correctXml)
		return
	}
}
