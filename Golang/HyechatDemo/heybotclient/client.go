package heybotclient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"HyechatDemo/heybotclient/eventinterface"
	"HyechatDemo/heybotclient/model"
)

const (
	PingInterval  = 30 * time.Second
	CheckInterval = 30 * time.Second
	SleepTime     = 1 * time.Second
	MaxSleepTime  = 60 * time.Second
)

type WebSocketClient struct {
	conn         *websocket.Conn
	isConnected  bool
	close        bool
	messageChan  chan []byte
	eventHandler eventinterface.EventHandler
	ctx          context.Context
	cancel       context.CancelFunc
	proxyURL     *url.URL
	token        string
}

func NewWebSocketClient(ctx context.Context, token string, proxyURL *url.URL, eventHandler eventinterface.EventHandler) *WebSocketClient {
	ctx, cancel := context.WithCancel(ctx)
	client := &WebSocketClient{
		proxyURL:     proxyURL,
		messageChan:  make(chan []byte, 1000),
		isConnected:  false,
		eventHandler: eventHandler,
		ctx:          ctx,
		cancel:       cancel,
		token:        token,
	}
	return client
}

func (c *WebSocketClient) Close() error {
	c.close = true
	c.cancel()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *WebSocketClient) Connect() error {
	header := http.Header{
		"Accept":          {"application/json, text/plain, */*"},
		"Accept-Language": {"zh-CN,zh;q=0.9"},
		"Cache-Control":   {"no-cache"},
		"Pragma":          {"no-cache"},
	}

	dialer := &websocket.Dialer{
		TLSClientConfig: &tls.Config{},
	}
	if c.proxyURL != nil {
		dialer.Proxy = http.ProxyURL(c.proxyURL)
	}

	conn, _, err := dialer.Dial(model.GetWssUrl(c.token), header)
	if err != nil {
		return err
	}

	c.conn = conn
	c.isConnected = true

	go func() {
		c.heartbeat()
	}()
	go func() {
		c.receive()
	}()
	go func() {
		c.handle()
	}()

	return nil
}

func (c *WebSocketClient) heartbeat() {
	ticker := time.NewTicker(PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if !c.isConnected {
				continue
			}
			if err := c.SendPing(); err != nil {
				log.Printf("Send PING error: %v", err)
				c.Reconnect()
				return
			}
		}
	}
}

func (c *WebSocketClient) SendPing() error {
	return c.conn.WriteMessage(websocket.PingMessage, []byte("PING"))
}

func (c *WebSocketClient) receive() {
	for {
		msgType, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			c.Reconnect()
			return
		}
		if msgType != websocket.TextMessage {
			// Only process text messages
			continue
		}
		select {
		case <-c.ctx.Done():
			return
		case c.messageChan <- message:
		}
	}
}

func (c *WebSocketClient) Reconnect() {
	c.isConnected = false
	sleepTime := SleepTime
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}
		if c.close {
			return
		}
		time.Sleep(sleepTime)
		if sleepTime < MaxSleepTime {
			sleepTime *= 2
		} else {
			sleepTime = SleepTime
		}
		if !c.isConnected {
			log.Println("Reconnecting to the server...")
			if err := c.Connect(); err != nil {
				log.Printf("Failed to reconnect: %v", err)
				continue
			}
			log.Println("Reconnected to the server")
			break
		}
	}
}

func (c *WebSocketClient) HandleData(msgData []byte) error {
	if strings.HasPrefix(string(msgData), "PONG") || strings.HasPrefix(string(msgData), "pong") {
		return nil
	}
	reqStruct := &model.GenericType{}
	if err := json.Unmarshal(msgData, reqStruct); err != nil {
		log.Printf("ReceiveData Failed: %v", err)
		return err
	}
	return c.eventHandler.OnMessage(c.ctx, reqStruct)
}

func (c *WebSocketClient) handle() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case data := <-c.messageChan:
			log.Printf("ReceiveData: %s", string(data))
			if c.close {
				log.Println("ReceiveData Failed: close")
				return
			}
			go func() {
				if err := c.HandleData(data); err != nil {
					log.Printf("HandleData Failed: %v", err)
				}
			}()
		}
	}
}
