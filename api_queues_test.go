package rmqtool

import (
	"testing"
)

func TestAPIListQueues(t *testing.T) {
	ret, err := GenerateTestClient().ListQueues("")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(ret)
	}
}

func TestAPIQueues(t *testing.T) {
	testName := "test.queue"
	err := GenerateTestClient().CreateQueue("/", testName, map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("Create Test Exchange Success: ", testName)
	}
	queue, err := GenerateTestClient().Queue("/", testName)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(queue)
	}
	ret, err := GenerateTestClient().QueueBindings("/", queue["name"].(string))
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(ret)
	}
	err = GenerateTestClient().DeleteQueue("/", queue["name"].(string))
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("Delete Test Queue Success")
	}
}

func TestAPIBindings(t *testing.T) {
	testName := "test.queue.binding"
	err := GenerateTestClient().CreateQueue("/", testName, nil)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("Create Test Exchange Success: ", testName)
	}
	ret, err := GenerateTestClient().QueueCreateBinding("/", "amq.topic", testName, "aaab", nil)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("Binding Key Success: ", ret)
	}
}