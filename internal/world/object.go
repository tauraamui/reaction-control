package world

import "github.com/hajimehoshi/ebiten/v2"

// Objects represented in the game world
var Objects []Object

type Object interface {
	Update()
	Draw(*ebiten.Image, *ebiten.DrawImageOptions, int)
}