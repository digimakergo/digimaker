package handler

type ValidationResult struct {
	Message []string                `json:"message"`
	Fields  []FieldValidationResult `json:"fields"`
}

func (v ValidationResult) Passed() bool {
	return len(v.Message) == 0 && len(v.Fields) == 0
}

type FieldValidationResult struct {
	Identifier string `json:"identifier"`
	Detail     string `json:"detail"` //1 means required, other message means real messgae.
}
