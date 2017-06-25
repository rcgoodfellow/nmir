package nmir

import (
	//"github.com/ajstarks/svgo"
	"log"
	"math"
	//"sort"
	//"os"
)

func PTree(net *Net) *Pinode {

	n := float64(len(net.Nodes))

	for i, x := range net.Nodes {
		theta := 2 * math.Pi / n * float64(i)

		x.Props["position"] = &Point{
			InitRadius * math.Cos(theta),
			InitRadius * math.Sin(theta),
		}
	}

	ptr := &Pinode{
		Width:    100000.0,
		Centroid: Point{0, 0},
	}

	for _, x := range net.Nodes {
		//log.Printf("insert %d", i)
		ptr.Insert(&Plnode{Data: x})
	}

	return ptr

}

func playout(net *Net) {

	log.Printf("building ptree")
	ptr := PTree(net)
	//log.Printf("%#v", ptr)
	log.Printf("qtree: size=%d height=%d",
		PSize(ptr),
		PHeight(ptr),
	)

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

func PNetSvg(name string, net *Net) error {
	log.Printf("net: size=%d", len(net.Nodes))

	for i := 0; i < 20; i++ {

		playout(net)

	}

	/*
		out, err := os.Create(name + ".svg")
		if err != nil {
			return err
		}
		defer out.Close()

		b := QBounds(qtr)
		cx, cy := int(b.Width/2.0), int(b.Height/2.0)

		canvas := svg.New(out)
		canvas.Start(int(b.Width), int(b.Height))

		defer canvas.End()

		doNetSvgLink(net, cx, cy, canvas)
		doNetSvgNode(net, cx, cy, canvas)
	*/

	return nil

}
