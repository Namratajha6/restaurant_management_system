package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"new_restaurant/database"
	"new_restaurant/servers"
	"os"
)

func main() {
	if err := database.ConnectAndMigrate(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		"disable"); err != nil {
		logrus.Panicf("Failed to initialize and migrate database with error: %+v", err)
	}
	logrus.Print("migration successful!!")

	r := server.SetupRoutes()

	log.Println("Server running on http://localhost:8005")
	if err := http.ListenAndServe(":8005", r); err != nil {
		log.Fatal(err)
	}
}
