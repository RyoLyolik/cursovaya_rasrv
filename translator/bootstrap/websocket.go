package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func WSConnect(ctx context.Context, wsCfg Websocket) (*websocket.Conn, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	address := fmt.Sprintf("%s:%s", wsCfg.Host, wsCfg.Port)
	u := url.URL{Scheme: wsCfg.Scheme, Host: address, Path: wsCfg.Path}
	connection := make(chan *websocket.Conn, 1)
	go func() {
		for {
			c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Println("dial failed, retrying:", err)
				time.Sleep(1 * time.Second)
				continue
			}
			log.Println("connected successfully")
			connection <- c
			break
		}
	}()
	select {
	case <-ctxWithTimeout.Done():
		return nil, fmt.Errorf("failed to connect")
	case c := <-connection:
		return c, nil
	}
}
