package rmqtool

import "testing"

func TestAPIListExchange(t *testing.T) {
	ret, err := GenerateTestClient().ListExchanges()
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(ret)
	}
}

func TestAPIListVhostExchange(t *testing.T) {
	ret, err := GenerateTestClient().ListVhostExchanges("/")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(ret)
	}
}

func TestAPIExhcnage(t *testing.T) {
	testName := "test.test"
	err := GenerateTestClient().CreateExchange("/", testName, map[string]interface{}{
		"type": "direct",
	})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("Create Test Exchange Success: ", testName)
	}
	exchange, err := GenerateTestClient().Exchange("/", testName)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(exchange)
	}
	err = GenerateTestClient().DeleteExchange("/", exchange["name"].(string))
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("Delete Test Exchange Success")
	}
}

func TestAPIListExchangeSourceBindings(t *testing.T) {
	ret, err := GenerateTestClient().ListExchangeSourceBindings("/", "amq.topic")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(ret)
	}
}

func TestAPIListExchangeDestinationBindings(t *testing.T) {
	ret, err := GenerateTestClient().ListExchangeDestinationBindings("/", "amq.topic")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(ret)
	}
}

func TestAPIExchangePublish(t *testing.T) {
	ret, err := GenerateTestClient().ExchangePublish("/", "amq.topic", map[string]interface{}{
		"properties":       struct{}{},
		"routing_key":      "my key",
		"payload":          "my body",
		"payload_encoding": "string",
	})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("Publish Success: ", ret)
	}
}
