package nmir

import (
	"github.com/satori/go.uuid"
)

type Props map[string]interface{}

type Net struct {
	Id    string  `json:"id"`
	Nodes []*Node `json:"nodes"`
	Links []*Link `json:"links"`
	Nets  []*Net  `json:"nets"`
}

func NewNet() Net {
	return Net{
		Id: uuid.NewV4().String(),
	}
}

func (n *Net) Node() *Node {
	node := &Node{
		Id:    uuid.NewV4().String(),
		Props: make(Props),
	}
	n.Nodes = append(n.Nodes, node)
	return node
}

func (n *Net) Link() *Link {
	link := &Link{
		Id:    uuid.NewV4().String(),
		Props: make(Props),
	}
	n.Links = append(n.Links, link)
	return link
}

type Node struct {
	Id        string      `json:"id"`
	Endpoints []*Endpoint `json:"endpoints"`
	Props     Props       `json:"props"`
}

func (n *Node) Endpoint() *Endpoint {
	ep := &Endpoint{
		Id:    uuid.NewV4().String(),
		Props: make(Props),
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

type Link struct {
	Id        string         `json:"id"`
	Endpoints [2][]*Endpoint `json:"endpoints"`
	Props     Props          `json:"props"`
}

func (l *Link) Set(p Props) *Link {
	for k, v := range p {
		l.Props[k] = v
	}
	return l
}

type Endpoint struct {
	Id    string `json:"id"`
	Props Props  `json:"props"`
}

func (e *Endpoint) Set(p Props) *Endpoint {
	for k, v := range p {
		e.Props[k] = v
	}
	return e
}
