package handler

import (
	"context"
	"github.com/xc/digimaker/core"
	"github.com/xc/digimaker/core/util"
	"github.com/xc/digimaker/core/util/debug"
	"fmt"
	"testing"

	_ "github.com/xc/digimaker/eth/entity"
)

func TestMain(m *testing.M) {

	core.Bootstrap("github.com/xc/digimaker/eth")
	fmt.Println("Test starting..")
	debug.Init(context.Background())
	m.Run()
}

func TestCanLogin(t *testing.T) {
	// htmlToPDF("<div>Test1</div>")
	hash, err := util.HashPassword("Test123")
	fmt.Println(hash)
	fmt.Println(err)
	result, _ := CanLogin("admin", "test")
	fmt.Println(result)
	// CanLogin("admin", "test")
}
