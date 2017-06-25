/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^'''''''''''''
 * d2n - deter to nmir
 * ===================
 *   This program grabs the hardware toplogy from a deter database and makes an
 *   nmir file out of it.
 *
 *  Copyright The Deter Project 2017 - All Rights Reserved
 *  Apache 2.0 License
 *=================-----------------------------------------................```*/

package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rcgoodfellow/nmir"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
)

const theQ = `
select 
    node_id1, node_id2
from wires
left join nodes as n1 on n1.node_id = node_id1
left join nodes as n2 on n2.node_id = node_id2
left join interfaces as i1 
    on i1.node_id = node_id1 and
       i1.card   = card1 and
       i1.port   = port1
    left join interface_types as it1
        on it1.type = i1.interface_type
left join interfaces as i2 
    on i2.node_id = node_id2 and
       i2.card   = card2 and
       i2.port   = port2
    left join interface_types as it2
        on it2.type = i2.interface_type
`

type Wire struct {
	a, b string
}

func main() {

	log.SetFlags(0)

	f, err := os.Create("prof")
	if err != nil {
		log.Fatal(f)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	wires := collectWires()
	net := buildNet(wires)

	buf, err := json.MarshalIndent(net, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("dnet.json", buf, 0644)

	nmir.PNetSvg("qdnet", net)

	/*
		err = nmir.NetSvg("dnet", net)
		if err != nil {
			log.Fatal(err)
		}
	*/

	/*
		nmir.VTag(net)

		buf, err = json.MarshalIndent(net, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		ioutil.WriteFile("dnet_vt.json", buf, 0644)
	*/

}

func buildNet(wires []Wire) *nmir.Net {

	net := nmir.NewNet()

	m := make(map[string]*nmir.Node)
	for _, w := range wires {
		_, ok := m[w.a]
		if !ok {
			m[w.a] = net.Node().Set(nmir.Props{"name": w.a})
		}
		_, ok = m[w.b]
		if !ok {
			m[w.b] = net.Node().Set(nmir.Props{"name": w.b})
		}
	}

	for _, w := range wires {
		//a := net.Node().Set(nmir.Props{"name": w.a}).Endpoint()
		//b := net.Node().Set(nmir.Props{"name": w.b}).Endpoint()
		net.Link(
			[]*nmir.Endpoint{m[w.a].Endpoint()},
			[]*nmir.Endpoint{m[w.b].Endpoint()},
		)
	}

	for _, n := range net.Nodes {

		log.Printf("%s %d", n.Props["name"], n.Valence())

	}

	return net

}

func collectWires() []Wire {

	db, err := sql.Open("mysql", "unix(/tmp/mysql.sock)/tbdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(theQ)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var wires []Wire

	for rows.Next() {
		var a, b string
		err := rows.Scan(&a, &b)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s %s", a, b)
		wires = append(wires, Wire{a, b})
	}

	return wires

}
