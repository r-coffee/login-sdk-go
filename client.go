package login

// host cannot end with /
const host = "https://login-service-296122.uc.r.appspot.com"

// Client :
type Client struct {
	EntityID string
}

func createClient(entityID string) *Client {
	var c Client
	c.EntityID = entityID
	return &c
}

func (c *Client) register(id, password string) error {
	return register(c.EntityID, id, password, host)
}

func (c *Client) login(id, password string) (string, error) {
	return login(c.EntityID, id, password, host)
}

func (c *Client) validate(id, token string) error {
	return validate(c.EntityID, token, host)
}
