package main

import (
	"fmt"
	"os/exec"

	"github.com/xc/digimaker/core"
	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/util"
)

func main() {
	//Init data
	log.Info("Init data...")
	testFolder := util.DMPath() + "/test"
	core.Bootstrap(testFolder)

	dbConfig := util.GetConfigSection("database")
	user := dbConfig["username"]
	password := dbConfig["password"]
	database := dbConfig["database"]
	host := dbConfig["host"]

	cmdStr := "mysql -h " + host + " -u " + user + " -p" + password + " " + database + " < "
	dataFolder := util.DMPath() + "/data"

	//schema
	log.Info("Init schema...")
	schemaCmd := cmdStr + dataFolder + "/schema.sql"
	cmd := exec.Command("bash", "-c", schemaCmd)
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err))
		return
	}

	//data
	log.Info("Init data...")
	dataCmd := cmdStr + dataFolder + "/initdata.sql"
	cmd = exec.Command("bash", "-c", dataCmd)
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err))
		return
	}

}
