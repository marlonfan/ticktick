package ticktick

import (
	"context"
	"errors"
	"net/http"

	"github.com/imroc/req"
)

// Client ticktick请求客户端
type Client struct {
	httpClient *req.Req
	username   string
	password   string
	ctx        context.Context
	user       User
	cookie     *http.Cookie
}

// NewClient 生成 ticktick Client
func NewClient(username string, password string, cookie string) (Client, error) {
	c := Client{
		httpClient: req.New(),
		username:   username,
		password:   password,
	}

	if cookie != "" {
		c.cookie = parseCookie(cookie)[0]
		return c, nil
	}

	return c, c.init()
}

func (c *Client) init() error {
	return c.RefreshCookie()
}

// Context 获取context
func (c *Client) Context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.Background()
}

// RefreshCookie 登录
func (c *Client) RefreshCookie() error {
	user := User{}

	params := map[string]string{
		"username": c.username,
		"password": c.password,
	}
	resp, err := c.httpClient.Post(TickTickLoginURL, req.BodyJSON(&params))
	if err != nil {
		return err
	}

	if resp.Response().StatusCode != 200 {
		return errors.New(resp.String())
	}

	resp.ToJSON(&user)

	cookies := resp.Response().Cookies()
	defer resp.Response().Body.Close()

	if len(cookies) != 2 {
		return errors.New("Cookie error")
	}

	c.cookie = cookies[1]
	c.user = user

	return nil
}

// CheckLoginStatus 检查当前登录状态
func (c *Client) CheckLoginStatus() bool {
	if c.cookie == nil {
		return false
	}

	resp, err := c.httpClient.Get(TickTickUserInfoURL, c.cookie)
	statusCode := resp.Response().StatusCode
	defer resp.Response().Body.Close()

	if err != nil || statusCode != http.StatusOK {
		return false
	}

	return true
}

func parseCookie(cookies string) []*http.Cookie {
	return (&http.Response{
		Header: http.Header{"Set-Cookie": {cookies}},
	}).Cookies()
}
