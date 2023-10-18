package wrapper

import "net/http"

const secret = "Cr4p0nT0pD54dljn4D"
const endpoint = "http://127.0.0.1:80"

type PrivateWrapper struct {
	client http.Client
}

type CreateUserResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

type GetUserResponse struct {
	Data struct {
		ID      string `json:"id"`
		Balance string `json:"balance"`
	} `json:"data"`
}

type AddBalanceResponse struct {
	Success bool `json:"success"`
}
