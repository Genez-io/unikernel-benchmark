package main

import (
	"benchmark-tool/cli"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.DateTime,
	})

	app := cli.CreateApp()
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
		return
	}
}
