package handler

type ValidationResult struct {
	Message string
	Fields  []FieldValidationResult
}

type FieldValidationResult struct {
	Identifier string
	Detail     string //1 means required, other message means real messgae.
}
