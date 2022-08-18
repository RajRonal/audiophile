package main

import (
	"audioPhile/database"
	"audioPhile/server"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	err := database.ConnectAndMigrate(os.Getenv("host"), os.Getenv("port"), os.Getenv("databaseName"), os.Getenv("user"), os.Getenv("password"), database.SSLModeDisable)
	if err != nil {
		logrus.Fatal(err)
		return
	}
	fmt.Println("connected")
	//cmd.Execute()

	srv := server.SetupRoutes()
	err = srv.Run(":8080")
	if err != nil {
		logrus.Error(err)
	}

}
