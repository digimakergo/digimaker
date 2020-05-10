//Author xc, Created on 2020-05-09 18:00
//{COPYRIGHTS}

//test package provides helpers and setting up test environment for unit test. eg. set up basic content structure in db.
package test

import (
	"context"
	"fmt"
	"os/exec"

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
		InitData()
		testFolder := util.DMPath() + "/test"
		core.Bootstrap(testFolder)
		ctx = context.Background()
	}
	return ctx
}

func InitSchema() {
	runSQL("schema.sql")
}

func InitData() {
	runSQL("initdata.sql")
}

func runSQL(file string) {
	dbConfig := util.GetConfigSection("database")
	user := dbConfig["username"]
	password := dbConfig["password"]
	database := dbConfig["database"]
	host := dbConfig["host"]

	dataFolder := util.DMPath() + "/data"
	cmdStr := "mysql -h " + host + " -u " + user + " -p" + password + " " + database + " < "

	dataCmd := cmdStr + dataFolder + "/" + file
	cmd := exec.Command("bash", "-c", dataCmd)
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err))
	}
}
