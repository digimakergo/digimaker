package rest

import (
	"dm/core"
	_ "dm/eth/entity"
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {

	core.Bootstrap("dm/eth")
	fmt.Println("Test starting..")
	m.Run()
}

func TestHtmlToPDF(t *testing.T) {
	// htmlToPDF("<div>Test1</div>")
}
