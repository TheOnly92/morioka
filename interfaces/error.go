package interfaces

type WebError struct {
	Error   error
	Message string
	Code    int
	Json    bool
}
