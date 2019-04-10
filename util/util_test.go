package util

import (
	"fmt"
	"testing"
)

func TestReadConfig(t *testing.T) {
	database := GetConfigSection("database", "site")

	t.Log(database["host"])

	database = GetConfigSection("database", "site")

	t.Log(database["host"], database["database"], database["username"])

	id := GenerateUID()
	fmt.Println(id)
}
