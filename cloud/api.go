package main

import (
	"Server/cc"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type turtleDetail struct {
	ID        string
	Type      string
	Label     string
	Position  cc.Vector
	Fuel      float64
	FuelLimit float64
	Inventory []cc.Item
}

func collectTurtleData(worker *cc.Worker) turtleDetail {
	data := (*worker).Turtle.Ping()
	inv := (*worker).Turtle.GetInventory()
	pos, _ := (*worker).Turtle.GetPosition()

	return turtleDetail{
		ID:        worker.Turtle.ID,
		Type:      string(worker.Type),
		Position:  pos,
		Fuel:      data["fuel"].(float64),
		FuelLimit: data["fuelLimit"].(float64),
		Inventory: inv,
	}
}

func HandleTurtleDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]

	worker, _ := zoo.Get(id)

	data := collectTurtleData(worker)

	_json, _ := json.Marshal(data)

	w.Header().Set("content-type", "application/json")
	w.Write(_json)
}

func HandleTurtleList(w http.ResponseWriter, r *http.Request) {
	var ts []turtleDetail

	for _, worker := range zoo.Workers {
		ts = append(ts, collectTurtleData(worker))
	}

	msg, _ := json.Marshal(ts)

	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

func toJson(t *cc.Turtle, cmd string, ok bool, err error) []byte {
	var eMsg string
	if err != nil {
		eMsg = err.Error()
	} else {
		eMsg = ""
	}

	var pos cc.Vector
	var fuel int

	if t != nil {
		pos, _ = t.GetPosition()
		fuel, _ = t.GetFuelLevel()
	}

	msg, e := json.Marshal(struct {
		Instruction string    `json:"instruction"`
		Successful  bool      `json:"success"`
		Error       string    `json:"error"`
		Position    cc.Vector `json:"position"`
		Fuel        int       `json:"fuel"`
	}{
		Instruction: cmd,
		Successful:  ok,
		Error:       eMsg,
		Position:    pos,
		Fuel:        fuel,
	})
	if e != nil {
		log.Fatalln(e)
	}

	return msg
}

func HandleTurtleCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]
	cmd, _ := vars["cmd"]
	var response []byte

	worker, foundWorker := zoo.Get(id)
	if foundWorker {
		t := worker.Turtle

		switch cmd {
		case "forward":
			ok, err := t.Forward()
			response = toJson(t, cmd, ok, err)
			break
		case "back":
			ok, err := t.Back()
			response = toJson(t, cmd, ok, err)
			break
		case "up":
			ok, err := t.Up()
			response = toJson(t, cmd, ok, err)
			break
		case "down":
			ok, err := t.Down()
			response = toJson(t, cmd, ok, err)
			break
		case "turnLeft":
			ok, err := t.TurnLeft()
			response = toJson(t, cmd, ok, err)
			break
		case "turnRight":
			ok, err := t.TurnRight()
			response = toJson(t, cmd, ok, err)
			break
		case "dig":
			ok, err := t.Dig()
			response = toJson(t, cmd, ok, err)
			break
		case "digDown":
			ok, err := t.DigDown()
			response = toJson(t, cmd, ok, err)
			break
		case "digUp":
			ok, err := t.DigUp()
			response = toJson(t, cmd, ok, err)
			break
		default:
			response = toJson(t, cmd, false, errors.New("cmd doesn't exist"))
		}
	} else {
		response = toJson(nil, cmd, false, errors.New("turtle doesn't exist"))
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	_, err := w.Write(response)
	if err != nil {
		log.Fatalln(err)
	}
}
