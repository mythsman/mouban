package main

import (
	"mouban/internal/app"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := app.Bootstrap(); err != nil {
		logrus.Fatalln("bootstrap failed:", err)
	}

	if err := app.RunHTTPServer(); err != nil {
		logrus.Fatalln("server stopped:", err)
	}
}
