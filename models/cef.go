package models

import (
	"fmt"
	"github.com/rcgoodfellow/nmir"
	"github.com/satori/go.uuid"
)

func NewOS(name string, arch ...string) *nmir.Software {

	arch_ := make([]interface{}, len(arch))
	for i, x := range arch {
		arch_[i] = x
	}

	return &nmir.Software{
		Props: nmir.Props{
			"kind": "os",
			"name": name,
			"id":   uuid.NewV4().String(),
		},
		Requirements: nmir.Props{
			"arch?": arch_,
		},
	}

}

func CEF_Software1() []*nmir.Software {

	return []*nmir.Software{
		//*nix
		NewOS("debian-stable", "x86_64", "x86"),
		NewOS("debian-testing", "x86_64", "x86"),
		NewOS("centos-7", "x86_64", "x86"),
		NewOS("ubuntu-1604", "x86_64", "x86"),
		NewOS("freebsd-11", "x86_64", "x86"),

		//m$ft
		NewOS("windows-7", "x86_64", "x86"),
		NewOS("windows-10", "x86_64", "x86"),

		//network
		NewOS("cumulus", "x86_64", "x86", "arm7l"),
		NewOS("vyatta", "x86_64", "x86", "mips"),
		NewOS("palo-alto", "x86_64", "x86"),

		NewOS("android", "arm7l", "armhf"),
		NewOS("android-7", "arm7l", "armhf"),
		NewOS("android-6", "arm7l", "armhf"),
		NewOS("lte-linux", "x86_64", "armhf"),

		//IoT/Smart Home/Entertainment
		NewOS("plex", "arm7l", "armhf"),
	}

}

