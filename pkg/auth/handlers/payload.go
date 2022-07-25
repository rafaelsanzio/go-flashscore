package handlers

type APIKeyPayload struct {
	Email string `json:"email"`
	Key   string
}
