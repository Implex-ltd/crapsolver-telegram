package wrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (w *PrivateWrapper) CreateUser() (string, error) {
	req, err := http.NewRequest("POST", endpoint+"/api/user/internal/new", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", secret)

	resp, err := w.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response CreateUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response.Data) == 0 {
		return "", fmt.Errorf("No user ID returned")
	}

	return response.Data[0].ID, nil
}

func (w *PrivateWrapper) GetUser(user string) (*GetUserResponse, error) {
	resp, err := w.client.Get(endpoint + "/api/user/" + user)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response GetUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	fmt.Println(response)

	return &response, nil
}

func (w *PrivateWrapper) AddBalance(user string, amount int) (*AddBalanceResponse, error) {
	data := map[string]interface{}{
		"amount": amount,
		"user":   user,
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint+"/api/user/internal/refill", bytes.NewBuffer(dataJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", secret)
	req.Header.Add("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response AddBalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
