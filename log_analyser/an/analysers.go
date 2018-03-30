package main

import (
	"fmt"
	"github.com/bclicn/color"
	"encoding/json"
)

// Analysers codes.
const (
	DefaultAnalyser       = iota + 1
	TimePerCallAnalyser
	CallFrequencyAnalyser
	UserActionsAnalyser
)

var AnalysersFuncs = [] func(data interface{}, resultChan chan Result, deeper *bool){
	DefaultAnalyser: func(data interface{}, resultChan chan Result, deeper *bool) {

		resultChan <- Result{"00000", "Default analyser done."}
	},
	TimePerCallAnalyser: func(data interface{}, resultChan chan Result, deeper *bool) {
		//todo: analyse session time from first action in session till last action

		*deeper = false
		resultChan <- Result{"00001", "TimePerCallAnalyser done his job"}
	},
	CallFrequencyAnalyser: func(data interface{}, resultChan chan Result, deeper *bool) {
		//todo: analyse amount of calls from one number & time between this calls

		resultChan <- Result{"00002", "CallFrequencyAnalyser done his job"}
	},
	UserActionsAnalyser: func(data interface{}, resultChan chan Result, deeper *bool) {
		//todo: analyse sessions from one number if they have same combinations of actions (by actions hash)

		resultChan <- Result{"00003", "UserActionsAnalyser done his job"}
	},
}

type Result struct {
	Identifier string
	Status     string
}

type Analyser interface {
	Analyse(interface{}, chan Result)
	GetAnalyticFunc() func(interface{}, chan Result, *bool)
}

type Composite interface {
	getChildren() []Analyser
	addChild(analyser Analyser)
}

type CompositeAnalyser interface {
	Analyser
	Composite
}

/////////////////////////////////////////
type CallsAnalyser struct {
	Id           int                   `json:"id"`
	Parent       int                   `json:"parent"`
	Children     map[int]CallsAnalyser `json:"children"`
	analyticFunc func(data interface{}, resultChan chan Result, deeper *bool)
}

func (n CallsAnalyser) getChildren() map[int]CallsAnalyser {
	return n.Children
}
func (n *CallsAnalyser) addChild(analyser CallsAnalyser) {

	if len(n.Children) == 0 {
		n.Children = make(map[int]CallsAnalyser)
	}

	n.Children[len(n.Children)] = analyser
}

func (n CallsAnalyser) GetAnalyticFunc() func(data interface{}, resultChan chan Result, deeper *bool) {
	return n.analyticFunc
}

func (n CallsAnalyser) Analyse(data interface{}, resultChan chan Result) {

	goDeeperFlag := true

	n.GetAnalyticFunc()(data, resultChan, &goDeeperFlag)

	if goDeeperFlag {
		subAnalysers := n.getChildren()
		for _, analyser := range subAnalysers {
			analyser.Analyse(data, resultChan)
		}
	}
}

func createTree(nodes *map[int]CallsAnalyser, parent map[int]CallsAnalyser) map[int]CallsAnalyser {
	tree := make(map[int]CallsAnalyser)

	nodesLoc := *nodes

	for _, node := range parent {
		if _, ok := nodesLoc[node.Id]; ok {
			children := createTree(nodes, nodesLoc[node.Id].Children)
			for _, child := range children {
				node.addChild(child)
			}
		}
		tree[len(tree)] = node
	}

	return tree
}

func load(jsonConfig string) CallsAnalyser {

	new := make(map[int]CallsAnalyser)
	root := make(map[int]CallsAnalyser)
	nodes := make(map[int]CallsAnalyser)

	analysersFuncs := AnalysersFuncs

	var analyserObjs []CallsAnalyser

	err := json.Unmarshal([]byte(jsonConfig), &analyserObjs)
	if err != nil {
		fmt.Printf("Error: %s \n", err.Error())
	}

	for k, analyserObj := range analyserObjs {
		analyserObj.analyticFunc = analysersFuncs[analyserObj.Id]
		analyserObj.Children = make(map[int]CallsAnalyser)

		nodes[k] = analyserObj
	}

	for _, node := range nodes {
		nNode := new[node.Parent]
		nNode.addChild(node)
		new[node.Parent] = nNode
	}

	root[0] = nodes[0]

	treeData := createTree(&new, root)
	tree := treeData[0]

	return tree
}

type CompositeNode struct {
	Id       int             `json:"id"`
	Children []CompositeNode `json:"children"`
	Parent   int             `json:"parent"`
}

func main() {

	data := "string long enought"
	resultChan := make(chan Result)
	go func() {
		for {
			status := <-resultChan
			fmt.Printf(color.Blue("[%+v] \n"), status)
		}
	}()

	jsonConfig := `[
					{ "id": 1, "parent": 0 },
					{ "id": 2, "parent": 1 },
					{ "id": 3, "parent": 2 },
					{ "id": 4, "parent": 2 }]`

	analyser := load(jsonConfig)
	analyser.Analyse(data, resultChan)
}
