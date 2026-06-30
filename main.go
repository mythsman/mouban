package main

import (
	"mouban/internal/app"

	"github.com/sirupsen/logrus"
)

// @title           Mouban API
// @version         1.0
// @description     Mouban 后端接口文档
// @BasePath        /
// @schemes         http https

func main() {
	if err := app.Bootstrap(); err != nil {
		logrus.Fatalln("bootstrap failed:", err)
	}

	if err := app.RunHTTPServer(); err != nil {
		logrus.Fatalln("server stopped:", err)
	}
}
