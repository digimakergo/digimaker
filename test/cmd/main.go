package main

import (
	"github.com/xc/digimaker/core"
	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/util"
	"github.com/xc/digimaker/test"
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
