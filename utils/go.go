package utils

// This is utils for enhace go language.

//IfElse implements ternary. Equvelant to <cond>?<trueV>:<falseV> in other language.
//Ref: https://github.com/golang/go/issues/20774 The reason was insane.
//
//There might be performance not perfect, since it needs type convertion.
func IfElse(cond bool, trueV interface{}, falseV interface{}) interface{} {
	result := trueV
	if !cond {
		result = falseV
	}
	return result
}
