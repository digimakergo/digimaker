package handler

import (
	"fmt"
	"testing"

	"github.com/xc/digimaker/core/util"
)

func TestCanLogin(t *testing.T) {
	// htmlToPDF("<div>Test1</div>")
	hash, err := util.HashPassword("Test123")
	fmt.Println(hash)
	fmt.Println(err)
	result, _ := CanLogin("admin", "test")
	fmt.Println(result)
	// CanLogin("admin", "test")
}