func CEF_3bed() *nmir.Net {

	tb := nmir.NewNet()

	node := func(format string, args ...interface{}) *nmir.Node {
		return tb.Node().Set(nmir.Props{"name": fmt.Sprintf(format, args...)})
	}

	// The Internet
	//internet := node("internet")

	// edge nodes
	egw := node("egw")
	mgw := node("mgw")
	hgw := node("hgw")

	// emulation-mobile site link
	tb.Link(
		egw.Endpoint().Set(nmir.Props{
			"bandwidth": nmir.Unit("=", 10, "gbps"),
		}),
		mgw.Endpoint().Set(nmir.Props{
			"bandwidth": nmir.Unit("=", 1, "gbps"),
		}),
	).Set(nmir.Props{
		"bandwidth": nmir.Unit("-", 3, "mbps"),
		"latency":   nmir.Unit("+", 11, "ms"),
		"loss":      nmir.Unit("-", 0.3, "%"),
	})

	// emulation-hpc site link
	tb.Link(
		egw.Endpoint().Set(nmir.Props{
			"bandwidth": nmir.Unit("=", 10, "gbps"),
		}),
		hgw.Endpoint().Set(nmir.Props{
			"bandwidth": nmir.Unit("=", 1, "gbps"),
		}),
	).Set(nmir.Props{
		"bandwidth": nmir.Unit("-", 22, "mbps"),
		"latency":   nmir.Unit("+", 7, "ms"),
		"loss":      nmir.Unit("-", 0.2, "%"),
	})

	// mobile-hpc site link
	tb.Link(
		mgw.Endpoint().Set(nmir.Props{
			"bandwidth": nmir.Unit("=", 1, "gbps"),
		}),
		hgw.Endpoint().Set(nmir.Props{
			"bandwidth": nmir.Unit("=", 1, "gbps"),
		}),
	).Set(nmir.Props{
		"bandwidth": nmir.Unit("-", 1.8, "mbps"),
		"latency":   nmir.Unit("+", 19, "ms"),
		"loss":      nmir.Unit("-", 0.5, "%"),
	})

	// Emulation Host

	hbe0 := node("hbe0")
	tb.Link(
		egw.Endpoint().Set(nmir.Props{
			"bandwidth": nmir.Unit("=", 40, "gbps"),
		}),
		hbe0.Endpoint().Set(nmir.Props{
			"bandwidth": nmir.Unit("=", 40, "gbps"),
		}),
	)

	var tor [3]*nmir.Node
	for i := 0; i < 3; i++ {
		tor[i] = node("tor%d", i)
		tb.Link(
			tor[i].Endpoint().Set(nmir.Props{
				"bandwidth": nmir.Unit("=", 10, "gbps"),
			}),
			hbe0.Endpoint().Set(nmir.Props{
				"bandwidth": nmir.Unit("=", 10, "gbps"),
			}),
		)
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			n := tb.Node().Set(nmir.Props{
				"name":  fmt.Sprintf("n%d", i*3+j),
				"vhost": "yes",
				"arch":  "x86_64",
			})
			tb.Link(
				tor[i].Endpoint().Set(nmir.Props{
					"bandwidth": nmir.Unit("=", 1, "gbps"),
				}),
				n.Endpoint().Set(nmir.Props{
					"bandwidth": nmir.Unit("=", 1, "gbps"),
				}),
			)
		}
	}

	// Mobile/IoT host

	hbe1 := node("hbe1")
	tb.Link(
		hbe1.Endpoint().Set(nmir.Props{
			"bandwidth": nmir.Unit("=", 10, "gbps"),
		}),
		mgw.Endpoint().Set(nmir.Props{
			"bandwidth": nmir.Unit("=", 10, "gbps"),
		}),
	)

	wap := node("wap")
	tb.Link(
		wap.Endpoint().Set(nmir.Props{
			"bandwidth": nmir.Unit("=", 10, "gbps"),
		}),
		hbe1.Endpoint().Set(nmir.Props{
			"bandwidth": nmir.Unit("=", 10, "gbps"),
		}),
	)

	for i := 0; i < 5; i++ {
		n := tb.Node().Set(nmir.Props{
			"name": fmt.Sprintf("droid%d", i),
			"arch": "armhf",
		})
		tb.Link(
			wap.Endpoint().Set(nmir.Props{
				"bandwidth": nmir.Unit("=", 100, "mbps"),
			}),
			n.Endpoint().Set(nmir.Props{
				"bandwidth": nmir.Unit("=", 100, "mbps"),
			}),
		)
	}
	for i := 0; i < 4; i++ {
		n := tb.Node().Set(nmir.Props{
			"name": fmt.Sprintf("rpi%d", i),
			"arch": "arm7l",
		})
		tb.Link(
			wap.Endpoint().Set(nmir.Props{
				"bandwidth": nmir.Unit("=", 100, "mbps"),
			}),
			n.Endpoint().Set(nmir.Props{
				"bandwidth": nmir.Unit("=", 100, "mbps"),
			}),
		)
	}

	// HPC

	hbe2 := node("hbe2")
	tb.Link(
		hbe2.Endpoint().Set(nmir.Props{
			"bandwidth": nmir.Unit("=", 10, "gbps"),
		}),
		hgw.Endpoint().Set(nmir.Props{
			"bandwidth": nmir.Unit("=", 10, "gbps"),
		}),
	)

	fab := node("fab")
	tb.Link(
		fab.Endpoint().Set(nmir.Props{
			"bandwidth": nmir.Unit("=", 10, "gbps"),
		}),
		hbe2.Endpoint().Set(nmir.Props{
			"bandwidth": nmir.Unit("=", 10, "gbps"),
		}),
	)

	for i := 0; i < 5; i++ {
		n := tb.Node().Set(nmir.Props{
			"name": fmt.Sprintf("hpc%d", i),
			"arch": "x86_64",
		})
		tb.Link(
			n.Endpoint().Set(nmir.Props{
				"bandwidth": nmir.Unit("=", 40, "gbps"),
			}),
			fab.Endpoint().Set(nmir.Props{
				"bandwidth": nmir.Unit("=", 40, "gbps"),
			}),
		)
	}

	return tb

}

