package main

import (
	"log"
	"mexemexe/internal/server"
	"mexemexe/internal/service"
	"net/http"
)

func main() {

	serverConfig := server.NewServerConfig(service.LEVEL_DEBUG)

	server := server.NewServer(serverConfig)
	http.HandleFunc("/ws", server.HandleConnections)

	log.Println("HTTP server started on :8888")
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
