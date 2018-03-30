package main

import "strings"


// Scenarios codes.
const (
	UserEnteredIvr = iota
	UserEnteredButton
	UserDisconected
	UserNotActive
)

var Scenarios = []Scenario{
	UserEnteredIvr:    UserEnteredIvrScenario{}.init(),
	UserEnteredButton: UserEnteredButtonScenario{}.init(),
	UserDisconected:   UserDisconnectedScenario{}.init(),
	UserNotActive:     UserNotActiveScenario{}.init(),
}

type Scenario interface {
	verify(rowObject) bool
	init()Scenario
	finalise(rowObject) rowObject
}

type ScenarioVerifier struct{
	Feature string
}
func (sv *ScenarioVerifier) verify(LogRowObject rowObject) bool {
	return strings.Contains(LogRowObject.commandRow, sv.Feature)
}
func (sv *ScenarioVerifier) setFeature(feature string){
	sv.Feature = feature
}

type UserEnteredIvrScenario struct {
	ScenarioVerifier
}
func (sc UserEnteredIvrScenario ) init() Scenario{
	scenario :=  &UserEnteredIvrScenario{}
	scenario.setFeature("flow.GetIvrTree: DNIS: ")
	return scenario
}
func (sc UserEnteredIvrScenario ) finalise(row rowObject)  rowObject{

	str := strings.Split(row.commandRow, sc.Feature)
	if len(str) == 2{
		row.Scenario = "UserEnteredIvrScenario"
		row.Attr = str[1]
	}

	return row
}

type UserEnteredButtonScenario struct {
	ScenarioVerifier
}
func (sc UserEnteredButtonScenario ) init() Scenario{
	scenario :=  &UserEnteredButtonScenario{}
	scenario.setFeature("# Button   : ")
	return scenario
}
func (sc UserEnteredButtonScenario ) finalise(row rowObject)  rowObject{

	str := strings.Split(row.commandRow, sc.Feature)
	if len(str) == 2{
		row.Scenario = "UserEnteredButtonScenario"
		row.Attr = str[1]
	}

	return row
}

type UserDisconnectedScenario struct {
	ScenarioVerifier
}
func (sc UserDisconnectedScenario) init() Scenario{
	scenario :=  &UserDisconnectedScenario{}
	scenario.setFeature("actions.DisconnectNodeAction: Disconnect")
	return scenario
}
func (sc UserDisconnectedScenario) finalise(row rowObject)  rowObject{

	row.Scenario = "UserDisconnectedScenario"
	return row
}

type UserNotActiveScenario struct {
	ScenarioVerifier
}
func (sc UserNotActiveScenario ) init() Scenario{
	scenario :=  &UserNotActiveScenario{}
	scenario.setFeature("actions.MenuNodeAction: # NoInput")
	return scenario
}
func (sc UserNotActiveScenario ) finalise(row rowObject)  rowObject{

	row.Scenario = "UserNotActiveScenario"
	return row
}

