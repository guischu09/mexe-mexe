package client

// import (
// 	"flag"
// 	"fmt"
// 	"log"
// 	"mexemexe/internal/server"
// 	"net/url"

// 	"github.com/gorilla/websocket"
// )

// var addr = flag.String("addr", "localhost:8888", "http service address")

// func main() {

// 	flag.Parse()
// 	log.SetFlags(0)

// 	var username string
// 	fmt.Println("Lets play mexe-mexe!")
// 	fmt.Println("To join a game room, please enter your username:")
// 	fmt.Scanf("%s", &username)

// 	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
// 	log.Printf("connecting to %s\n", u.String())

// 	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
// 	if err != nil {
// 		log.Fatal("dial:", err)
// 	}
// 	defer ws.Close()

// 	joinMessage := server.JoinServerMessage{
// 		Username: username,
// 	}
// 	err = ws.WriteJSON(joinMessage)
// 	if err != nil {
// 		log.Printf("error writing to websocket: %v", err)
// 		return
// 	}

// 	var msg string
// 	err = ws.ReadJSON(&msg)
// 	if err != nil {
// 		log.Printf("error: %v", err)
// 		return
// 	}
// 	fmt.Println(msg)

// }