func CEF_3bed_spine_leaf() *nmir.Net {

	tb := nmir.NewNet()

	node := func(format string, args ...interface{}) *nmir.Node {
		return tb.Node().Set(nmir.Props{"name": fmt.Sprintf(format, args...)})
	}

	// The Internet
	internet := node("internet")

	// Emulation Host
	egw := node("egw")
	tb.Link(egw.Endpoint(), internet.Endpoint())

	stem := node("stem")
	tb.Link(egw.Endpoint(), stem.Endpoint())

	var spine, leaf [3]*nmir.Node
	for i := 0; i < 3; i++ {
		spine[i] = node("spine%d", i)
		leaf[i] = node("leaf%d", i)
		tb.Link(spine[i].Endpoint(), stem.Endpoint())
		tb.Link(leaf[i].Endpoint(), stem.Endpoint())
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			tb.Link(leaf[i].Endpoint(), spine[j].Endpoint())
			n := tb.Node().Set(nmir.Props{
				"name":  fmt.Sprintf("n%d", i*3+j),
				"vhost": "yes",
				"arch":  "x86_64",
			})
			tb.Link(leaf[i].Endpoint(), n.Endpoint())
		}
	}

	// Mobile/IoT host
	mgw := node("mgw")
	tb.Link(mgw.Endpoint(), internet.Endpoint())

	wap := node("wap")
	tb.Link(wap.Endpoint(), mgw.Endpoint())

	for i := 0; i < 5; i++ {
		n := tb.Node().Set(nmir.Props{
			"name": fmt.Sprintf("droid%d", i),
			"arch": "armhf",
		})
		tb.Link(wap.Endpoint(), n.Endpoint())
	}
	for i := 0; i < 4; i++ {
		n := tb.Node().Set(nmir.Props{
			"name": fmt.Sprintf("rpi%d", i),
			"arch": "arm7l",
		})
		tb.Link(wap.Endpoint(), n.Endpoint())
	}

	// HPC
	hgw := node("hgw")
	tb.Link(hgw.Endpoint(), internet.Endpoint())

	fab := node("fab")
	tb.Link(fab.Endpoint(), hgw.Endpoint())

	for i := 0; i < 5; i++ {
		n := tb.Node().Set(nmir.Props{
			"name": fmt.Sprintf("hpc%d", i),
			"arch": "x86_64",
		})
		tb.Link(n.Endpoint(), fab.Endpoint())
	}

	return tb

}

