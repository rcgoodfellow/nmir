/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 * Deter 2.0 nmir
 * ==================
 *   This file defines the network modeling intermediate representation (nmir)
 *   data structures. nmir is a simple network represenatation where
 *
 *   The primary components are:
 *     - (sub)networks
 *     - nodes
 *		 - links
 *
 *   Everything is extensible through a property map member called Props.
 *
 *   Interconnection model supports node-neighbor traversal as well as
 *   link-endpoint traversal.
 *
 *   Endpoints are the glue that bind nodes to links. Everything is also
 *   upwards traversable. You can follow pointers from an endpoint to a
 *   parent node, and then to a parent network, and then to another parent
 *   network etc...
 *
 *   Serialization to json is inherent. It does not include the traversal
 *   mechanisims as this would create recursively repeditive output.
 *
 *	Copyright the Deter Project 2017 - All Rights Reserved
 *	License: Apache 2.0
 *~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

package nmir

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"sort"
)

// Data structures ------------------------------------------------------------

/*

guest : {
	"name": "gx",
	"props": {
		"hardware": {
			"memory=": "4G",
		},
		"software": [
			{ "name=": "debian-stable" }
		]
	}
}

host: {
	"name": "hx",
	"props": {
		"hardware": {
			"memory>": "32G",
			"arch=": "x86_64"
		}
	}
}

software: [
	{
		"kind": "os",
		"name": "debian-stable",
		"requirements": {
			"arch?": ["x86_64", "x86"]
		}
	}
]

*/

type Props map[string]interface{}
type Prop struct {
	Key   string
	Value interface{}
}
type SortedProps []Prop

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
	Parent    *Net        `json:"-"`
	Props     Props       `json:"props"`
	Visited   bool        `json:"-"`
}

type Link struct {
	Id        string      `json:"id"`
	Endpoints []*Endpoint `json:"endpoints"`
	Props     Props       `json:"props"`
}

type Endpoint struct {
	Id        string               `json:"id"`
	Props     Props                `json:"props"`
	Neighbors map[string]*Neighbor `json:"-"`
	Parent    *Node                `json:"-"`
}

type Software struct {
	Props        Props `json:"props"`
	Requirements Props `json:"target"`
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
		Id:    uuid.NewV4().String(),
		Props: make(Props),
		//Hardware: make(Props),
		Parent: n,
	}
	n.Nodes = append(n.Nodes, node)
	return node
}

func (n *Net) Link(es ...*Endpoint) *Link {
	link := &Link{
		Id:        uuid.NewV4().String(),
		Props:     make(Props),
		Endpoints: es,
	}
	setNeighbors(link)
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

func (n *Net) GetNodeEndpointById(id string) *Endpoint {
	for _, x := range n.Nodes {
		for _, e := range x.Endpoints {
			if e.Id == id {
				return e
			}
		}
	}
	return nil
}

// Convinience methods ~~~

func (n *Net) String() string {

	s := "nodes\n"
	s += "-----\n"

	for _, n := range n.Nodes {
		s += propString(n.Props)
		s += "\n"
	}

	s += "links\n"
	s += "-----\n"
	for _, l := range n.Links {
		s += propString(l.Props)

		s += "  endpoints: "
		for _, e := range l.Endpoints {
			//if the node has a name print that, otherwise print the id
			name, ok := e.Parent.Props["name"]
			if ok {
				s += name.(string) + " "
			} else {
				s += e.Id
			}
		}
		s += "\n\n"
	}

	return s

}

func FromFile(filename string) (*Net, error) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	net := NewNet()
	err = json.Unmarshal(data, net)
	if err != nil {
		return nil, err
	}

	linkNetwork(net)

	return net, nil

}

func (n *Net) ToFile(filename string) error {

	js, err := json.MarshalIndent(*n, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, js, 0644)
	if err != nil {
		return err
	}

	return nil

}

func (n *Net) Json() ([]byte, error) {
	js, err := json.MarshalIndent(*n, "", "  ")
	if err != nil {
		return nil, err
	}
	return js, nil
}

// When we read a network from a file or database, the traversal pointers are
// not linked. This function ensures that all the traversal pointers in a network
// data structure are linked.
func linkNetwork(net *Net) {

	//recurse networks
	for _, n := range net.Nets {
		n.Parent = net
		linkNetwork(net)
	}

	for _, n := range net.Nodes {
		n.Parent = net

		for _, e := range n.Endpoints {
			e.Parent = n
			e.Neighbors = make(map[string]*Neighbor)
		}
	}

	for _, l := range net.Links {
		for i, e := range l.Endpoints {
			e_ := net.GetNodeEndpointById(e.Id)
			l.Endpoints[i] = e_
		}
		setNeighbors(l)
	}

}

func ToFile(net *Net, filename string) error {

	js, err := json.MarshalIndent(net, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, js, 0660)

}

func (x SortedProps) Len() int           { return len(x) }
func (x SortedProps) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x SortedProps) Less(i, j int) bool { return x[i].Key < x[j].Key }

func sortProps(p Props) SortedProps {

	result := make(SortedProps, 0, len(p))
	for k, v := range p {
		result = append(result, Prop{k, v})
	}
	sort.Sort(result)
	return result

}

func propString(p Props) string {
	st := ""
	ps := sortProps(p)
	for _, x := range ps {
		st += fmt.Sprintf("  %s: %v\n", x.Key, x.Value)
	}
	return st
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

func (n *Node) AddSoftware(p Props) *Node {
	sw, ok := n.Props["software"]
	if !ok {
		n.Props["software"] = []Props{p}
		return n
	}
	_, ok = sw.([]Props)
	if !ok {
		n.Props["software"] = []Props{p}
		return n
	}

	n.Props["software"] = append(
		n.Props["software"].([]Props),
		p,
	)
	return n

	//n.Software = append(n.Software, p)
	return n
}

func (n *Node) Valence() int {
	v := 0
	for _, e := range n.Endpoints {
		v += len(e.Neighbors)
	}
	return v
}

func (n *Node) Neighbors() []*Node {

	var result []*Node
	for _, e := range n.Endpoints {
		for _, n := range e.Neighbors {
			result = append(result, n.Endpoint.Parent)
		}
	}
	return result

}

func (n *Node) Label() string {
	label := n.Id
	name, ok := n.Props["name"]
	if ok {
		label = name.(string)
	}
	return label
}

// Link methods ---------------------------------------------------------------

func (l *Link) IsLocal() bool {

	for _, x := range l.Endpoints {
		for _, y := range l.Endpoints {
			if x.Parent.Parent.Id != y.Parent.Parent.Id {
				return false
			}
		}
	}

	return true

}

func (l *Link) Set(p Props) *Link {
	for k, v := range p {
		l.Props[k] = v
	}
	return l
}

// Endpoint methods -----------------------------------------------------------

func (e *Endpoint) Set(p Props) *Endpoint {
	for k, v := range p {
		e.Props[k] = v
	}
	return e
}

// helpers ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func setNeighbors(link *Link) {
	for _, x := range link.Endpoints {
		for _, y := range link.Endpoints {
			if x == y {
				continue
			}
			x.Neighbors[y.Id] = &Neighbor{link, y}
			y.Neighbors[x.Id] = &Neighbor{link, x}
		}
	}

}
