package rmqtool

import (
	"sync"
	"time"

	"github.com/streadway/amqp"
)

var (
	DefaultPublisherRetryTime   time.Duration = 3 * time.Second
	DefaultPublisherChannelSize int32         = 2000
)

type PublisherTool struct {
	channels        map[string]*publishChannel
	RetryTime       time.Duration
	url             string
	exchange        string
	safeChannelMaps *sync.Map
}

func NewPublisherTool(url, exchange string, routeKeys []string) (*PublisherTool, error) {
	tool := &PublisherTool{
		channels:        make(map[string]*publishChannel, 0), //channels
		RetryTime:       DefaultPublisherRetryTime,           //default retry
		url:             url,
		exchange:        exchange,
		safeChannelMaps: new(sync.Map),
	}
	err := tool.conn(url, exchange, routeKeys)
	return tool, err
}

func (c *PublisherTool) conn(url, exchange string, routeKeys []string) error {
	//test link
	testConn, err := amqp.Dial(url)
	if testConn != nil {
		go testConn.Close()
	} //close test conn
	if err != nil {
		return err
	}
	for _, route := range routeKeys {
		c.channels[route] = &publishChannel{
			amqpUrl:  url,
			exchange: exchange,
			routeKey: route,
			Channel:  make(chan interface{}, DefaultPublisherChannelSize),
		}
		go c.channels[route].Receive()
	}
	return nil
}

func (c *PublisherTool) GetSafeChannel(route string) (*amqp.Channel, error) {
	if channel, ok := c.safeChannelMaps.Load(route); ok {
		return channel.(*amqp.Channel), nil
	} else {
		conn, err := amqp.Dial(c.url)
		if err != nil {
			go conn.Close()
			return nil, err
		}
		channel, err := conn.Channel()
		if err != nil {
			go conn.Close()
			return nil, err
		}
		c.safeChannelMaps.Store(route, channel)
		return channel, nil
	}

}

func (c *PublisherTool) SafePublish(route string, msg amqp.Publishing) error {
	channel, err := c.GetSafeChannel(route)
	if err != nil {
		return err
	} else {
		err := channel.Publish(c.exchange, route, false, false, msg)
		if err != nil {
			c.safeChannelMaps.Delete(route)
			go channel.Close()
		}
		return err
	}
}

func (c *PublisherTool) publish(route string, msg interface{}) {
	if s, ok := c.channels[route]; ok {
		s.Channel <- msg
	}
}

func (c *PublisherTool) Publish(route string, msg amqp.Publishing) {
	c.publish(route, msg)
}

func (c *PublisherTool) PublishExt(route, fix string, msg amqp.Publishing) {
	c.publish(route, &publishingExt{
		routeKeyFix: fix,
		msg:         msg,
	})
}

type publishingExt struct {
	routeKeyFix string
	msg         amqp.Publishing
}

func (c *publishingExt) Key(prefix string) string {
	return prefix + c.routeKeyFix
}

func (c *publishingExt) Msg() amqp.Publishing {
	return c.msg
}

type publishChannel struct {
	amqpUrl   string
	exchange  string
	routeKey  string
	retryTime time.Duration
	Channel   chan interface{}
}

func (c *publishChannel) Receive() {
RetryConnect:
	conn, err := amqp.Dial(c.amqpUrl)
	if err != nil {
		Log.Error("Channel Connection Error 1", c.routeKey, err)
		if conn != nil {
			go conn.Close()
		}
		time.Sleep(3 * time.Second)
		goto RetryConnect
	}
	channel, err := conn.Channel()
	if err != nil {
		Log.Error("Channel Connection Error 2", c.routeKey, err)
		go conn.Close()
		time.Sleep(3 * time.Second)
		goto RetryConnect
	}
	/*
		err = channel.ExchangeDeclare(c.exchange, "topic", true, false, false, false, nil)
		if err != nil {
			Log.Error("Channel Connection Error 3", c.routeKey, err)
			conn.Close()
			time.Sleep(3 * time.Second)
			goto RetryConnect
		}
	*/
BreakFor:
	for {
		select {
		case msg := <-c.Channel:
			switch msg.(type) {
			case string:
				if msg.(string) == "quit" {
					Log.Info("Channel Connection Quit", c.routeKey)
					go conn.Close()
					return
				} //quit
			case amqp.Publishing:
				err := channel.Publish(c.exchange, c.routeKey, false, false, msg.(amqp.Publishing))
				if err != nil {
					c.Channel <- msg
					go conn.Close()
					Log.Error("Channel Connection Error 4", c.routeKey, err)
					break BreakFor
				}
			case *publishingExt:
				err := channel.Publish(c.exchange, msg.(*publishingExt).Key(c.routeKey), false, false, msg.(*publishingExt).Msg())
				if err != nil {
					c.Channel <- msg
					go conn.Close()
					Log.Error("Channel Connection Error 4", c.routeKey, err)
					break BreakFor
				}
			}
		}
	}
	time.Sleep(3 * time.Second)
	goto RetryConnect
}
