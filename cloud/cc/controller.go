package cc

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	SegmentGap = 3
	SideLength = 7
)

func ManualControl(zoo *Zoo) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		var t *Turtle

		fmt.Println("SELECT TURTLE")
		for scanner.Scan() {
			in := scanner.Text()

			v, ok := zoo.Get(in)
			if ok {
				t = v.Turtle
				break
			} else {
				fmt.Println("turtle", in, "doesnt exist!")
			}
		}

		fmt.Println("TURTLE AWAITING INPUT")
		for scanner.Scan() {
			in := scanner.Text()
			fmt.Println(">", in)

			switch in {
			case "w", "forward":
				log.Println(t.Forward())
				break
			case "s", "back":
				log.Println(t.Back())
				break
			case "d", "turnRight":
				log.Println(t.TurnRight())
				break
			case "a", "turnLeft":
				log.Println(t.TurnLeft())
				break
			case "q", "up":
				log.Println(t.Up())
				break
			case "e", "down":
				log.Println(t.Down())
				break
			case "f", "dig":
				log.Println(t.Dig())
				break
			case "v", "digDown":
				log.Println(t.DigDown())
				break
			case "r", "digUp":
				log.Println(t.DigUp())
				break
			case "p", "ping":
				data := t.Ping()
				log.Println(data)
				log.Println("PING", data["id"], "\t", data["ping"])
				break
			case "inspect":
				log.Println(t.InspectAll())
				break
			case "scan":
				log.Println(t.Scan())
				break
			case "pos", "position":
				log.Println(t.UpdatePosition())
				break
			case "heading", "face":
				log.Println(t.GetHeading())
				break
			case "inv", "inventory":
				log.Println(t.GetInventory())
				break
			case "mine":
				log.Println("begin strip mine")
				t.StripMine(10)
			}

			if in == "exit()" {
				break
			}

		}
	}
}

func (t *Turtle) StripMine(segments int) {
	for i := 0; i < 2; i++ {
		for i := 0; i < segments; i++ {
			for i := 0; i < SegmentGap+1; i++ {
				t.sliceForward()
			}

			t.TurnRight()
			for i := 0; i < SideLength; i++ {
				t.sliceForward()
			}
			t.TurnRight()
			t.TurnRight()
			for i := 0; i < SideLength; i++ {
				t.Forward()
			}
			t.TurnRight()
		}
		t.TurnLeft()
		t.sliceForward()
		t.TurnRight()
		t.Dig()
		t.TurnRight()
		t.TurnRight()
		t.Back()
	}

}

func (t *Turtle) sliceForward() {
	t.Dig()
	t.Forward()
	t.DigUp()
}
