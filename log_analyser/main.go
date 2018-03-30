package main

import (
	"crypto/sha256"
	"fmt"
	"encoding/hex"
	"crypto/md5"
	"github.com/bclicn/color"
)

var sessionPool SessionsPool
var nTree Tree

func init() {
	//todo: тут подключаю Nats

	sessionPool = make(map[string]Session)

	nTree.nodes = make(map[string]NumberNode)
	nTree.ch = make(chan Status)

}

func main() {

	//todo: подписка на канал сканера
	//todo: луп по каналу.
	//todo: распихивание приходящих структур по полочкам
	//todo: реагирование по ходу дела.

	go func() {
		for {
			status := <-nTree.ch
			fmt.Printf(color.Blue("Statuses :[%+v] \n"), status)
		}
	}()

	testActions := Actions // test data

	for _, action := range testActions {
		if sessionPool.sessionExist(action.Hash) {
			sessionPool.setAction(action)
		} else {
			sessionPool.setSession(action)
		}
	}

	fmt.Printf("Tree :[%+v] \n \n", nTree)

	for number, treeBranch := range nTree.nodes {
		for _, session := range treeBranch.Sessions {
			fmt.Printf(color.BBlue("number[%s] : hash[%s] \n"), number, session.ActionsHash)
		}
	}
	fmt.Printf("Sessions :[%+v] \n \n", nTree.nodes["123456789"].Sessions[0].ActionsHash)

}

type Session struct {
	ActionsHash string
	Actions     []Action
	Number      string
}

func (s *Session) addAction(action Action) {
	s.Actions = append(s.Actions, action)
	s.ActionsHash = makeHash(s.Actions)
	if action.Scenario == "Number" {
		number := action.Attr

		numNode := NumberNode{}

		//если номер существует, дополняем узел ссылкой на сессию
		if nTree.numberExist(number) {
			numNode = nTree.getNumberNode(number)
		}
		numNode.setSession(s)
		nTree.setNumberNode(number, numNode)
	}
}

type tmpAttr struct {
	Scenario string
	Attr     string
}

func makeHash(actions []Action) string {

	var tmpSlice []tmpAttr
	for _, action := range actions {
		tmpSlice = append(tmpSlice, tmpAttr{Scenario: action.Scenario, Attr: action.Attr})
	}

	h := sha256.New()
	s := fmt.Sprintf("%v", tmpSlice)
	sum := h.Sum([]byte(s))

	hasher := md5.New()
	hasher.Write(sum)

	return hex.EncodeToString(hasher.Sum(nil))
}

type Action struct {
	Date     string
	Time     string
	Host     string
	Status   string
	Hash     string
	Scenario string
	Attr     string
}

type SessionsPool map[string]Session

func (sp SessionsPool) sessionExist(hash string) bool {
	if _, ok := sp[hash]; ok {
		return true
	}
	return false
}

//func (sp SessionsPool ) getSession (hash string)	(*Session, error) {}
func (sp *SessionsPool) setSession(action Action) {
	sP := *sp

	session := Session{}
	session.addAction(action)

	sP[action.Hash] = session

	*sp = sP
}
func (sp *SessionsPool) setAction(action Action) {
	sP := *sp

	session := sP[action.Hash]
	session.addAction(action)
	sP[action.Hash] = session

	*sp = sP
}

type NumberNode struct {
	Sessions []*Session
	Status   string
}

func (n *NumberNode) setSession(session *Session) {
	n.Sessions = append(n.Sessions, session)
}

type Result struct {
	Identifier string
	Status     string
}

type Tree struct {
	//ch chan Result
	nodes          map[string]NumberNode
	nodesAnalysers Analysers
}

func (t Tree) numberExist(number string) bool {
	if _, ok := t.nodes[number]; ok {
		return true
	}
	return false
}

func (t Tree) setNumberNode(number string, node NumberNode) {
	t.nodesAnalysers.Analyse(node)
	t.nodes[number] = node
}

func (t Tree) getNumberNode(number string) NumberNode {
	return t.nodes[number]
}
