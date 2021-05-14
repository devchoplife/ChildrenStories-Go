package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/alecthomas/template"
)

var useCLI = flag.Bool("useCLI", false, "Do you want to use the Command Line Interface of the game")

type storyOption struct {
	Text string
	Arc  string
}

type storyArc struct {
	Title   string
	Story   []string
	Options []storyOption
}

type storyHandler struct {
	storyArcs map[string]storyArc
	template  *template.Template
}

type storyCLI struct {
	storyArcs map[string]storyArc
	reader    bufio.Reader
}

func (sh storyHandler) serveHTTP(res http.ResponseWriter, req *http.Request) {
	var path string
	if req.URL.Path == "/" {
		path = "intro"
	} else {
		path = strings.TrimLeft(req.URL.Path, "/")
	}
	storyArc := sh.storyArcs[path]
	sh.template.Execute(res, storyArc)
}

func main() {
	flag.Parse()

	f, err := os.Open("stories.json")
	if err != nil {
		panic(err)
	}

	var storyArcs map[string]storyArc
	err = json.Unmarshal(buf.Bytes(), &storyArcs)
	if err != nil {
		panic(err)
	}

	if *useCLI {
		storyCLI{storyArcs, bufio.NewReader(os.Stdin)}.start()
	} else {
		t, err := template.ParseFiles("templates/main.html")
		if err != nil {
			panic(err)
		}

		fmt.Println("Your advenure awaits, go to your browser and visit localhost:8080 to begin")
		http.ListenAndServe(":8080", storyHandler{storyArcs, t})
	}
}
