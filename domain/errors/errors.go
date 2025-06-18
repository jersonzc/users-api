package errors

const (
	ServerInvalidPort   = AppError("server: invalid port number")
	ServerMissingPrefix = AppError("server: missing prefix")
)

type AppError string

func (e AppError) Error() string {
	return string(e)
}
