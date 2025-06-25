package requests

type MultipleIDRequest struct {
	Users []int `json:"users" required:"true"`
}
