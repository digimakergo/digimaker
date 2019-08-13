package dm

type HanlderError struct {
	Code    string
	Message string
}

func (e HanlderError) Error() string {
	return "[" + e.Code + "]" + e.Message
}
