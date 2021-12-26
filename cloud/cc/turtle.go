package cc

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
)

type Instruction int

const (
	Ping Instruction = iota
	Forward
	Back
	Down
	Up
	TurnLeft
	TurnRight
	Dig
	DigUp
	DigDown
	Place
	PlaceUp
	PlaceDown
	Drop
	DropUp
	DropDown
	InspectAll   = 100
	Scan         = 101
	DetectAll    = 102
	GetPosition  = 103
	GetHeading   = 104
	GetInventory = 105
)

type Vector struct {
	X int64
	Y int64
	Z int64
}

type Turtle struct {
	ID         string
	connection *websocket.Conn
	lastUpdate time.Time
	position   Vector
	fuelLevel  int
	fuelLimit  int
	inventory  []Item
}

type response struct {
	Data []interface{}
	Head map[string]interface{}
}

func (t *Turtle) handleHead(h map[string]interface{}) {

	fuel, ok := h["fuel"]
	if ok {
		t.fuelLevel = int(fuel.(float64))
	}

	pos, ok := h["position"]
	if ok {
		pos := pos.(map[string]interface{})
		x, xOK := pos["x"]
		y, yOK := pos["y"]
		z, zOK := pos["z"]

		if xOK && yOK && zOK {
			t.position.X = int64(x.(float64))
			t.position.Y = int64(y.(float64))
			t.position.Z = int64(z.(float64))
			log.Println("updated pos", t.position)
		} else {
			log.Println("WARNING: coordinates dirty")
		}
	} else {
		log.Println("position not ok")
	}
}

func (t Turtle) handleResponse(raw []byte) response {
	var tmp map[string]interface{}
	if err := json.Unmarshal(raw, &tmp); err != nil {
		log.Println(err)
	}

	switch tmp["data"].(type) {
	case []interface{}:
		return response{
			Data: tmp["data"].([]interface{}),
			Head: tmp["head"].(map[string]interface{}), //tmp["head"].(map[string]interface{}),
		}
	default:
		arr := make([]interface{}, 1)
		arr[0] = tmp["data"].(interface{})

		return response{
			Data: arr,
			Head: tmp["head"].(map[string]interface{}),
		}
	}

}

func (res response) Return() (bool, error) {
	ok := res.Data[0].(bool)

	var err error
	if !ok {
		err = errors.New(res.Data[1].(string))
	} else {
		err = nil
	}

	return ok, err
}

func (t *Turtle) execute(instruction Instruction) response {
	msg, _ := json.Marshal(struct {
		CMD Instruction
	}{CMD: instruction})

	t.connection.WriteMessage(
		websocket.TextMessage,
		msg)

	_, response, err := t.connection.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}

	t.lastUpdate = time.Now()
	log.Println("TURTLE RESPONSE:", string(response))
	r := t.handleResponse(response)
	t.handleHead(r.Head)
	return r
}

func (t *Turtle) Forward() (bool, error) {
	res := t.execute(Forward)
	return res.Return()
}

func (t *Turtle) Back() (bool, error) {
	res := t.execute(Back)
	return res.Return()
}

func (t *Turtle) Up() (bool, error) {
	res := t.execute(Up)
	return res.Return()
}

func (t *Turtle) Down() (bool, error) {
	res := t.execute(Down)
	return res.Return()
}

func (t *Turtle) TurnLeft() (bool, error) {
	res := t.execute(TurnLeft)
	return res.Return()
}

func (t *Turtle) TurnRight() (bool, error) {
	res := t.execute(TurnRight)
	return res.Return()
}

func (t *Turtle) Dig() (bool, error) {
	res := t.execute(Dig)
	return res.Return()
}

func (t *Turtle) DigUp() (bool, error) {
	res := t.execute(DigUp)
	return res.Return()
}

func (t *Turtle) DigDown() (bool, error) {
	res := t.execute(DigDown)
	return res.Return()
}

func (t *Turtle) Place() (bool, error) {
	res := t.execute(Place)
	return res.Return()
}

func (t *Turtle) PlaceUp() (bool, error) {
	res := t.execute(PlaceUp)
	return res.Return()
}

func (t *Turtle) PlaceDown() (bool, error) {
	res := t.execute(PlaceDown)
	return res.Return()
}

func (t *Turtle) Drop() (bool, error) {
	res := t.execute(Drop)
	return res.Return()
}

func (t *Turtle) DropUp() (bool, error) {
	res := t.execute(DropUp)
	return res.Return()
}

func (t *Turtle) DropDown() (bool, error) {
	res := t.execute(DropDown)

	return res.Return()
}

func (t *Turtle) GetFuelLevel() (int, time.Time) {
	return t.fuelLevel, t.lastUpdate
}

func (t *Turtle) GetInventory() []Item {
	res := t.execute(GetInventory)

	var inv = make([]Item, 16)

	for i, slot := range res.Data {
		slot := slot.(map[string]interface{})

		count, countOk := slot["count"]
		name, nameOk := slot["name"]

		if countOk && nameOk {
			inv[i] = Item{
				Name:  name.(string),
				Count: int8(count.(float64)),
			}
		}
	}
	return inv
}

func (t *Turtle) InspectAll() map[string]interface{} {
	res := t.execute(InspectAll)
	return res.Data[0].(map[string]interface{})
}

func (t *Turtle) Scan() map[string]interface{} {
	res := t.execute(Scan)
	return res.Data[0].(map[string]interface{})
}

func (t *Turtle) DetectAll() map[string]interface{} {
	res := t.execute(DetectAll)
	return res.Data[0].(map[string]interface{})
}

func (t *Turtle) UpdatePosition() []interface{} {
	res := t.execute(GetPosition)
	return res.Data
}

func (t *Turtle) GetPosition() (Vector, time.Time) {
	return t.position, t.lastUpdate
}

func (t *Turtle) GetHeading() []interface{} {
	res := t.execute(GetHeading)
	return res.Data
}

func (t *Turtle) Ping() map[string]interface{} {
	startOfReq := time.Now()

	msg, _ := json.Marshal(struct {
		CMD Instruction
	}{CMD: Ping})

	t.connection.WriteMessage(
		websocket.TextMessage,
		msg)

	_, response, _ := t.connection.ReadMessage()

	var tmp map[string]interface{}
	if err := json.Unmarshal(response, &tmp); err != nil {
		log.Println(err)
	}

	var data map[string]interface{}
	data = tmp["data"].(map[string]interface{})

	t.handleHead(tmp["head"].(map[string]interface{}))

	data["ping"] = time.Since(startOfReq)
	return data
}

func (t *Turtle) Disconnect() {
	t.connection.Close()
}

func NewTurtle(conn *websocket.Conn) (*Turtle, error) {
	t := Turtle{
		connection: conn,
	}

	data := t.Ping()
	log.Println("HANDSHAKE:", data)

	switch data["id"].(type) {
	case float64:
		t.ID = strconv.FormatFloat(data["id"].(float64), 'f', 0, 64)
	default:
		t.ID = data["id"].(string)
	}

	return &t, nil
}
