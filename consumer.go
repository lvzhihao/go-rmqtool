package rmqtool

import (
	"strings"
	"time"

	"github.com/lvzhihao/goutils"
	"github.com/streadway/amqp"
)

var (
	DefaultConsumerRetryTime time.Duration = 3 * time.Second
	DefaultConsumerToolName  string        = "golang.rmqtool"
)

func GenerateConsumerName(name string) string {
	return strings.Join([]string{
		name,
		goutils.RandStr(20),
		time.Now().Format(time.RFC3339),
	}, ".")
}

type ConsumerTool struct {
	amqpUrl   string
	conn      *amqp.Connection
	name      string
	RetryTime time.Duration
	isClosed  bool
}

func NewConsumerTool(url string) (*ConsumerTool, error) {
	c := &ConsumerTool{
		amqpUrl:   url,                      //rmq link
		RetryTime: DefaultConsumerRetryTime, //default retry
		isClosed:  false,
		name:      DefaultConsumerToolName,
	}
	// first test dial
	_, err := amqp.Dial(url)
	return c, err
}

func (c *ConsumerTool) Name() string {
	return c.name
}

func (c *ConsumerTool) SetName(name string) string {
	c.name = name
	return c.Name()
}

func (c *ConsumerTool) Link(queue string, prefetchCount int) (<-chan amqp.Delivery, error) {
	var err error
	c.conn, err = amqp.Dial(c.amqpUrl)
	if err != nil {
		return nil, err
	}
	channel, err := c.conn.Channel()
	if err != nil {
		c.conn.Close()
		return nil, err
	}
	if err := channel.Qos(prefetchCount, 0, false); err != nil {
		c.conn.Close()
		return nil, err
	}
	deliveries, err := channel.Consume(queue, GenerateConsumerName(c.name), false, false, false, false, nil)
	if err != nil {
		c.conn.Close()
		return nil, err
	}
	return deliveries, nil
}

func (c *ConsumerTool) Close() {
	// close
	if c.isClosed == false {
		if c.conn != nil {
			c.conn.Close()
		}
		c.isClosed = true
	}
}

func (c *ConsumerTool) Consume(queue string, prefetchCount int, handle func(amqp.Delivery)) {
	defer c.Close()
	for {
		if c.isClosed == true {
			Log.Error("Consumer Link Closed, Quit...", queue)
			break
		}
		time.Sleep(c.RetryTime)
		deliveries, err := c.Link(queue, prefetchCount)
		if err != nil {
			Log.Error("Consumer Link Error", err)
			continue
		}
		for msg := range deliveries {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						Log.Error("Consumer Recover", r)
					}
				}()
				handle(msg)
			}()
		}
		c.conn.Close()
		Log.Debug("Consumer ReConnection After RetryTime", c.RetryTime)
	}
}
