package main

import (
	"context"
	"os"

	"github.com/alecthomas/kingpin"
	log "github.com/sirupsen/logrus"
	"gitlab.com/sdce/exlib/config"
	mgo "gitlab.com/sdce/exlib/mongo"
)

var (
	version string = "0.0.1"
)

func main() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})

	showversion := kingpin.Flag("version", "Show version information.").Short('v').Bool()
	validatorURL := kingpin.Flag("validator", "validator Url.").Default("tcp://localhost:4040").Short('d').String()
	kingpin.Parse()

	if *showversion == true {
		log.Infof("Version: %s\n", version)
		os.Exit(0)
	}
	if validatorURL == nil {
		log.Fatalln("no validator url is defined")
	}

	//load config from exlib
	v, err := config.LoadConfig("service.currency")
	if err != nil {
		log.Fatalln("Failed to load config file ", ".env")
	}
	//init repository
	mcfg, err := mgo.GetConfig(v)
	if err != nil {
		log.Fatalln("Failed to load repository config.")
	}

	ctx := context.Background()
	db := mgo.Connect(ctx, mcfg)
	if db == nil {
		log.Fatalln(err.Error())
	}
	defer db.Close(ctx)
	service := NewService(db)

	service.Run()
}
