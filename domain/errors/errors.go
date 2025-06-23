package errors

const (
	ServerInvalidPort   = AppError("server: invalid port number")
	ServerMissingPrefix = AppError("server: missing prefix")
	AppUserExists       = AppError("app: user already exists")
	AppUserNotFound     = AppError("app: user not found")
	PostgresMissingHost = AppError("postgres: missing host")
	PostgresMissingPort = AppError("postgres: missing port")
	PostgresMissingDB   = AppError("postgres: missing database")
	PostgresMissingUser = AppError("postgres: missing username")
	PostgresMissingPwd  = AppError("postgres: missing password")
)

type AppError string

func (e AppError) Error() string {
	return string(e)
}
