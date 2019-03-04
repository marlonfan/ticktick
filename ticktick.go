package ticktick

import (
	"context"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/imroc/req"
)

// Client ticktick请求客户端
type Client struct {
	httpClient *req.Req
	username   string
	password   string
	ctx        context.Context
	user       User
	cookies    []*http.Cookie
}

// NewClient 生成 ticktick Client
func NewClient(username string, password string, cookies []*http.Cookie) (Client, error) {
	c := Client{
		httpClient: req.New(),
		username:   username,
		password:   password,
	}

	c.httpClient.Client().Jar, _ = cookiejar.New(nil)

	if cookies != nil && c.SetCookies(cookies) != nil {
		return c, errors.New("set cookies error")
	}

	return c, c.init()
}

func (c *Client) init() error {
	if c.cookies == nil || c.CheckLoginStatus() == false {
		return c.RefreshCookie()
	}
	return nil
}

// SetCookies 主动设置已登录cookie
func (c *Client) SetCookies(cookies []*http.Cookie) error {
	c.cookies = cookies
	website, err := url.Parse(TickTickWebSiteURL)
	if err != nil {
		return err
	}

	c.httpClient.Client().Jar.SetCookies(website, c.cookies)

	api, err := url.Parse(TickTickAPIRootURL)
	if err != nil {
		return err
	}
	c.httpClient.Client().Jar.SetCookies(api, c.cookies)
	return nil
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

	c.SetCookies(cookies)
	c.user = user

	return nil
}

// Cookie 获取当前登录cookie
func (c *Client) Cookie() []*http.Cookie {
	return c.cookies
}

// CheckLoginStatus 检查当前登录状态
func (c *Client) CheckLoginStatus() bool {
	resp, err := c.httpClient.Get(TickTickUserInfoURL)
	statusCode := resp.Response().StatusCode
	defer resp.Response().Body.Close()

	if err != nil || statusCode != http.StatusOK {
		return false
	}

	return true
}
