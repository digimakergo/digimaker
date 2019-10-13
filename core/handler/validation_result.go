package handler

type ValidationResult struct {
	Message []string          `json:"message"`
	Fields  map[string]string `json:"fields"` //1 means required, other message means real messgae.
}

func (v ValidationResult) Passed() bool {
	return len(v.Message) == 0 && len(v.Fields) == 0
}
