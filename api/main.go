package main

import (
	"fmt"
	"github.com/joho/godotenv"
	//"gotemplate/database"
	"gotemplate/middleware"
	"gotemplate/routes"
	"log"
	"net/http"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//err = database.ConnectToDatabase()
	//if err != nil {
	//log.Fatal("Connection to DB failed")
	//}
}

func main() {
	router := http.NewServeMux()

	// handle static resources
	staticFs := http.FileServer(http.Dir("./static"))
	router.Handle("GET /static/", http.StripPrefix("/static", staticFs))

	// routes
	pageRouter := http.NewServeMux()
	routes.SetHomePageHandlers(pageRouter)

	stack := middleware.CreateMiddlewareStack(
		middleware.Logging,
	)
	router.Handle("/", stack(pageRouter))

	// serve
	port := os.Getenv("PORT")
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}
	fmt.Println("Server listening on port", port)
	server.ListenAndServe()
}
