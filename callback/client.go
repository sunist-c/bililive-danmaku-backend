package callback

import (
	"log"
	"net/http"
	"net/url"

	"github.com/sunist-c/bililive-danmaku-backend/common/threading"
)

type Client struct {
	client   *http.Client
	Options  *ClientOptions
	executor *threading.Goroutine
	buffer   chan *Context
}

func (c *Client) Execute(exit chan struct{}) {
	for {
		select {
		case <-exit:
			return
		case ctx := <-c.buffer:
			c.sendMessage(ctx)
		}
	}
}

func (c *Client) Stop() (success bool) {
	return true
}

func (c *Client) Close() (success bool) {
	return c.executor.Close()
}

func (c *Client) Serve() (success bool) {
	return c.executor.Serve()
}

func (c *Client) sendMessage(ctx *Context) {
	if ctx == nil {
		return
	}

	u, err := url.JoinPath(c.Options.BaseRouter, ctx.RouterPath)
	if err != nil {
		log.Println("call-back request initialized error: ", err)
		return
	}

	request, err := http.NewRequest(http.MethodPost, u, ctx.Payload)
	if err != nil {
		log.Println("call-back data initialized error: ", err)
		return
	}

	request.Header.Add("Content-Type", "application/json")
	_, err = c.client.Do(request)
	if err != nil {
		log.Println("call-back request send failed: ", err)
	}
}

func (c *Client) SendMessage(ctx *Context) (success bool) {
	select {
	case c.buffer <- ctx:
		return true
	default:
		return false
	}
}

func NewClient() *Client {
	client := &Client{
		client: &http.Client{},
		Options: &ClientOptions{
			BaseRouter:     "http://127.0.0.1:8086",
			DanmakuRouter:  "dm",
			GiftRouter:     "gift",
			GuardRouter:    "welcome",
			MasterRouter:   "custom",
			AudienceRouter: "welcome",
			FansRouter:     "custom",
			CustomMessage:  "custom",
		},
		buffer: make(chan *Context, 128),
	}
	client.executor = threading.NewGoroutine(client)

	return client
}

func NewClientWithOption(opts *ClientOptions) *Client {
	client := &Client{
		client:  &http.Client{},
		Options: opts,
		buffer:  make(chan *Context, 128),
	}
	client.executor = threading.NewGoroutine(client)

	return client
}
