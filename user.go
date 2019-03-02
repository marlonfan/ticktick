package ticktick

// User 用户
type User struct {
	DS            bool        `json:"ds"`
	InboxID       string      `json:"inboxId"`
	NeedSubscribe bool        `json:"needSubscribe"`
	Pro           bool        `json:"pro"`
	ProEndDate    string      `json:"proEndDate"`
	ProStartDate  interface{} `json:"proStartDate"`
	SubscribeFreq interface{} `json:"subscribeFreq"`
	SubscribeType string      `json:"subscribeType"`
	Token         string      `json:"token"`
	UserCode      interface{} `json:"userCode"`
	UserID        string      `json:"userId"`
	Username      string      `json:"username"`
}

func (c *Client) getUserInfo() User {
	return c.user
}
