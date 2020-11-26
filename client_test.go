package login

import (
	"testing"
)

const (
	eid       = "Eg4KBmVudGl0eRoER29vZA"
	id        = "foo"
	pass      = "password"
	adminPass = "}V[=jJyQnkmN63hRBHRq"
)

func TestRegister(t *testing.T) {
	client := CreateClient(eid)
	err := client.Register(id, pass)
	if err != nil {
		t.Errorf("should not return error: %v", err)
	}
}

func TestRegisterWithInvalidEID(t *testing.T) {
	client := CreateClient("invalid")
	err := client.Register(id, pass)
	if err.Error() != "invalid entity" {
		t.Error("expected error")
	}
}

func TestLogin(t *testing.T) {
	client := CreateClient(eid)
	token, admin, err := client.Login(id, pass)

	if err != nil {
		t.Errorf("expected no error: %v", err)
	}

	if admin {
		t.Error("expected normal user")
	}

	if len(token) == 0 {
		t.Error("invalid token")
	}
}

func TestLoginAdmin(t *testing.T) {
	client := CreateClient(eid)
	token, admin, err := client.Login(eid, adminPass)

	if err != nil {
		t.Errorf("expected no error: %v", err)
	}

	if admin == false {
		t.Error("expected admin user")
	}

	if len(token) == 0 {
		t.Error("invalid token")
	}
}

func TestLoginLocked(t *testing.T) {
	client := CreateClient(eid)
	token, _, err := client.Login(eid, adminPass)
	if err != nil {
		t.Errorf("login error: %v", err)
	}
	client.Lock(token, id)
	_, _, err = client.Login(id, pass)
	if err != nil && err.Error() != "user account locked" {
		t.Errorf("expected a different error: %v", err)
	}

	client.Unlock(token, id)
}

func TestValidate(t *testing.T) {
	client := CreateClient(eid)
	token, _, _ := client.Login(id, pass)
	admin, site, err := client.Validate(token)

	if admin {
		t.Error("expected normal user")
	}

	if err != nil {
		t.Errorf("expected no error: %v", err)
	}

	if site != eid {
		t.Errorf("wrong site: %s", site)
	}
}

func TestValidateBadToken(t *testing.T) {
	client := CreateClient(eid)
	_, _, err := client.Validate("invalid")

	if err == nil || err.Error() != "invalid token" {
		t.Errorf("expected differnt error: %v", err)
	}
}

func TestValidateLockedUser(t *testing.T) {
	client := CreateClient(eid)
	tokenAdmin, _, _ := client.Login(eid, adminPass)

	token, _, _ := client.Login(id, pass)
	client.Lock(tokenAdmin, id)

	_, _, err := client.Validate(token)
	if err == nil || err.Error() != "user account locked" {
		t.Errorf("expected differnt error: %v", err)
	}

	client.Unlock(tokenAdmin, id)
}

func TestUpgrade(t *testing.T) {
	client := CreateClient(eid)
	tokenAdmin, _, _ := client.Login(eid, adminPass)

	// user should start as user
	_, admin, _ := client.Login(id, pass)
	if admin {
		t.Error("expected normal user")
	}

	// upgrade user
	client.Upgrade(tokenAdmin, id)

	// user should be admin
	_, admin, _ = client.Login(id, pass)
	if !admin {
		t.Error("expected admin user")
	}

	client.Downgrade(tokenAdmin, id)
}
