package rest

import (
	"context"
	"dm/core"
	"dm/core/util/debug"
	_ "dm/eth/entity"
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {

	core.Bootstrap("dm/eth")
	fmt.Println("Test starting..")
	debug.Init(context.Background())
	m.Run()
}

func TestHtmlToPDF(t *testing.T) {
	// htmlToPDF("<div>Test1</div>")
}
