package nmir

import (
	"log"
)

type Pquad [4]Pnode

type Pnode interface{}

type Pinode struct {
	Quad     Pquad
	Centroid Point
	velocity Point
	Width    float64
	Mass     float64
}

type Plnode struct {
	Data Positional
}

func (pi *Pinode) Insert(x *Plnode) {
	/*
		log.Printf("newnode %f,%f",
			x.Data.Position().X,
			x.Data.Position().Y,
		)
	*/
	if x == nil {
		return
	}

	i := pi.Select(x)
	p := pi.Quad[i]
	pi.Mass++

	if p == nil {
		pi.Quad[i] = x
		return
	}

	//subdivide
	lnode, ok := p.(*Plnode)
	if ok {
		if distance(lnode.Data, x.Data) < 1e-3 {
			log.Fatalf("fuck")
		}
		new_node := pi.NewQuad(i)
		pi.Quad[i] = new_node
		new_node.Insert(lnode)
		new_node.Insert(x)
		return
	}

	//recurse
	inode, ok := p.(*Pinode)
	if ok {
		inode.Insert(x)
		return
	}

}

func (p *Pinode) NewQuad(sector int) *Pinode {
	//log.Printf("newquad %d", sector)
	result := &Pinode{
		Width: p.Width / 2.0,
	}
	shift := p.Width / 4.0
	switch sector {
	case 0:
		result.Centroid = Point{
			X: p.Centroid.X - shift,
			Y: p.Centroid.Y + shift,
		}
	case 1:
		result.Centroid = Point{
			X: p.Centroid.X + shift,
			Y: p.Centroid.Y + shift,
		}
	case 2:
		result.Centroid = Point{
			X: p.Centroid.X + shift,
			Y: p.Centroid.Y - shift,
		}
	case 3:
		result.Centroid = Point{
			X: p.Centroid.X - shift,
			Y: p.Centroid.Y - shift,
		}
	}
	return result
}

func (p *Pinode) Select(x *Plnode) int {

	if x.Data.Position().Y >= p.Centroid.Y {
		if x.Data.Position().X <= p.Centroid.X {
			return 0
		} else {
			return 1
		}
	} else {
		if x.Data.Position().X > p.Centroid.X {
			return 2
		} else {
			return 3
		}
	}

}

func (p *Pinode) Position() *Point {
	return &p.Centroid
}

func (p *Pinode) Velocity() *Point {
	return &p.velocity
}

func (p *Pinode) Weight() float64 {
	return p.Mass
}
