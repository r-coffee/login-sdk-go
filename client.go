package login

// host cannot end with /
const host = "https://login-service-296122.uc.r.appspot.com"

// Client :
type Client struct {
	EntityID string
}

// CreateClient : create a new login client
func CreateClient(entityID string) *Client {
	var c Client
	c.EntityID = entityID
	return &c
}

// Register : register a new user
func (c *Client) Register(id, password string) error {
	return register(c.EntityID, id, password, host)
}

// Login : log in user and get jwt
func (c *Client) Login(id, password string) (string, error) {
	return login(c.EntityID, id, password, host)
}

// Validate : validate jwt token
func (c *Client) Validate(id, token string) error {
	return validate(c.EntityID, token, host)
}
