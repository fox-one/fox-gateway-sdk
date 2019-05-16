package gateway

func (client *Client) ServiceP(service string) *Client {
	return client.Group("/p/" + service)
}
