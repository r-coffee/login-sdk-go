package login

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	reflect "reflect"

	proto "github.com/golang/protobuf/proto"
)

func do(req, res interface{}, url string) error {
	raw, err := proto.Marshal(req.(proto.Message))
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/x-protobuf", bytes.NewReader(raw))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(body, res.(proto.Message))
	if err != nil {
		return err
	}

	remoteErr := reflect.ValueOf(res).Elem().FieldByName("Error").String()

	if len(remoteErr) > 0 {
		return errors.New(remoteErr)
	}

	return nil
}

func register(e, i, p, h string) error {
	var req RegisterRequest
	var res RegisterResponse

	req.EntityID = e
	req.Id = i
	req.Password = p
	return do(&req, &res, fmt.Sprintf("%s/v1/register", h))
}

func login(e, i, p, h string) (string, error) {
	var req LoginRequest
	var res LoginResponse

	req.EntityID = e
	req.Id = i
	req.Password = p

	err := do(&req, &res, fmt.Sprintf("%s/v1/login", h))
	if err != nil {
		return "", err
	}

	return res.Token, nil
}

func validate(e, t, h string) error {
	var req ValidateRequest
	var res ValidateResponse

	req.EntityID = e
	req.Token = t

	return do(&req, &res, fmt.Sprintf("%s/v1/validate", h))
}
