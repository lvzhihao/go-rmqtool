package rmqtool

func (c *APIClient) ClusterName() (map[string]interface{}, error) {
	return c.readMap("cluster-name")
}

func (c *APIClient) ChangeClusterName(params map[string]interface{}) error {
	return c.create("/cluster-name", params)
}

func APIClusterName(api, user, passwd string) (map[string]interface{}, error) {
	return NewAPIClient(api, user, passwd).ClusterName()
}

func APIChnageClusterName(api, user, passwd string, params map[string]interface{}) error {
	return NewAPIClient(api, user, passwd).ChangeClusterName(params)
}
