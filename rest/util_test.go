package rest

import (
	"fmt"
	"testing"

	"github.com/xc/digimaker/core"
	_ "github.com/xc/digimaker/core/entity"
)

func TestMain(m *testing.M) {

	core.Bootstrap("github.com/xc/digimaker/eth")
	fmt.Println("Test starting..")
	m.Run()
}

func TestHtmlToPDF(t *testing.T) {
	// htmlToPDF("<div>Test1</div>")
}
