package rmqtool

import (
	"fmt"

	"github.com/streadway/amqp"
)

type ConnectConfig struct {
	Host     string                 `json:"host" yaml:"host"`         //127.0.0.1:5672
	Api      string                 `json:"api" yaml:"api"`           //http://127.0.0.1:15672
	User     string                 `json:"user" yaml:"user"`         //username
	Passwd   string                 `json:"passwd" yaml:"passwd"`     //passwd
	Vhost    string                 `json:"vhost" yaml:"vhost"`       //vhost
	MetaData map[string]interface{} `json:"metadata" yaml:"metadata"` //metadata
}

func (c *ConnectConfig) Scheme() string {
	return fmt.Sprintf("amqp://%s:%s@%s/%s", c.User, c.Passwd, c.Host, c.Vhost)
}

func Conn(config ConnectConfig) *Connect {
	return &Connect{
		config: config,
	}
}

type Connect struct {
	config ConnectConfig //queue config
}

//todo check link

func (c *Connect) Scheme() string {
	return c.config.Scheme()
}

func (c *Connect) Api() string {
	return c.config.Api
}

func (c *Connect) User() string {
	return c.config.User
}

func (c *Connect) Passwd() string {
	return c.config.Passwd
}

func (c *Connect) Vhost() string {
	return c.config.Vhost
}

func (c *Connect) Clone() *Connect {
	return &Connect{
		config: c.config,
	}
}

func (c *Connect) Dial() (*amqp.Connection, error) {
	return amqp.Dial(c.Scheme())
}

func (c *Connect) CreateExchange(exchange string) error {
	return APICreateExchange(c.Api(), c.User(), c.Passwd(), c.Vhost(), exchange, nil)
}

func (c *Connect) ApplyQueue(name string) *Queue {
	return &Queue{
		conn: c,
		name: name,
	}
}

func (c *Connect) ApplyPublisher(exchange string, routeKeys []string) (*PublisherTool, error) {
	return NewPublisherTool(c.Scheme(), exchange, routeKeys)
}
