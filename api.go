package rmqtool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func ApiClientDo(user, passwd string, req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	req.SetBasicAuth(user, passwd)
	req.Header.Add("Content-Type", "application/json")
	return client.Do(req)
}

func CreateExchange(api, user, passwd, vhost, exchange string) error {
	b := bytes.NewBufferString(`{"type":"topic","auto_delete":false,"durable":true,"internal":false,"arguments":[]}`)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/exchanges/%s/%s", api, vhost, exchange), b)
	if err != nil {
		return err
	}
	// enusre exchange
	resp, err := ApiClientDo(user, passwd, req)
	if err != nil {
		return err
	}
	if (resp.StatusCode == http.StatusNoContent) || (resp.StatusCode == http.StatusCreated) {
		return nil
	} else {
		return fmt.Errorf("CreateExchange StatusError: %d, %v", resp.StatusCode, resp)
	}
}

func ListQueues(api, user, passwd, vhost string) ([]map[string]interface{}, error) {
	var req *http.Request
	var err error
	if vhost == "" {
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/queues", api), nil)
	} else {
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/queues/%s", api, vhost), nil)
	}
	if err != nil {
		return nil, err
	}
	resp, err := ApiClientDo(user, passwd, req)
	if err != nil {
		return nil, err
	}
	b, _ := ioutil.ReadAll(resp.Body)
	var ret []map[string]interface{}
	err = json.Unmarshal(b, &ret)
	return ret, err
}

func CreateQueue(api, user, passwd, vhost, name string) error {
	b := bytes.NewBufferString(`{"auto_delete":false, "durable":true, "arguments":[]}`)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/queues/%s/%s", api, vhost, name), b)
	if err != nil {
		return err
	}
	// enusre queue
	resp, err := ApiClientDo(user, passwd, req)
	if err != nil {
		return err
	}
	if (resp.StatusCode == http.StatusNoContent) || (resp.StatusCode == http.StatusCreated) {
		return nil
	} else {
		return fmt.Errorf("CreateQueue StatusError: %d, %v", resp.StatusCode, resp)
	}
}

func DeleteQueue(api, user, passwd, vhost, name string) error {
	b := bytes.NewBufferString(`{}`)
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/queues/%s/%s", api, vhost, name), b)
	if err != nil {
		return err
	}
	// enusre queue
	resp, err := ApiClientDo(user, passwd, req)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNoContent {
		return nil
	} else {
		return fmt.Errorf("CreateQueue StatusError: %d, %v", resp.StatusCode, resp)
	}
}

func BindRoutingKey(api, user, passwd, vhost, name, exchange, key string, args map[string]interface{}) error {
	params := map[string]interface{}{
		"routing_key": key,
		"arguments":   args,
	}
	bt, err := json.Marshal(params)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(bt)
	// ensure binding
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/bindings/%s/e/%s/q/%s", api, vhost, exchange, name),
		b,
	)
	resp, err := ApiClientDo(user, passwd, req)
	if err != nil {
		return err
	}
	if (resp.StatusCode == http.StatusNoContent) || (resp.StatusCode == http.StatusCreated) {
		return nil
	} else {
		return fmt.Errorf("BindRoutingKey StatusError: %d, %v", resp.StatusCode, resp)
	}
}

func RegisterQueue(api, user, passwd, vhost, name, exchange string, keys []string) error {
	err := CreateQueue(api, user, passwd, vhost, name)
	if err != nil {
		return err
	}
	for _, key := range keys {
		if exchange != "" && key != "" {
			err := BindRoutingKey(api, user, passwd, vhost, name, exchange, key, nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
