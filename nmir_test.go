package nmir

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func net4() *Net {

	host := NewNet()
	zwitch := host.Node().Set(Props{"name": "leaf"})
	for i := 0; i < 4; i++ {
		zwitch.Endpoint().Set(Props{"bandwidth": "1G"})
	}

	a := host.Node().Set(Props{"name": "a"})
	a.Endpoint().Set(Props{"bandwidth": "1G"})

	b := host.Node().Set(Props{"name": "b"})
	b.Endpoint().Set(Props{"bandwidth": "1G"})

	c := host.Node().Set(Props{"name": "c"})
	c.Endpoint().Set(Props{"bandwidth": "1G"})

	d := host.Node().Set(Props{"name": "d"})
	d.Endpoint().Set(Props{"bandwidth": "1G"})

	host.Link(
		[]*Endpoint{zwitch.Endpoints[0]},
		[]*Endpoint{a.Endpoints[0]},
	)

	host.Link(
		[]*Endpoint{zwitch.Endpoints[1]},
		[]*Endpoint{b.Endpoints[0]},
	)

	host.Link(
		[]*Endpoint{zwitch.Endpoints[2]},
		[]*Endpoint{c.Endpoints[0]},
	)

	host.Link(
		[]*Endpoint{zwitch.Endpoints[3]},
		[]*Endpoint{d.Endpoints[0]},
	)

	return host

}

func TestModelA(t *testing.T) {

	a := net4()

	buf, _ := json.MarshalIndent(a, "", "  ")
	ioutil.WriteFile("4net.json", buf, 0644)

	VTag(a)

	buf, _ = json.MarshalIndent(a, "", "  ")
	ioutil.WriteFile("4net_vt.json", buf, 0644)

}

func TestModelAB(t *testing.T) {

	a := net4()
	b := net4()

	ab := NewNet()
	ab.Nets = append(ab.Nets, a)
	ab.Nets = append(ab.Nets, b)

	ta := a.GetNodeByName("leaf").Endpoint().Set(Props{"bandwidth": "10G"})
	tb := b.GetNodeByName("leaf").Endpoint().Set(Props{"bandwidth": "10G"})
	ab.Link(
		[]*Endpoint{ta},
		[]*Endpoint{tb},
	)

	buf, err := json.MarshalIndent(ab, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("44net.json", buf, 0644)

	VTag(ab)

	buf, err = json.MarshalIndent(ab, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("44net_vt.json", buf, 0644)

}
