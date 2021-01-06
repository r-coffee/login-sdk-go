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
		return errors.New(string(body))
	}

	remoteErr := reflect.ValueOf(res).Elem().FieldByName("Error").String()

	if len(remoteErr) > 0 {
		return errors.New(remoteErr)
	}

	return nil
}

func doAuthRequired(req, res interface{}, url, token, entityID string) error {
	// construct protobuf message
	raw, err := proto.Marshal(req.(proto.Message))
	if err != nil {
		return err
	}

	// construct http request
	client := &http.Client{}
	reqq, err := http.NewRequest("POST", url, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	reqq.Header.Add("Content-Type", "application/x-protobuf")
	reqq.Header.Add("Authorization", "Bearer "+token)
	reqq.Header.Add("EntityID", entityID)
	resp, err := client.Do(reqq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// parse protobuf response
	err = proto.Unmarshal(body, res.(proto.Message))
	if err != nil {
		return errors.New(string(body))
	}

	remoteErr := reflect.ValueOf(res).Elem().FieldByName("Error").String()

	if len(remoteErr) > 0 {
		return errors.New(remoteErr)
	}

	return nil
}

func register(e, i, p, h string) error {
	var req RegisterRequest
	var res GenericResponse

	req.EntityID = e
	req.Id = i
	req.Password = p
	return do(&req, &res, fmt.Sprintf("%s/v1/register", h))
}

func login(e, i, p, h string) (string, bool, error) {
	var req LoginRequest
	var res LoginResponse

	req.EntityID = e
	req.Id = i
	req.Password = p

	err := do(&req, &res, fmt.Sprintf("%s/v1/login", h))
	if err != nil {
		return "", false, err
	}

	return res.Token, res.Admin, nil
}

func validate(e, t, h string) (bool, string, string, error) {
	var req ValidateRequest
	var res ValidateResponse

	req.EntityID = e
	req.Token = t

	err := do(&req, &res, fmt.Sprintf("%s/v1/validate", h))
	return res.Admin, res.Site, res.Id, err
}

func upgrade(e, t, h, u string) error {
	var req GenericRequest
	var res GenericResponse

	req.UserID = u

	return doAuthRequired(&req, &res, fmt.Sprintf("%s/v1/upgrade", h), t, e)
}

func downgrade(e, t, h, u string) error {
	var req GenericRequest
	var res GenericResponse

	req.UserID = u

	return doAuthRequired(&req, &res, fmt.Sprintf("%s/v1/downgrade", h), t, e)
}

func lock(e, t, h, u string) error {
	var req GenericRequest
	var res GenericResponse

	req.UserID = u

	return doAuthRequired(&req, &res, fmt.Sprintf("%s/v1/lock", h), t, e)
}

func unlock(e, t, h, u string) error {
	var req GenericRequest
	var res GenericResponse

	req.UserID = u

	return doAuthRequired(&req, &res, fmt.Sprintf("%s/v1/unlock", h), t, e)
}
