package fieldtype

import (
	"fmt"
	"testing"
)

func TestDatetime(t *testing.T) {
	handler := DatetimeHandler{}
	value, err := handler.LoadInput("2021-10-12", "new")
	fmt.Println(value, err)

	value, err = handler.LoadInput("2021-10-12 01:10:15", "new")
	fmt.Println(value, err)
}
