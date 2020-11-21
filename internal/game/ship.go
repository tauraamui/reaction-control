package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ship space ship
type Ship struct {
	Object
	image            *ebiten.Image
	xPos, yPos, rPos float64
	xSpd, ySpd, rSpd float64
	rmax             float64
	vmax             float64
	cwThrusters      bool
	ccwThrusters     bool
	fwdThrusters     bool
	revThrusters     bool
}

// NewShip is initialized and returned
func NewShip(x float64, y float64) *Ship {
	return &Ship{
		image: shipImage,
		xPos:  x,
		yPos:  y,
		xSpd:  0,
		ySpd:  0,
		rSpd:  0,
		rmax:  10,
		vmax:  1,
	}
}

// Update the ship state
func (ship *Ship) Update() {
	ship.xPos += ship.xSpd
	ship.yPos += ship.ySpd
	ship.rPos = math.Mod(ship.rPos+ship.rSpd, 360)

	if ship.fwdThrusters {
		radAng := (ship.rPos + 90) * (math.Pi / 180)
		xSpd := ship.xSpd - 0.01*math.Cos(radAng)
		ySpd := ship.ySpd - 0.01*math.Sin(radAng)

		if math.Abs(xSpd)+math.Abs(ySpd) < ship.vmax {
			ship.xSpd = xSpd
			ship.ySpd = ySpd
		}
	}

	if ship.revThrusters {
		radAng := (ship.rPos + 90) * (math.Pi / 180)
		xSpd := ship.xSpd + 0.01*math.Cos(radAng)
		ySpd := ship.ySpd + 0.01*math.Sin(radAng)

		if math.Abs(xSpd)+math.Abs(ySpd) < ship.vmax {
			ship.xSpd = xSpd
			ship.ySpd = ySpd
		}
	}

	if ship.cwThrusters {
		if ship.rSpd < ship.rmax {
			ship.rSpd += 0.05
		}
	}

	if ship.ccwThrusters {
		if ship.rSpd > -ship.rmax {
			ship.rSpd -= 0.05
		}
	}

	if !ship.isThrusting() && rcsSound.IsPlaying() {
		stopRcsSound()
	}
}

// Draw the ship on
func (ship *Ship) Draw(screen *ebiten.Image, g *Game) {

	op := &ebiten.DrawImageOptions{}

	imgWidth, imgHeight := ship.image.Size()
	op.GeoM.Translate(-float64(imgWidth)/2, -float64(imgHeight)/2)
	op.GeoM.Rotate(float64(ship.rPos) * 2 * math.Pi / 360)

	x := (ship.xPos - g.viewPort.xPos) + (g.viewPort.width / 2)
	y := (ship.yPos - g.viewPort.yPos) + (g.viewPort.height / 2)

	op.GeoM.Translate(x, y)
	screen.DrawImage(ship.image, op)

	frame := (g.count / 2) % 2

	if ship.ccwThrusters {
		screen.DrawImage(rcsfl.SubImage(image.Rect(frame*32, 0, 32+(frame*32), 32)).(*ebiten.Image), op)
		screen.DrawImage(rcsbr.SubImage(image.Rect(frame*32, 0, 32+(frame*32), 32)).(*ebiten.Image), op)
	}

	if ship.cwThrusters {
		screen.DrawImage(rcsfr.SubImage(image.Rect(frame*32, 0, 32+(frame*32), 32)).(*ebiten.Image), op)
		screen.DrawImage(rcsbl.SubImage(image.Rect(frame*32, 0, 32+(frame*32), 32)).(*ebiten.Image), op)
	}

	if ship.fwdThrusters {
		if !ship.cwThrusters {
			screen.DrawImage(rcsbl.SubImage(image.Rect(frame*32, 0, 32+(frame*32), 32)).(*ebiten.Image), op)
		}

		if !ship.ccwThrusters {
			screen.DrawImage(rcsbr.SubImage(image.Rect(frame*32, 0, 32+(frame*32), 32)).(*ebiten.Image), op)
		}
	}

	if ship.revThrusters {
		if !ship.ccwThrusters {
			screen.DrawImage(rcsfl.SubImage(image.Rect(frame*32, 0, 32+(frame*32), 32)).(*ebiten.Image), op)
		}

		if !ship.cwThrusters {
			screen.DrawImage(rcsfr.SubImage(image.Rect(frame*32, 0, 32+(frame*32), 32)).(*ebiten.Image), op)
		}
	}
}

// CwThrustersOn clockwise thrusters on
func (ship *Ship) CwThrustersOn() {
	if !ship.cwThrusters {
		ship.cwThrusters = true
		startRcsSound()
	}
}

// CwThrustersOff clockwise thruters off
func (ship *Ship) CwThrustersOff() {
	ship.cwThrusters = false
}

// CcwThrustersOn counter clockwise thrusters on
func (ship *Ship) CcwThrustersOn() {
	if !ship.ccwThrusters {
		ship.ccwThrusters = true
		startRcsSound()
	}
}

// CcwThrustersOff counter clockwise thrusters off
func (ship *Ship) CcwThrustersOff() {
	ship.ccwThrusters = false
}

// FwdThrustersOn forward thrusters on
func (ship *Ship) FwdThrustersOn() {
	if !ship.fwdThrusters {
		ship.fwdThrusters = true
		startRcsSound()
	}
}

// FwdThrustersOff forward thrusters off
func (ship *Ship) FwdThrustersOff() {
	ship.fwdThrusters = false
}

// RevThrustersOn reverse thrusters on
func (ship *Ship) RevThrustersOn() {
	if !ship.revThrusters {
		ship.revThrusters = true
		startRcsSound()
	}
}

// RevThrustersOff reverse thrusters off
func (ship *Ship) RevThrustersOff() {
	ship.revThrusters = false
}

func (ship *Ship) isThrusting() bool {
	return ship.fwdThrusters || ship.revThrusters || ship.cwThrusters || ship.ccwThrusters
}