func CEF_SmallWorld() *nmir.Net {

	world := nmir.NewNet()

	var routers [3]*nmir.Node

	//Internet Model
	for i := 0; i < 3; i++ {
		routers[i] = world.Node().
			Set(nmir.Props{
				"name": fmt.Sprintf("r%d", i),
			}).
			AddSoftware(nmir.Props{
				"name=": "cumulus",
			})

	}
	for i := 0; i < 3; i++ {
		world.Link(routers[i].Endpoint(), routers[(i+1)%3].Endpoint()).
			Set(nmir.Props{
				"bandwidth": nmir.Unit("-", 2.5, "mbps"),
				"latency":   nmir.Unit("+", 7, "ms"),
				"loss":      nmir.Unit("-", 3, "%"),
			})
	}

	//Cellular Network
	tower := world.Node().
		Set(nmir.Props{
			"name": "tower",
		}).
		AddSoftware(nmir.Props{
			"name=": "lte-linux",
		})

	for i := 0; i < 3; i++ {
		n := world.Node().
			Set(nmir.Props{
				"name": fmt.Sprintf("droid%d", i),
			}).
			AddSoftware(nmir.Props{
				"name=": "android-7",
			})
		world.Link(n.Endpoint(), tower.Endpoint()).
			Set(nmir.Props{
				"bandwidth": nmir.Unit("-", 1, "mbps"),
				"latency":   nmir.Unit("+", 27, "ms"),
				"loss":      nmir.Unit("-", 10, "%"),
			})
	}
	world.Link(tower.Endpoint(), routers[0].Endpoint()).
		Set(nmir.Props{
			"bandwidth": nmir.Unit("-", 100, "mbps"),
			"latency":   nmir.Unit("+", 3, "ms"),
			"loss":      nmir.Unit("-", 2, "%"),
		})

	//Home Network
	rtr := world.Node().
		Set(nmir.Props{
			"name": "rtr",
		}).
		AddSoftware(nmir.Props{
			"name=": "vyatta",
		})
	world.Link(rtr.Endpoint(), routers[1].Endpoint()).
		Set(nmir.Props{
			"bandwidth": nmir.Unit("-", 25, "mbps"),
			"latency":   nmir.Unit("+", 11, "ms"),
			"loss":      nmir.Unit("-", 4, "%"),
		})

	n := world.Node().
		Set(nmir.Props{
			"name": "bob-phone",
		}).
		AddSoftware(nmir.Props{
			"name=": "android",
		})
	world.Link(n.Endpoint(), rtr.Endpoint()).
		Set(nmir.Props{
			"bandwidth": nmir.Unit("-", 100, "mbps"),
			"latency":   nmir.Unit("+", 5, "ms"),
			"loss":      nmir.Unit("-", 2, "%"),
		})

	n = world.Node().
		Set(nmir.Props{
			"name": "alice-phone",
		}).
		AddSoftware(nmir.Props{
			"name=": "android-6",
		})
	world.Link(n.Endpoint(), rtr.Endpoint()).
		Set(nmir.Props{
			"bandwidth": nmir.Unit("-", 100, "mbps"),
			"latency":   nmir.Unit("+", 8, "ms"),
			"loss":      nmir.Unit("-", 3, "%"),
		})

	n = world.Node().
		Set(nmir.Props{
			"name": "tv",
		}).
		AddSoftware(nmir.Props{
			"name=": "plex",
		})
	world.Link(n.Endpoint(), rtr.Endpoint()).
		Set(nmir.Props{
			"bandwidth": nmir.Unit("-", 100, "mbps"),
			"latency":   nmir.Unit("+", 4, "ms"),
			"loss":      nmir.Unit("-", 2, "%"),
		})

	n = world.Node().
		Set(nmir.Props{
			"name": "homepc",
		}).
		AddSoftware(nmir.Props{
			"name=": "windows-10",
		})
	world.Link(n.Endpoint(), rtr.Endpoint()).
		Set(nmir.Props{
			"bandwidth": nmir.Unit("-", 1, "gbps"),
			"latency":   nmir.Unit("+", 0.1, "ms"),
			"loss":      nmir.Unit("-", 0.1, "%"),
		})

	//Service Provider
	gw := world.Node().
		Set(nmir.Props{
			"name": "edge",
			"phys": "true",
		}).
		AddSoftware(nmir.Props{
			"name=": "freebsd-11",
		})
	world.Link(gw.Endpoint(), routers[1].Endpoint()).
		Set(nmir.Props{
			"bandwidth": nmir.Unit("-", 100, "gbps"),
			"latency":   nmir.Unit("+", 0.1, "ms"),
			"loss":      nmir.Unit("-", 0.1, "%"),
		})

	for i := 0; i < 4; i++ {
		n = world.Node().
			Set(nmir.Props{
				"name": fmt.Sprintf("srv%d", i),
			}).
			AddSoftware(nmir.Props{
				"name=": "centos-7",
			})
		world.Link(gw.Endpoint(), n.Endpoint()).
			Set(nmir.Props{
				"bandwidth": nmir.Unit("-", 40, "gbps"),
				"latency":   nmir.Unit("+", 0.1, "ms"),
				"loss":      nmir.Unit("-", 0.1, "%"),
			})
	}

	//Enterprise
	fw := world.Node().
		Set(nmir.Props{
			"name": "firewall",
		}).
		AddSoftware(nmir.Props{
			"name=": "palo-alto",
		})
	world.Link(fw.Endpoint(), routers[2].Endpoint()).
		Set(nmir.Props{
			"bandwidth": nmir.Unit("-", 100, "mbps"),
			"latency":   nmir.Unit("+", 3, "ms"),
			"loss":      nmir.Unit("-", 0.8, "%"),
		})

	bk := world.Node().
		Set(nmir.Props{
			"name": "backbone",
		}).
		AddSoftware(nmir.Props{
			"name=": "cumulus",
		})
	world.Link(bk.Endpoint(), fw.Endpoint()).
		Set(nmir.Props{
			"bandwidth": nmir.Unit("-", 10, "gbps"),
			"latency":   nmir.Unit("+", 3, "ms"),
			"loss":      nmir.Unit("-", 0.8, "%"),
		})

	for i := 0; i < 2; i++ {
		wk := world.Node().Set(nmir.Props{
			"name": fmt.Sprintf("workgroup%d", i),
		}).
			AddSoftware(nmir.Props{
				"name=": "cumulus",
			})
		world.Link(wk.Endpoint(), bk.Endpoint()).
			Set(nmir.Props{
				"bandwidth": nmir.Unit("-", 10, "gbps"),
				"latency":   nmir.Unit("+", 2, "ms"),
				"loss":      nmir.Unit("-", 0.1, "%"),
			})

		for j := 0; j < 3; j++ {
			n := world.Node().Set(
				nmir.Props{
					"name": fmt.Sprintf("workstation%d", i*2+j),
					"phys": "true",
				}).
				AddSoftware(nmir.Props{
					"name=": "windows-7",
				})
			world.Link(wk.Endpoint(), n.Endpoint()).
				Set(nmir.Props{
					"bandwidth": nmir.Unit("-", 1, "gbps"),
					"latency":   nmir.Unit("+", 2, "ms"),
					"loss":      nmir.Unit("-", 0.2, "%"),
				})
		}
	}

	return world

}
