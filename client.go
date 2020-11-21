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

// Login : log in user and return jwt token and Admin status
func (c *Client) Login(id, password string) (string, bool, error) {
	return login(c.EntityID, id, password, host)
}

// Validate : validate jwt token. returns Admin status
func (c *Client) Validate(token string) (bool, error) {
	return validate(c.EntityID, token, host)
}

// Upgrade : upgrade user to Admin
func (c *Client) Upgrade(token, userID string) error {
	return upgrade(c.EntityID, token, host, userID)
}

// Downgrade : downgrade user (remove Admin rights)
func (c *Client) Downgrade(token, userID string) error {
	return downgrade(c.EntityID, token, host, userID)
}

// Lock : locks user. user cannot login or validate a token
func (c *Client) Lock(token, userID string) error {
	return lock(c.EntityID, token, host, userID)
}

// Unlock : unlocks a user
func (c *Client) Unlock(token, userID string) error {
	return unlock(c.EntityID, token, host, userID)
}
