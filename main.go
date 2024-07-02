package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/yaonkey/gokws/core"
)

var Environment map[string]string

func main() {
	Environment, err := godotenv.Read()
	if err != nil {
		logrus.Fatal("loading dot env file is failed")
	}
	app := core.NewEngine()

	logrus.Fatal(app.App.Listen(":" + Environment["PORT"]))
}
