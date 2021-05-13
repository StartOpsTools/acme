package acme

// 4 新建订单

type NewOrderRequest struct {
	protected string
	payload   string
	signature string
}

// protected
type NewOrderRequestProtected struct {
	nonce string
	url   string
	alg   string // RS256
	kid   string
}

// payload
type NewOrderRequestPayload struct {
	identifiers []NewOrderRequestPayloadIdentifiers
}
type NewOrderRequestPayloadIdentifiers struct {
	Type  string `json:"type"` // dns, ipv4
	Value string `json:"value"`
}

// response
type NewOrderResponse struct {
	Status         string                        `json:"status"`         //	pending/ready/processing/valid/invalid
	Expires        string                        `json:"expires"`        //	订单失效时间, 由服务商或CA决定
	Identifiers    []NewOrderResponseIdentifiers `json:"identifiers"`    //
	Authorizations []string                      `json:"authorizations"` // 订单需要依次完成的授权验证资源（Auth-Z）的链接    不允许为空数组（必须至少有一个流程）
	Finalize       string                        `json:"finalize"`       //  授权验证完成后，调用finalize接口签发证书（包括CSR也是在这一步提交的）
}

type NewOrderResponseIdentifiers struct {
	Type  string `json:"type"`  //
	Value string `json:"value"` //
}

func NewOrder() {

}
