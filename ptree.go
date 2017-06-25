package nmir

import (
//"log"
)

type Pquad [4]Pnode

type Pnode interface{}

type Pinode struct {
	Quad     Pquad
	Centroid Point
	Width    float64
	Weight   float64
}

type Plnode struct {
	Data Positional
}

func (pi *Pinode) Insert(x *Plnode) {
	//log.Printf("insert %v", x.Data.Position())
	//log.Printf("centroid %v", pi.Centroid)
	if x == nil {
		return
	}

	i := pi.Select(x)
	//log.Printf("sector %d", i)
	p := pi.Quad[i]

	if p == nil {
		pi.Quad[i] = x
		return
	}

	//subdivide
	lnode, ok := p.(*Plnode)
	if ok {
		//log.Printf("subd")
		new_node := pi.NewQuad(i)
		pi.Quad[i] = new_node

		//log.Printf("ins lnode")
		new_node.Insert(lnode)

		//log.Printf("ins x")
		new_node.Insert(x)
		//log.Printf("KERPOO")
		return
	}

	//recurse
	inode, ok := p.(*Pinode)
	if ok {
		//log.Printf("rec")
		inode.Insert(x)
		return
	}

}

func (p *Pinode) NewQuad(sector int) *Pinode {
	result := &Pinode{
		Width: p.Width / 2.0,
	}
	shift := p.Width / 4.0
	switch sector {
	case 0:
		result.Centroid = Point{
			X: p.Centroid.X - shift,
			Y: p.Centroid.X + shift,
		}
	case 1:
		result.Centroid = Point{
			X: p.Centroid.X + shift,
			Y: p.Centroid.X + shift,
		}
	case 2:
		result.Centroid = Point{
			X: p.Centroid.X + shift,
			Y: p.Centroid.X - shift,
		}
	case 3:
		result.Centroid = Point{
			X: p.Centroid.X - shift,
			Y: p.Centroid.X - shift,
		}
	}
	return result
}

func (p *Pinode) Select(x *Plnode) int {

	if x.Data.Position().Y <= p.Centroid.Y {
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
