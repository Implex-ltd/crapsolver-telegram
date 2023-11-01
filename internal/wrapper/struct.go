package wrapper

import "net/http"

const secret = "Cr4p0nT0pD54sdljn4D"
const endpoint = "http://127.0.0.1:123"

type PrivateWrapper struct {
	client http.Client
}

type User struct {
	Balance               int    `json:"balance"`
	ID                    string `json:"id"`
	SolvedHcaptcha        int    `json:"solved_hcaptcha"`
	ThreadMaxHcaptcha     int    `json:"thread_max_hcaptcha"`
	ThreadUsedHcaptcha    int    `json:"thread_used_hcaptcha"`
	BypassRestrictedSites bool   `json:"bypass_restricted_sites"`
}

type CreateUserResponse struct {
	Data    []User `json:"data"`
	Success bool   `json:"success"`
}

type GetUserResponse struct {
	Data    User `json:"data"`
	Success bool `json:"success"`
}

type AddBalanceResponse struct {
	//Data    User `json:"data"`
	Success bool `json:"success"`
}
