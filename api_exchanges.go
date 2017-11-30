package rmqtool

import (
	"net/url"
)

func (c *APIClient) ListExchanges() ([]map[string]interface{}, error) {
	return c.readSlice("exchanges")
}

func APIListExchanges(api, user, passwd string) ([]map[string]interface{}, error) {
	return NewAPIClient(api, user, passwd).ListExchanges()
}

func (c *APIClient) ListVhostExchanges(name string) ([]map[string]interface{}, error) {
	return c.readSlice([]string{"exchanges", name})
}

func APIListVhostExchanges(api, user, passwd, name string) ([]map[string]interface{}, error) {
	return NewAPIClient(api, user, passwd).ListVhostExchanges(name)
}

func (c *APIClient) Exchange(vhost, name string) (map[string]interface{}, error) {
	return c.readMap([]string{"exchanges", vhost, name})
}

func APIExchange(api, user, passwd, vhost, name string) (map[string]interface{}, error) {
	return NewAPIClient(api, user, passwd).Exchange(vhost, name)
}

func (c *APIClient) CreateExchange(vhost, name string, params map[string]interface{}) error {
	// `{"type":"topic","auto_delete":false,"durable":true,"internal":false,"arguments":[]}`
	return c.create([]string{"exchanges", vhost, name}, params)
}

func APICreateExchange(api, user, passwd, vhost, name string, params map[string]interface{}) error {
	return NewAPIClient(api, user, passwd).CreateExchange(vhost, name, params)
}

func (c *APIClient) DeleteExchange(vhost, name string) error {
	req, err := c.NewRequest("delete", []string{"exchanges", vhost, name}, nil)
	if err != nil {
		return err
	}
	query := &url.Values{}
	query.Add("if-unused", "true")
	req.URL.RawQuery = query.Encode()
	resp, err := c.Do(req)
	return c.ScanDelete(resp, err)
}

func APIDeleteExchange(api, user, passwd, vhost, name string) error {
	return NewAPIClient(api, user, passwd).DeleteExchange(vhost, name)
}

func (c *APIClient) ForceDeleteExchange(vhost, name string) error {
	return c.delete([]string{"exchanges", vhost, name})
}

func APIForceDeleteExchange(api, user, passwd, vhost, name string) error {
	return NewAPIClient(api, user, passwd).ForceDeleteExchange(vhost, name)
}

func (c *APIClient) ListExchangeSourceBindings(vhost, name string) ([]map[string]interface{}, error) {
	return c.readSlice([]string{"exchanges", vhost, name, "bindings", "source"})
}

func APIListExchangeSourceBindings(api, user, passwd, vhost, name string) ([]map[string]interface{}, error) {
	return NewAPIClient(api, user, passwd).ListExchangeSourceBindings(vhost, name)
}

func (c *APIClient) ListExchangeDestinationBindings(vhost, name string) ([]map[string]interface{}, error) {
	return c.readSlice([]string{"exchanges", vhost, name, "bindings", "destination"})
}

func APIListExchangeDestinationBindings(api, user, passwd, vhost, name string) ([]map[string]interface{}, error) {
	return NewAPIClient(api, user, passwd).ListExchangeDestinationBindings(vhost, name)
}

func (c *APIClient) ExchangePublish(vhost, name string, params map[string]interface{}) (map[string]interface{}, error) {
	req, err := c.NewRequest("post", []string{"exchanges", vhost, name, "publish"}, params)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	return c.ScanMap(resp, err)
}

func APIExchangePublish(api, user, passwd, vhost, name string, params map[string]interface{}) (map[string]interface{}, error) {
	return NewAPIClient(api, user, passwd).ExchangePublish(vhost, name, params)
}
