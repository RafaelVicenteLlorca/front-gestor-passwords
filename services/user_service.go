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

var UserServiceRequest = &UserService{HttpClient: client.HttpClient}

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

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func (us *UserService) Login(lr LoginRequest) (LoginResponse, error) {
	bodyReq := bytes.NewBuffer(request.JSON(lr))
	url := us.HttpClient.BackendUri + "users/login"
	req, _ := http.NewRequest(http.MethodPost, url, bodyReq)
	resp, err := us.HttpClient.Do(req)

	if err != nil || resp.StatusCode != 200 {
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

func (us *UserService) Register(lr RegisterRequest) (RegisterResponse, error) {
	bodyReq := bytes.NewBuffer(request.JSON(lr))
	url := us.HttpClient.BackendUri + "users"
	req, _ := http.NewRequest(http.MethodPost, url, bodyReq)
	resp, err := us.HttpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return RegisterResponse{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return RegisterResponse{}, err
	}

	if err != nil || resp.StatusCode != 201 {
		errorMessage := ErrorMessage{}
		err = json.Unmarshal(body, &errorMessage)
		if err != nil {
			return RegisterResponse{}, errors.New("error al crear el usuario")
		}
		return RegisterResponse{}, errors.New(errorMessage.Message)
	}
	regResp := RegisterResponse{}

	defer resp.Body.Close()

	err = json.Unmarshal(body, &regResp)
	if err != nil {
		fmt.Println(err)
		return RegisterResponse{}, err
	}
	return regResp, nil
}
