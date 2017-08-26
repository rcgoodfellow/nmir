package models

import (
	"github.com/deter-project/testbed/core/nmir"
)

func Tutorial_small() *nmir.Net {

	n := nmir.NewNet()

	a := n.Node().Set(nmir.Props{
		"name": "a",
		"software": nmir.Props{
			"os=": "ubuntu16.04",
		},
	})

	b := n.Node().Set(nmir.Props{
		"name": "b",
		"software": nmir.Props{
			"os=": "debian9",
		},
	})

	c := n.Node().Set(nmir.Props{
		"name": "c",
		"software": nmir.Props{
			"os=": "fedora25",
		},
	})

	s := n.Node().Set(nmir.Props{
		"name": "s",
	})

	d := n.Node().Set(nmir.Props{
		"name": "d",
		"software": nmir.Props{
			"os=": "freebsd11",
		},
	})

	n.Link(a.Endpoint(), s.Endpoint())
	n.Link(b.Endpoint(), s.Endpoint())
	n.Link(c.Endpoint(), s.Endpoint())
	n.Link(d.Endpoint(), s.Endpoint())

	return n
}
