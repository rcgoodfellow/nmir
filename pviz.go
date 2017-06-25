package nmir

import (
	"github.com/ajstarks/svgo"
	"log"
	"math"
	"os"
)

const (
	BHC = 0.5
)

func pinit(net *Net) {
	n := float64(len(net.Nodes))

	for i, x := range net.Nodes {
		theta := 2 * math.Pi / n * float64(i)

		x.Props["position"] = &Point{
			InitRadius * math.Cos(theta),
			InitRadius * math.Sin(theta),
		}
		x.Props["dp"] = &Point{0, 0}
	}

	for i := 0; i < 0; i++ {
		constrain(net)
		step(net)
	}

}

func PTree(net *Net) *Pinode {

	ptr := &Pinode{
		Width:    10000.0,
		Centroid: Point{0, 0},
	}

	for _, x := range net.Nodes {
		//log.Printf("insert %d", i)
		ptr.Insert(&Plnode{Data: x})
	}

	return ptr

}

func playout(net *Net) *Pinode {

	ptr := PTree(net)

	/*
		log.Printf("qtree: size=%d height=%d",
			PSize(ptr),
			PHeight(ptr),
		)
	*/

	Pspread(net, ptr)
	step(net)
	constrain(net)
	step(net)

	return ptr

}

func PHeight(ptr *Pinode) int {

	h := 1
	hs := make([]int, 4)
	for i, x := range ptr.Quad {
		_, ok := x.(*Plnode)
		if ok {
			hs[i] = 1
		}
		p, ok := x.(*Pinode)
		if ok {
			hs[i] = PHeight(p)
		}
	}
	max := 0
	for _, x := range hs {
		if x > max {
			max = x
		}
	}
	return h + max

}

func PSize(ptr *Pinode) int {

	s := 0
	for _, x := range ptr.Quad {
		_, ok := x.(*Plnode)
		if ok {
			s += 1
		}
		p, ok := x.(*Pinode)
		if ok {
			s += PSize(p)
		}
	}
	return s

}

func Pspread(net *Net, ptr *Pinode) {

	for _, n := range net.Nodes {
		Pforce(n, ptr)
	}

}

func Pforce(node *Node, ptr *Pinode) {

	for _, p := range ptr.Quad {

		if p == nil {
			continue
		}

		lnode, ok := p.(*Plnode)
		if ok {
			pfab(lnode.Data, node)
			if node.Props["name"] == "bpc223" {
				/*
					log.Printf("push %s[%d] @%f %f,%f",
						lnode.Data.(*Node).Props["name"],
						lnode.Data.(*Node).Valence(),
						distance(lnode.Data, node),
						node.Props["dp"].(*Point).X,
						node.Props["dp"].(*Point).Y,
					)
				*/
			}
		}

		inode, ok := p.(*Pinode)
		if ok {
			s := inode.Width
			d := distance(node, inode)

			if s/d < BHC { //far enough away to aggregate
				pfab(inode, node)
				//Pforce(node, inode)
			} else { //need to recurse down futher
				Pforce(node, inode)
			}
		}

	}

}

func PBounds(ptr *Pinode) Bounds {
	ul := PUpperLeft(ptr)
	lr := PLowerRight(ptr)
	return Bounds{
		ul.Y - lr.Y,
		lr.X - ul.X,
	}
}

func PUpperLeft(ptr *Pinode) *Point {
	return PBound(ptr, 0)
}

func PLowerRight(ptr *Pinode) *Point {
	return PBound(ptr, 2)
}

func PBound(ptr *Pinode, quad int) *Point {

	x := ptr.Quad[quad]
	if x == nil {
		return &ptr.Centroid
	}

	lnode, ok := x.(*Plnode)
	if ok {
		return lnode.Data.Position()
	}

	inode, ok := x.(*Pinode)
	if x != nil {
		return PBound(inode, quad)
	}

	return &ptr.Centroid
}

func PNetSvg(name string, net *Net) error {
	log.Printf("net: size=%d", len(net.Nodes))

	//var ptr *Pinode
	pinit(net)

	for i := 0; i < 100; i++ {

		log.Printf(">>==={%d}-------->", i)
		/*ptr = */ playout(net)

	}

	out, err := os.Create(name + ".svg")
	if err != nil {
		return err
	}
	defer out.Close()

	/*ptr = PTree(net)*/
	//b := PBounds(ptr)
	b := Bounds{
		Height: 5000,
		Width:  5000,
	}
	cx, cy := int(b.Width/2.0), int(b.Height/2.0)

	canvas := svg.New(out)
	canvas.Start(int(b.Width), int(b.Height))

	defer canvas.End()

	doNetSvgLink(net, cx, cy, canvas)
	doNetSvgNode(net, cx, cy, canvas)

	return nil

}
