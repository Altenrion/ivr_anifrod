package main

import (
	"strings"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"os"
	"bufio"
	"github.com/bclicn/color"

	"log"
)

// Print rows benchmark for logs
func CheckLogFiles(){
	//logs := []string{
	//	"15/02/2018 17:00:00:022 http-7080-79 DEBUG - C7E0F9D8AFD0EA9A83B3CB023FAEF49C:/DynamicIVR : ## TRACE ## model.ivrnodes.Menu:     timeout = 0 ms",
	//	"15/03xdsf/2018 17:00:00:022 http-7080-79 DEBUG - C7E0F9D8AFD0EA9A83B3CB023FAEF49C:/DynamicIVR : ## DEBUG ## model.ivrnodes.Menu:     result: '1'",
	//	"15/03/2018 17:00:00:022 http-7080-79 DEBUG - C7E0F9D8AFD0EA9A83B3CB023FAEF49C:/DynamicIVR : ## DEBUG ## model.ivrnodes.Menu:   Evaluating '_card_number_with_pin_count_more_or_equal_1' from: if (_card_number_with_pin_count >= 1) {result = true;} else {result = false;}",
	//	"15/03/2018 17:00:00:022 http-7080-79 DEBUG - C7E0F9D8AFD0EA9A83B3CB023FAEF49C:/DynamicIVR : ## DEBUG ## model.ivrnodes.Menu:    loading into JS engine: variable: _card_number_with_pin_count = '1' [Long]",
	//	"15/03/2018 17:0022:00:022 http-7080-79 DEBUG - C7E0BUG:/DynamicIVR : ## model.ivrnodes.Menu:    loading into JS engine: vari",
	//	"15/03/2018 17:00:00:023 http-7080-79 DEBUG - C7E0F9D8AFD0EA9A83B3CB023FAEF49C:/DynamicIVR : ## DEBUG ## model.ivrnodes.Menu:   script result: 'true'",
	//	"15/03/2018 17:00:00:023 http-7080-79 DEBUG - C7E0F9D8AFD0EA9A83B3CB023FAEF49C:/DynamicIVR : ## DEBUG ## model.ivrnodes.Menu:   script result: 'true'",
	//	"sdfsdfsdfsdfsd.wav",
	//	"sdfsdfsdfsdfsd.wav",
	//	"15/03/2018 17:00:00:023 http-7080-79 DEBUG - C7E0F9D8AFD0EA9A83B3CB023FAEF49C:/DynamicIVR : ## DEBUG ## model.ivrnodes.Menu:   script result: 'true'",
	//	"sdfsdfsdfsdfsd.wav",
	//}
	//for i, s := range logs {
	//	CheckLogRow(i,s)
	//}
	files, err := ioutil.ReadDir("logs")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(color.Blue(" filename     |  total rows |  valid rows  |  failed rows |  failed (perc) |  Errors  |  MemoryLeak \n"))

	for _, file := range files {
		if strings.Contains(file.Name(), "trace.log"){
			path := filepath.Join("logs", file.Name())
			printStatus(path)
		}
	}
}

func printStatus(filePath string){

	statusBuffer := StatusBuffer{}
	readLine(filePath, &statusBuffer)

	fmt.Printf(color.Blue(" %s  |"), filePath)
	fmt.Printf(color.Blue("   %d  |"), statusBuffer.SuccessCounter+statusBuffer.FailCounter)
	fmt.Printf(color.Yellow("   %d  |"), statusBuffer.SuccessCounter)
	fmt.Printf(color.Red("   %d   |   %f   |"), statusBuffer.FailCounter, (float64(statusBuffer.FailCounter)/float64(statusBuffer.SuccessCounter+statusBuffer.FailCounter))*100)
	fmt.Printf(color.Red(" %d  |"), len(statusBuffer.ErrorsMap))
	fmt.Printf(color.Red(" %f mb \n"), float64(statusBuffer.MemoryLeak) / 1000000)

}

func readLine(path string, statusBuffer *StatusBuffer) {
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		//CheckLogRow(0,scanner.Text(), statusBuffer)
	}
}


type StatusBuffer struct{
	SuccessCounter int
	FailCounter int
	ErrorsMap []string
	MemoryLeak uintptr

}

//func CheckLogRow(i int, s string, statusBuffer *StatusBuffer){
//	_, err := parseToStruct(s)
//
//	if err != nil {
//		statusBuffer.FailCounter++
//		statusBuffer.ErrorsMap = append(statusBuffer.ErrorsMap, err.Error())
//		statusBuffer.MemoryLeak = statusBuffer.MemoryLeak + unsafe.Sizeof(s)
//	}else{
//		statusBuffer.SuccessCounter++
//		//fmt.Printf(color.Blue("Row %d valid :[%s, %s, %s] \n"),i, rowData.Time, rowData.Hash, rowData.CommandRow)
//	}
//}

//func parseToStruct(s string ) (rowData LogRowData, err error) {
//
//	var validDate = regexp.MustCompile(`[\d]{2}/[\d]{2}/[\d]{2}`)
//	if ok := validDate.MatchString(s); !ok {
//		return LogRowData{}, errors.New("string doesn't contain Date")
//	}
//
//
//	var validTime = regexp.MustCompile(`[\d]{2}:[\d]{2}:[\d]{2}:[\d]{3}`)
//	if ok := validTime.MatchString(s); !ok {
//		return LogRowData{}, errors.New("string doesn't contain Time")
//	}
//
//	str := strings.Split(s, ":/DynamicIVR : ")
//	if len(str) == 2{
//		sysInfo := strings.Split(str[0], " - ")
//		rowData.CommandRow = str[1]
//
//		if len(sysInfo) == 2 {
//
//			rowData.Hash = sysInfo[1]
//
//			infoData := strings.Fields(sysInfo[0])
//			if len(infoData) == 4{
//				rowData.Date = infoData[0]
//				rowData.Time = infoData[1]
//				rowData.Host = infoData[2]
//				rowData.Status = infoData[3]
//			}else {
//				fmt.Printf(color.Red("struct :[%+v] \n"), rowData)
//			}
//		}
//	}else{
//		return LogRowData{}, errors.New("string doesn't contain DynamicIVR substr")
//	}
//
//	return rowData, nil
//}
