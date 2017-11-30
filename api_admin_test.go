package rmqtool

import "testing"

func TestAPIClientClusterName(t *testing.T) {
	ret, err := GenerateTestClient().ClusterName()
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(ret)
	}
}

func TestAPIClientChangeClusterName(t *testing.T) {
	params := map[string]interface{}{
		"name": "rabbit@test-rabbit-changed",
	}
	err := GenerateTestClient().ChangeClusterName(params)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("ChangeClusterName Success")
	}
}
