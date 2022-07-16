package main

import (
	"fmt"
	"math/rand"

	"github.com/costaconrado/mew"
)

type GridPos struct {
	X, Y int
}

type PlayerController struct {
	filter mew.Filter
}

func (system *PlayerController) Init() {
	system.filter = mew.NewFilter(GridPos{})
}

func (system *PlayerController) Update() {
	for _, entity := range system.filter.Entities() {
		pos := (*GridPos)(mew.GetComponent[GridPos](entity))
		fmt.Printf("\n\nold pos value [%d,%d]", pos.X, pos.Y)
		pos.X = rand.Intn(5)
		pos.Y = rand.Intn(5)
		fmt.Printf("\nnew pos value [%d,%d]", pos.X, pos.Y)
	}
}

func main() {
	mainLayer := mew.NewLayer()
	mew.AddLayer(mainLayer)

	mainLayer.Add(&PlayerController{})

	player := mew.NewEntity()

	g := mew.AddComponent[GridPos](player)
	g.X = 5
	g.Y = 7

	mew.Init()

	fmt.Println("starting")
	for true {
		mew.Update()
	}
}
