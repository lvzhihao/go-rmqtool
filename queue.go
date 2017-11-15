package rmqtool

import (
	"github.com/streadway/amqp"
)

type QueueConfig struct {
	Name     string       `yaml:"name"`
	Bindlist []*QueueBind `yaml:"bindlist"`
}

type QueueBind struct {
	Exchange  string                 `yaml:"exchange"`
	Key       string                 `yaml:"key"`
	Arguments map[string]interface{} `yaml:"arguments"`
}

type Queue struct {
	conn     *Connect
	name     string
	consumer *ConsumerTool
}

func (c *Queue) Clone(name string) *Queue {
	return &Queue{
		conn: c.conn,
		name: name,
	}
}

func (c *Queue) Scheme() string {
	return c.conn.Scheme()
}

func (c *Queue) Api() string {
	return c.conn.Api()
}

func (c *Queue) User() string {
	return c.conn.User()
}

func (c *Queue) Passwd() string {
	return c.conn.Passwd()
}

func (c *Queue) Vhost() string {
	return c.conn.Vhost()
}

func (c *Queue) Name() string {
	return c.name
}

func (c *Queue) ApplyConsumer() (*ConsumerTool, error) {
	return NewConsumerTool(c.Scheme())
}

func (c *Queue) Consume(prefetchCount int, handle func(amqp.Delivery)) error {
	var err error
	c.consumer, err = c.ApplyConsumer()
	if err != nil {
		return err
	}
	c.consumer.Consume(c.Name(), prefetchCount, handle)
	defer c.consumer.Close()
	return nil
}

func (c *Queue) Ensure(bindList []*QueueBind) error {
	err := c.Create()
	if err != nil {
		return err
	} else {
		return c.Bind(bindList)
	}
}

func (c *Queue) Create() error {
	return CreateQueue(
		c.Api(),
		c.User(),
		c.Passwd(),
		c.Vhost(),
		c.Name(),
	)
}

func (c *Queue) Close() {
	c.consumer.Close()
}

func (c *Queue) Bind(bindList []*QueueBind) error {
	for _, bind := range bindList {
		err := BindRoutingKey(
			c.Api(),
			c.User(),
			c.Passwd(),
			c.Vhost(),
			c.Name(),
			bind.Exchange,
			bind.Key,
			bind.Arguments,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
