package main


import (
	//"github.com/bclicn/color"
	"log"

	"fmt"
	"time"
	"github.com/radovskyb/watcher"
	//"github.com/tideland/golib/scroller"
	"os"
	"bufio"
	"strings"
	"io"
	//"github.com/tideland/golib/errors"
	"errors"
	"regexp"
	"github.com/bclicn/color"
	"encoding/json"

	"github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
)

	//todo 3.1: Доделай все тудушки в этом севрисе
	//todo 4: Написать второй сервис - Обработчик этого спама событий

//todo: доделать список ошибок. Сделать отдельным файлом при финализированном оформлении проекта.
//// Error codes of the scroller package.
//const (
//	ErrNoReader = iota + 1
//	ErrNoWriter
//)
//
//var errorMessages = errors.Messages{
//	ErrNoReader:      "cannot start Analyser: no reader",
//	ErrNoWriter:      "cannot start Analyser: no writer",
//}
var NatsClient *nats.EncodedConn
var NatsHost = os.Getenv("NATS_HOST")
var NatsErr error


var (
	targetRowIdentifier = "/DynamicIVR : "
)


func init(){
	nc, err1 := nats.Connect(NatsHost)
	if err1 != nil {
		logrus.Printf("Nats error 1: %s", err1.Error())
	}
	NatsClient, NatsErr = nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if NatsErr != nil {
		logrus.Printf("Nats error 2: %s", NatsErr.Error())
	} else {
		logrus.Printf("Nats success connection: %s", NatsClient.Conn.ConnectedServerId())
	}
}

func main(){

	newFilesChan := make(chan string)

	//todo: move this in environment variable
	someDirectory := "logs"

	go initWatcher(someDirectory, newFilesChan)

	initScan(newFilesChan)

}

func initScan( channel chan string){
	for logFile := range channel{
		go scanFile(logFile)
	}
}

func scanFile(file string){

	reader, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	////defer reader.Close() //todo: узнать, зачем закрывать reader и когда..

	writer := new(NatsWriter)
	options := []string{"", "",} //todo: либо использовать, либо убрать..


	analyser, err := NewRowAnalyser(reader, writer, options...)
	if err != nil {}

	analyser.Analyse()

}

func NewRowAnalyser(reader io.Reader, writer io.Writer, options ...string) ( *RowAnalyser, error){
	if reader == nil {
		return nil, errors.New("no reader")
	}
	if writer == nil {
		return nil, errors.New("no writer")
	}
	a := &RowAnalyser{
		reader:     reader,
		writer:     writer,
		bufferSize: 4068,
	}

	//todo: тут могут добавляться параметры типа фильтров, сценариев, и тд..
	//for _, option := range options {
	//	//if err := option(s); err != nil {
	//	//	return nil, err
	//	//}
	//}

	a.scenarios = Scenarios

	return a, nil
}

type RowAnalyser struct{
	writer io.Writer
	reader io.Reader
	bufferSize int
	row string
	scenarios []Scenario
}

type rowObject struct{
	Date string
	Time string
	Host string
	Status string
	Hash string
	commandRow string
	Scenario string
	Attr string
}

func (ra *RowAnalyser) basicCheck() bool {
	switch {
	case len(ra.row) < 1 : return false
	case strings.Contains(ra.row, targetRowIdentifier) : return true

	default:
		return false
	}
}

func (ra *RowAnalyser) Analyse (){
	scanner := bufio.NewScanner(ra.reader)
	for scanner.Scan() {
		lineStr := scanner.Text()

		ra.row = lineStr
		if ra.basicCheck(){
			ra.CommitByScenarios()
		}
	}
}


func (ra *RowAnalyser) ParseRow () (row rowObject, err error) {

	s := ra.row

	var validDate = regexp.MustCompile(`[\d]{2}/[\d]{2}/[\d]{2}`)
	if ok := validDate.MatchString(s); !ok {
		return rowObject{}, errors.New("string doesn't contain Date")
	}

	var validTime = regexp.MustCompile(`[\d]{2}:[\d]{2}:[\d]{2}:[\d]{3}`)
	if ok := validTime.MatchString(s); !ok {
		return rowObject{}, errors.New("string doesn't contain Time")
	}

	str := strings.Split(s, targetRowIdentifier)
	if len(str) < 2 {
		return rowObject{}, errors.New("string doesn't contain DynamicIVR substr")
	}

	sysInfo := strings.Split(str[0], " - ")
	row.commandRow = str[1]

	if len(sysInfo) < 2 {
		return rowObject{}, errors.New("string doesn't have correct Data-Hash structure")
	}
	row.Hash = sysInfo[1]

	infoData := strings.Fields(sysInfo[0])
	if len(infoData) != 4{
		return rowObject{}, errors.New("string doesn't have correct Info structure")
	}

	row.Date = infoData[0]
	row.Time = infoData[1]
	row.Host = infoData[2]
	row.Status = infoData[3]

	return row, nil
}

func (ra *RowAnalyser) CommitByScenarios () {

	rowObject, err := ra.ParseRow()
	if err != nil {
		fmt.Printf(color.Red("row error :[%s] \n"), err.Error())
	}else{

		for scenario := range ra.scenarios{

			if ra.scenarios[scenario].verify(rowObject) {

				obj := ra.scenarios[scenario].finalise(rowObject)

				bytesRow, err := json.Marshal(obj)
				if err != nil {
					fmt.Printf(color.Red("error on marshalling :[%s] \n"), err.Error())
				}else{
					ra.writer.Write(bytesRow)
				}
			}
		}
	}
}



func initWatcher(directory string, channel chan string){
	w := watcher.New()

	w.FilterOps(watcher.Rename, watcher.Move, watcher.Create, watcher.Write)

	go func() {
		for {
			select {
			case event := <-w.Event:
				if !event.IsDir() && event.Op == watcher.Create{
					channel <- event.Path
				}

			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.Add(directory); err != nil {
		log.Fatalln(err)
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.

	//todo: надо додумать как не обрабатывать файлы повторно
	for path, f := range w.WatchedFiles() {
		if f.Name() != directory{
			//fmt.Printf("%s: %s\n", path, f.Name())
			channel <- path
		}
	}

	//fmt.Println()

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

