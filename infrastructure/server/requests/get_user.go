package requests

type MultipleIDRequest struct {
	Users []string `json:"users" required:"true"`
}
