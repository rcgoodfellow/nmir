/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 * Deter 2.0 nmir-viz
 * ==================
 *   This file implements a PR type quadtree that is used for network
 *   visualization. The primary motivation for using this data structure is
 *   that it allows us to compute a force directed layout in l*log(n) time
 *   using the Barnes-Hut algorithm, as opposed to the typical n^2 approach
 *
 *   The quadtree is composed of two node types - interior (Pinode) and leaf
 *   (Pleaf). Leaf nodes contain data and interior nodes are purely for data
 *   organization. A good reference on this data structure is 'Foundations of
 *   Multidimensional and Metric Data Structures' by Hanan Samet - see section
 *   1.4.2.2.
 *
 *	Copyright the Deter Project 2017 - All Rights Reserved
 *	License: Apache 2.0
 *~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

package viz

import (
	"github.com/deter-project/testbed/core/nmir"
)

const (
	BHC = 0.5 //Barnes-Hut Constant
)

// A Pnode is a generic node type that is always either a Pinode or Plnode
type Pnode interface{}

type Pinode struct {
	Quad     [4]Pnode   // The four quadrants this interior node presides over
	Centroid nmir.Point // Center of mass for this interior node
	velocity nmir.Point // Used in force calculations
	Width    float64    // The collective width of the underlying 4 quadrants
	Mass     float64    // The collective mass of the underlying 4 quadrants
}

type Plnode struct {
	Data Positional // The data that is encapsulated by this leaf node
}

// PTree constructs a new PTree from an nmir network. The return value is
// the root of the tree. This is an O(n*log(n)) operation where n is the
// number of nodes in the provided network.
func PTree(net *nmir.Net) *Pinode {

	// compute the width of the root internal node based on the bounds of the
	// network
	var b Bounds
	bounds(net, &b, 0, 0)

	width := b.Width
	if b.Height > b.Width {
		width = b.Height
	}

	ptr := &Pinode{
		Width:    width,
		Centroid: nmir.Point{0, 0},
	}

	for _, x := range net.Nodes {
		ptr.Insert(&Plnode{Data: x})
	}

	return ptr

}

// This method inserts a new Pleaf node below this Pinode. This is a recursive
// insertion with O(log(n)) complexity where n is the number of nodes currently
// below this Pinode.
func (pi *Pinode) Insert(x *Plnode) {

	// safeguard
	if x == nil {
		return
	}

	pi.Mass++

	// figure out which quad to insert into
	i := pi.Select(x)
	p := pi.Quad[i]

	// if the selected quad is empty, we have found a new home of the leaf and
	// we are done
	if p == nil {
		pi.Quad[i] = x
		return
	}

	// if the selected quad is a leaf node, then we subdivide that leaf node
	// by replacing it with a new interior node and then inserting both the
	// old and new leaf into that interior node
	lnode, ok := p.(*Plnode)
	if ok {
		new_node := pi.NewQuad(i)
		pi.Quad[i] = new_node
		new_node.Insert(lnode)
		new_node.Insert(x)
		return
	}

	// if the selected quad is an interior node, then we recurse into that node
	// and continue on
	inode, ok := p.(*Pinode)
	if ok {
		inode.Insert(x)
		return
	}

}

// This method determines which quadrant a new leaf node node should fall in.
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

// This method creates a new quad for the specified sector of this Pinode.
func (p *Pinode) NewQuad(sector int) *Pinode {

	// the width of any new quad is half that of its parent
	result := &Pinode{
		Width: p.Width / 2.0,
	}

	// calculate the centroid of the new quad
	shift := p.Width / 4.0
	switch sector {
	case 0:
		result.Centroid = nmir.Point{
			X: p.Centroid.X - shift,
			Y: p.Centroid.Y + shift,
		}
	case 1:
		result.Centroid = nmir.Point{
			X: p.Centroid.X + shift,
			Y: p.Centroid.Y + shift,
		}
	case 2:
		result.Centroid = nmir.Point{
			X: p.Centroid.X + shift,
			Y: p.Centroid.Y - shift,
		}
	case 3:
		result.Centroid = nmir.Point{
			X: p.Centroid.X - shift,
			Y: p.Centroid.Y - shift,
		}
	}
	return result
}

// Compute the height and width bounds for this internal node
func (ptr *Pinode) Bounds() Bounds {
	ul := ptr.Limit(0)
	lr := ptr.Limit(2)
	return Bounds{
		ul.Y - lr.Y,
		lr.X - ul.X,
	}
}

// Compute the extreme corner value of a given quadrant
func (ptr *Pinode) Limit(quad int) *nmir.Point {

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
		return inode.Limit(quad)
	}

	return &ptr.Centroid
}

// Compute the height of this interior node
func (ptr *Pinode) Height() int {

	h := 1
	hs := make([]int, 4)
	for i, x := range ptr.Quad {
		_, ok := x.(*Plnode)
		if ok {
			hs[i] = 1
		}
		p, ok := x.(*Pinode)
		if ok {
			hs[i] = p.Height()
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

// Compute the total number of leaf nodes presiding within this Pinode
func (ptr *Pinode) Size() int {

	s := 0
	for _, x := range ptr.Quad {
		_, ok := x.(*Plnode)
		if ok {
			s += 1
		}
		p, ok := x.(*Pinode)
		if ok {
			s += p.Size()
		}
	}
	return s

}

// Compute the force on the specified node within this Pinode
func (ptr *Pinode) Force(node Positional) {

	for _, p := range ptr.Quad {

		if p == nil {
			continue
		}

		lnode, ok := p.(*Plnode)
		if ok {
			fab(lnode.Data, node)
		}

		inode, ok := p.(*Pinode)
		if ok {
			s := inode.Width
			d := distance(node, inode)

			if s/d < BHC { //far enough away to aggregate
				fab(inode, node)
			} else { //need to recurse down futher
				inode.Force(node)
			}
		}

	}

}

// Compute the forces on all nodes in the network within this Pinode
func (ptr *Pinode) Forces(net *nmir.Net) {

	for _, n := range net.Nodes {
		ptr.Force(n)
	}

}

// Implementation of the Positional interface for interior nodes
func (p *Pinode) Position() *nmir.Point {
	return &p.Centroid
}

func (p *Pinode) Velocity() *nmir.Point {
	return &p.velocity
}

func (p *Pinode) Weight() float64 {
	return p.Mass
}
