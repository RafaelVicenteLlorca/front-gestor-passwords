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

var PasswordsServiceRequest = &PasswordsService{HttpClient: client.HttpClient}

type PasswordsService struct {
	HttpClient *client.HTTPClientCustom
}

type PasswordsCreateRequest struct {
	Content string `json:"content"`
}

type PasswordsCreateResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	UpdatedAt string `json:"updated_at"`
}

func (ps *PasswordsService) Create(pcr PasswordsCreateRequest) (*PasswordsCreateResponse, error) {
	bodyReq := bytes.NewBuffer(request.JSON(pcr))
	url := ps.HttpClient.BackendUri + "passwords"
	req, _ := http.NewRequest(http.MethodPost, url, bodyReq)
	resp, err := ps.HttpClient.Do(req)
	if err != nil || resp.StatusCode != 201 {
		return &PasswordsCreateResponse{}, errors.New("error al crear la contraseña")
	}
	body, err := ioutil.ReadAll(resp.Body)
	pcres := PasswordsCreateResponse{}
	if err != nil {
		fmt.Println(err)
		return &PasswordsCreateResponse{}, err
	}

	defer resp.Body.Close()

	err = json.Unmarshal(body, &pcres)
	if err != nil {
		fmt.Println(err)
		return &PasswordsCreateResponse{}, err
	}
	return &pcres, nil
}
