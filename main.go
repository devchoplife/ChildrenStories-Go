package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

func (scli storyCLI) getStoryOption(options []storyOption) storyOption {
	numChoices := len(options)

	choice, err := scli.reader.ReadString('\n')
	if err != nil {
		fmt.Println("Sorry, your choice could not be read, please try again............")
		return scli.getStoryOption(options)
	}

	choiceNr, err := strconv.Atoi(strings.TrimRight(choice, "\n"))
	if err != nil {
		fmt.Println("Sorry, your choice has tp be a number, please try again..............")
		return scli.getStoryOption(options)
	}
	if choiceNr <= 0 || choiceNr > numChoices {
		fmt.Printf("Sorry your choice ha to be a number between %d and %d. Please try again..............", 1, numChoices)
		return scli.getStoryOption(options)
	}
	return options[choiceNr-1]
}

func (scli storyCLI) presentStoryArc(storyArcName string) {
	storyArc := scli.storyArcs[storyArcName]

	fmt.Printf("\n--- %s ---\n\n", storyArc.Title)

	for _, p := range storyArc.Story {
		fmt.Println(p)
	}

	if len(storyArc.Options) == 0 {
		fmt.Println("\n\n...............Your adventure has ended...............\n\n")
		os.Exit(0)
	}

	fmt.Println("\nWhat will you do next?")
	for i, opt := range storyArc.Options {
		fmt.Printf("%d: %s", i+1, opt.Text)
	}
	fmt.Println("")

	so := scli.getStoryOption(storyArc.Options)
	scli.presentStoryArc(so.Arc)
}

func (scli storyCLI) start() {
	scli.presentStoryArc("intro")
}

func main() {
	flag.Parse()

	f, err := os.Open("stories.json")
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
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
