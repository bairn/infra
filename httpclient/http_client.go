package httpclient

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bairn/infra/lb"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultHttpTimeout = 30 * time.Second
)

var parseUrl = url.Parse

type Option struct {
	Timeout time.Duration
}

type HttpClient struct {
	client *http.Client
	Option Option
	apps *lb.Apps
}

func NewHttpClient(apps *lb.Apps, opt *Option) *HttpClient {
	c := &HttpClient{
		apps:   apps,
	}
	if opt == nil {
		c.Option = Option{Timeout:defaultHttpTimeout}
	} else {
		c.Option = *opt
	}

	c.client = &http.Client{
		Timeout: c.Option.Timeout,
	}

	return c
}

func (c *HttpClient) NewRequest (method , url string, body io.Reader, headers http.Header) (*http.Request, error) {
	if method == "" {
		method = http.MethodGet
	}

	u, err := parseUrl(url)
	if err != nil {
		return nil, err
	}

	name := u.Host
	app := c.apps.Get(name)
	if app == nil {
		return nil, errors.New("没有可用的微服务应用,应用名称:" + name + ",请求" + url)
	}

	ins := app.Get(url)
	if ins == nil {
		return nil, errors.New("没有可用的应用实例，应用名称：" + name + ",请求：" + url)
	}

	u.Host = ins.Address

	//使用新构造URL创建一个Request
	fmt.Println("修改前", url)
	url = u.String()
	fmt.Println("修改后", url)

	r, err := http.NewRequest(method, url, body)
	if len(headers) > 0 {
		for key, value := range headers {
			for _, val := range value {
				r.Header.Add(key, val)
			}
		}
	}

	return r, err
}

func (c *HttpClient) Do(r *http.Request) (*http.Response, error) {
	res, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return res, err
}


