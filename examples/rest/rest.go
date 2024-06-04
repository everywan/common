package main

import "github.com/everywan/common/application"

func main() {
	app := application.New()
	// e := gin.New()

	// httpBundle := rest.New()
	// httpBundle.Run()
	app.AddBundle()
}
