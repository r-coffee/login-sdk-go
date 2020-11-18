package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	reflect "reflect"

	proto "github.com/golang/protobuf/proto"
)

func do(req, res interface{}, endpoint string) error {
	raw, err := proto.Marshal(req.(proto.Message))
	if err != nil {
		return err
	}
	resp, err := http.Post("http://localhost:8080"+endpoint, "application/x-protobuf", bytes.NewReader(raw))
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

func register(e, i, p string) error {
	var req RegisterRequest
	var res RegisterResponse

	req.EntityID = e
	req.Id = i
	req.Password = p
	return do(&req, &res, "/v1/register")
}

func login(e, i, p string) (string, error) {
	var req LoginRequest
	var res LoginResponse

	req.EntityID = e
	req.Id = i
	req.Password = p

	err := do(&req, &res, "/v1/login")
	if err != nil {
		return "", err
	}

	return res.Token, nil
}

func validate(e, t string) error {
	var req ValidateRequest
	var res ValidateResponse

	req.EntityID = e
	req.Token = t

	return do(&req, &res, "/v1/validate")
}

func main() {
	eid := "Eg4KBmVudGl0eRoER29vZA"

	err := register(eid, "foo", "password")
	if err != nil {
		log.Println(err)
	}

	token, err := login(eid, "foo", "password")
	if err != nil {
		log.Println(err)
	}
	log.Println(token)

	err = validate(eid, token)
	if err != nil {
		log.Println(err)
	}
}
