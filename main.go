package main

import (
	"log"
	"net/http"
)

func main() {
	server := newServer()
	server.router.ServeFiles("/static/*filepath", http.Dir(server.resourcePath))
	server.router.GET("/", server.index)
	server.router.GET("/habits", server.Get)
	server.router.PUT("/habits", server.Put)
	server.router.POST("/habits", server.Post)
	server.router.DELETE("/habits", server.Delete)

	log.Fatal(http.ListenAndServe(":8080", server.router))
}
