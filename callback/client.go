package callback

import (
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	client  *http.Client
	Options *ClientOptions
}

func NewClient(baseRouter, danmakuRouter, giftRouter, guardRouter, masterRouter, audienceRouter, fansRouter, customMessageRouter string) *Client {
	return &Client{
		client: &http.Client{},
		Options: &ClientOptions{
			BaseRouter:     baseRouter,
			DanmakuRouter:  danmakuRouter,
			GiftRouter:     giftRouter,
			GuardRouter:    guardRouter,
			MasterRouter:   masterRouter,
			AudienceRouter: audienceRouter,
			FansRouter:     fansRouter,
			CustomMessage:  customMessageRouter,
		},
	}
}

func NewClientWithOption(option *ClientOptions) *Client {
	return &Client{
		client:  &http.Client{},
		Options: option,
	}
}

func (c *Client) SendMessage(ctx *Context) {
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

	//log.Printf("call-back request will sent to %v\n", u)

	request.Header.Add("Content-Type", "application/json")
	_, err = c.client.Do(request)
	if err != nil {
		log.Println("call-back request send failed: ", err)
	} else {
		//log.Printf("call-back request sent successfully to %v\n", u)
	}
}
