package main

import (
	"Server/cc"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options
var zoo cc.Zoo

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}

	t, err := cc.NewTurtle(conn)
	zoo.Set(t.ID, t, cc.Mining)
	log.Println("saved turtle under", t.ID)
	//defer conn.Close()

}

func main() {
	zoo = cc.MakeZoo()

	go cc.ManualControl(&zoo)

	r := mux.NewRouter()
	r.HandleFunc("/socket", socketHandler)
	r.HandleFunc("/turtle", HandleTurtleList)
	r.HandleFunc("/turtle/{id}", HandleTurtleDetails)
	r.HandleFunc("/turtle/{id}/{cmd}", HandleTurtleCommand)
	//http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
