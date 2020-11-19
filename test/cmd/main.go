package main

import (
	"github.com/digimakergo/digimaker/core"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/digimakergo/digimaker/test"
)

func main() {
	//Init data
	log.Info("Init data...")
	testFolder := util.DMPath() + "/test"
	core.Bootstrap(testFolder)

	//schema
	log.Info("Init schema...")
	test.InitSchema()

	//data
	log.Info("Init data...")
	test.InitData()

}
