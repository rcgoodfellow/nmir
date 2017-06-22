package nmir

import (
	"github.com/satori/go.uuid"
)

// Data structures ------------------------------------------------------------

type Props map[string]interface{}

type Net struct {
	Id     string  `json:"id"`
	Nodes  []*Node `json:"nodes"`
	Links  []*Link `json:"links"`
	Nets   []*Net  `json:"nets"`
	Props  Props   `json:"props"`
	Parent *Net    `json:"-"`
}

type Node struct {
	Id        string      `json:"id"`
	Endpoints []*Endpoint `json:"endpoints"`
	Props     Props       `json:"props"`
	Parent    *Net        `json:"-"`
}

type Link struct {
	Id        string         `json:"id"`
	Endpoints [2][]*Endpoint `json:"endpoints"`
	Props     Props          `json:"props"`
}

type Endpoint struct {
	Id        string               `json:"id"`
	Props     Props                `json:"props"`
	Neighbors map[string]*Neighbor `json:"-"`
	Parent    *Node                `json:"-"`
}

type Neighbor struct {
	Link     *Link
	Endpoint *Endpoint
}

// Net methods ----------------------------------------------------------------

// Factory methods ~~~

func NewNet() *Net {
	return &Net{
		Id:    uuid.NewV4().String(),
		Props: make(Props),
	}
}

func (n *Net) Net() *Net {
	return &Net{
		Id:     uuid.NewV4().String(),
		Props:  make(Props),
		Parent: n,
	}
}

func (n *Net) Node() *Node {
	node := &Node{
		Id:     uuid.NewV4().String(),
		Props:  make(Props),
		Parent: n,
	}
	n.Nodes = append(n.Nodes, node)
	return node
}

func (n *Net) Link(a, b []*Endpoint) *Link {
	link := &Link{
		Id:        uuid.NewV4().String(),
		Props:     make(Props),
		Endpoints: [2][]*Endpoint{a, b},
	}
	setNeighbors(link, a, b)
	n.Links = append(n.Links, link)
	link.Props["local"] = link.IsLocal()
	return link
}

// Search methods ~~~

func (n *Net) GetNode(uuid string) *Node {
	for _, x := range n.Nodes {
		if x.Id == uuid {
			return x
		}
		for _, e := range x.Endpoints {
			if e.Id == uuid {
				return x
			}
		}
	}
	return nil
}

func (n *Net) GetNodeByName(name string) *Node {
	for _, x := range n.Nodes {
		if x.Props["name"] == name {
			return x
		}
		for _, e := range x.Endpoints {
			if e.Props["name"] == name {
				return x
			}
		}
	}
	return nil
}

// Node methods ---------------------------------------------------------------

func (n *Node) Endpoint() *Endpoint {
	ep := &Endpoint{
		Id:        uuid.NewV4().String(),
		Props:     make(Props),
		Neighbors: make(map[string]*Neighbor),
		Parent:    n,
	}
	n.Endpoints = append(n.Endpoints, ep)
	return ep
}

func (n *Node) Set(p Props) *Node {
	for k, v := range p {
		n.Props[k] = v
	}
	return n
}

func (l *Link) Set(p Props) *Link {
	for k, v := range p {
		l.Props[k] = v
	}
	return l
}

// Link methods ---------------------------------------------------------------

func (l *Link) IsLocal() bool {

	return len(l.Endpoints[0]) == 0 ||
		len(l.Endpoints[1]) == 0 ||
		l.Endpoints[0][0].Parent.Parent.Id == l.Endpoints[1][0].Parent.Parent.Id

}

// Endpoint methods -----------------------------------------------------------

func (e *Endpoint) Set(p Props) *Endpoint {
	for k, v := range p {
		e.Props[k] = v
	}
	return e
}

// helpers ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func setNeighbors(link *Link, a, b []*Endpoint) {
	for _, x := range a {
		for _, y := range b {
			x.Neighbors[y.Id] = &Neighbor{link, y}
			y.Neighbors[x.Id] = &Neighbor{link, x}
		}
	}

}
