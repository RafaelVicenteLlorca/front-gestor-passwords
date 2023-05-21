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

	"github.com/TwiN/go-color"
)

var PasswordsServiceRequest = &PasswordsService{HttpClient: client.HttpClient}

type PasswordsService struct {
	HttpClient *client.HTTPClientCustom
}

type PasswordsCreateRequest struct {
	Content string `json:"content"`
}

type PasswordsResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	UpdatedAt string `json:"updated_at"`
}

type ErrorPasswordMessage struct {
	Message string `json:"message"`
}

func (ps *PasswordsService) Create(pcr PasswordsCreateRequest) (*PasswordsResponse, error) {
	bodyReq := bytes.NewBuffer(request.JSON(pcr))
	url := ps.HttpClient.BackendUri + "passwords"
	req, _ := http.NewRequest(http.MethodPost, url, bodyReq)
	resp, err := ps.HttpClient.Do(req)
	if err != nil || resp.StatusCode != 201 {
		return &PasswordsResponse{}, errors.New("error al crear la contraseña")
	}
	body, err := ioutil.ReadAll(resp.Body)
	pcres := PasswordsResponse{}
	if err != nil {
		fmt.Println(err)
		return &PasswordsResponse{}, err
	}

	defer resp.Body.Close()

	err = json.Unmarshal(body, &pcres)
	if err != nil {
		fmt.Println(err)
		return &PasswordsResponse{}, err
	}
	return &pcres, nil
}

func (ps *PasswordsService) GetById(id string) (PasswordsResponse, error) {
	url := ps.HttpClient.BackendUri + "passwords/" + id
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := ps.HttpClient.Do(req)
	if err != nil {
		fmt.Println(color.Colorize(color.Red, err.Error()))
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(color.Colorize(color.Red, err.Error()))
		return PasswordsResponse{}, err
	}

	if resp.StatusCode != 200 {
		errorMessage := ErrorPasswordMessage{}
		err = json.Unmarshal(body, &errorMessage)
		if err != nil {
			fmt.Println(err)
			return PasswordsResponse{}, err
		}
		return PasswordsResponse{}, errors.New(errorMessage.Message)
	}
	pcres := PasswordsResponse{}

	err = json.Unmarshal(body, &pcres)
	if err != nil {
		fmt.Println(err)
		return PasswordsResponse{}, err
	}
	return pcres, nil
}

func (ps *PasswordsService) GetAll() ([]PasswordsResponse, error) {
	url := ps.HttpClient.BackendUri + "passwords"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := ps.HttpClient.Do(req)
	if err != nil {
		fmt.Println(color.Colorize(color.Red, err.Error()))
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(color.Colorize(color.Red, err.Error()))
		return []PasswordsResponse{}, err
	}

	if resp.StatusCode != 200 {
		errorMessage := ErrorPasswordMessage{}
		err = json.Unmarshal(body, &errorMessage)
		if err != nil {
			fmt.Println(err)
			return []PasswordsResponse{}, err
		}
		return []PasswordsResponse{}, errors.New(errorMessage.Message)
	}
	pcres := []PasswordsResponse{}

	err = json.Unmarshal(body, &pcres)
	if err != nil {
		fmt.Println(err)
		return []PasswordsResponse{}, err
	}
	return pcres, nil
}
