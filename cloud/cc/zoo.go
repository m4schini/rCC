package cc

var zoo = make(map[string]Worker)

type TurtleType rune

const (
	Transport TurtleType = 'T'
	Mining               = 'M'
	Attack               = 'A'
	Farming              = 'F'
)

type Zoo struct {
	Workers map[string]*Worker
}

type Worker struct {
	Turtle *Turtle
	Type   TurtleType
}

func MakeZoo() Zoo {
	return Zoo{
		Workers: make(map[string]*Worker),
	}
}

func (zoo Zoo) Set(key string, turtle *Turtle, typ TurtleType) {
	//zoo.Workers[key].Turtle.Disconnect()
	zoo.Workers[key] = &Worker{
		Turtle: turtle,
		Type:   typ,
	}
}

func (zoo Zoo) Get(key string) (*Worker, bool) {
	t, ok := zoo.Workers[key]

	return t, ok
}
