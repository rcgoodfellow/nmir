package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"html/template"
	"io/ioutil"
	"log"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {

	model := c.Params.Route.Get("model")
	log.Printf("model=%s", model)

	title := model
	moreStyles := []string{
		"css/style.css",
	}
	moreScripts := []string{
		"js/three.js",
		"js/viz.js",
	}

	buf, err := ioutil.ReadFile(fmt.Sprintf("/tmp/nmir/%s_vt.json", model))
	if err != nil {
		log.Println(err)
	}
	model_js := template.JS("var topo = " + string(buf) + ";")

	return c.Render(title, moreStyles, moreScripts, model_js)
}
