package nmir

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func net4() Net {

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

	host.Link().Endpoints = [2][]*Endpoint{
		[]*Endpoint{zwitch.Endpoints[0]},
		[]*Endpoint{a.Endpoints[0]},
	}

	host.Link().Endpoints = [2][]*Endpoint{
		[]*Endpoint{zwitch.Endpoints[1]},
		[]*Endpoint{b.Endpoints[0]},
	}

	host.Link().Endpoints = [2][]*Endpoint{
		[]*Endpoint{zwitch.Endpoints[2]},
		[]*Endpoint{c.Endpoints[0]},
	}

	host.Link().Endpoints = [2][]*Endpoint{
		[]*Endpoint{zwitch.Endpoints[3]},
		[]*Endpoint{d.Endpoints[0]},
	}

	return host

}

func TestBuilModels(t *testing.T) {

	host := net4()

	buf, _ := json.MarshalIndent(host, "", "  ")
	ioutil.WriteFile("4net.json", buf, 0644)

	host_vt := VTag(host)

	buf, _ = json.MarshalIndent(host_vt, "", "  ")
	ioutil.WriteFile("4net_vt.json", buf, 0644)

}
