package user

import (
	"context"
	"dm/core"
	"dm/core/util/debug"
	"fmt"
	"testing"

	_ "dm/eth/entity"
)

func TestMain(m *testing.M) {

	core.Bootstrap("dm/eth")
	fmt.Println("Test starting..")
	debug.Init(context.Background())
	m.Run()
}

func TestCanLogin(t *testing.T) {
	// htmlToPDF("<div>Test1</div>")
	hash, err := HashPassword("test")
	fmt.Println(hash)
	fmt.Println(err)
	result, _ := CanLogin("admin", "test")
	fmt.Println(result)
	// CanLogin("admin", "test")
}
