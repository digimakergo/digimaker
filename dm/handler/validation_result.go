package handler

type ValidationResult struct {
	Message []string
	Fields  []FieldValidationResult
}

func (v ValidationResult) Passed() bool {
	return len(v.Message) == 0 && len(v.Fields) == 0
}

type FieldValidationResult struct {
	Identifier string
	Detail     string //1 means required, other message means real messgae.
}
