package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"passwordsAdmin/client"
	"passwordsAdmin/pkg/request"
)

type UserService struct {
	HttpClient *client.HTTPClientCustom
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (us *UserService) Login(lr LoginRequest) (LoginResponse, error) {
	bodyReq := bytes.NewBuffer(request.JSON(lr))
	url := us.HttpClient.BackendUri + "users/login"
	req, _ := http.NewRequest(http.MethodPost, url, bodyReq)
	resp, err := us.HttpClient.Do(req)

	if err != nil || resp.Status != "200 OK" {
		return LoginResponse{}, errors.New("usuario o contraseña incorrectas")
	}
	body, err := ioutil.ReadAll(resp.Body)
	lresp := LoginResponse{}
	if err != nil {
		fmt.Println(err)
		return LoginResponse{}, err
	}

	defer resp.Body.Close()

	err = json.Unmarshal(body, &lresp)
	if err != nil {
		fmt.Println(err)
		return LoginResponse{}, err
	}
	return lresp, nil
}
