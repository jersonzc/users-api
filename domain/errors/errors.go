package errors

const (
	ServerInvalidPort   = AppError("server: invalid port number")
	ServerMissingPrefix = AppError("server: missing prefix")
	AppUserExists       = AppError("app: user already exists")
	AppUserNotFound     = AppError("app: user not found")
)

type AppError string

func (e AppError) Error() string {
	return string(e)
}
