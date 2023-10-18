package wrapper

import "net/http"

const secret = "Cr4p0nT0pD54dljn4D"
const endpoint = "http://127.0.0.1:80"

type PrivateWrapper struct {
	client http.Client
}

type CreateUserResponse struct {
	Data []struct {
		Balance            int    `json:"balance"`
		ID                 string `json:"id"`
		SolvedHcaptcha     int    `json:"solved_hcaptcha"`
		ThreadMaxHcaptcha  int    `json:"thread_max_hcaptcha"`
		ThreadUsedHcaptcha int    `json:"thread_used_hcaptcha"`
	} `json:"data"`
	Success bool `json:"success"`
}

type GetUserResponse struct {
	Data struct {
		Balance            int    `json:"balance"`
		ID                 string `json:"id"`
		SolvedHcaptcha     int    `json:"solved_hcaptcha"`
		ThreadMaxHcaptcha  int    `json:"thread_max_hcaptcha"`
		ThreadUsedHcaptcha int    `json:"thread_used_hcaptcha"`
	} `json:"data"`
	Success bool `json:"success"`
}

type AddBalanceResponse struct {
	Success bool `json:"success"`
}
