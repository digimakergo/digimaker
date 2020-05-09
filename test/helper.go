//Author xc, Created on 2020-05-09 18:00
//{COPYRIGHTS}

//test package provides helpers and setting up test environment for unit test. eg. set up basic content structure in db.
package test

import (
	"context"
	"fmt"

	"github.com/xc/digimaker/core"
	"github.com/xc/digimaker/core/util"
)

var RootID = 3

var started = false
var ctx context.Context

//Start bootstraps environment. Note it can be invoke multiple times because the tests can be in different package, but inside it will only init once.
func Start() context.Context {
	if !started {
		fmt.Println("Starting testing...")
		testFolder := util.DMPath() + "/test"
		core.Bootstrap(testFolder)
		ctx = context.Background()
	}
	return ctx
}

func InitData() {

}

func CleanData() {

}
