//Author xc, Created on 2019-03-27 20:00
//{COPYRIGHTS}

package util

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

func Contains(strings []string, element string) bool {
	for _, s := range strings {
		if s == element {
			return true
		}
	}
	return false
}
