package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/xuri/excelize/v2"
)

var (
	path          string
	file          string
	prefixOutput  = "ret_"
	listSheet     bool
	sheetList     = []string{}
	sheetString   string
	partnerString string
)

type File struct {
	*excelize.File
}

func main() {
	readArgument()
	f := readFile()
	f.setSheetList()

	fmt.Println("Selected List: ", sheetString)
	fmt.Printf("Selected Partner: %s\n\n", partnerString)
	fmt.Println("-> Processing ")
	var wg sync.WaitGroup
	for _, sheet := range sheetList {
		wg.Add(1)
		go func(wg *sync.WaitGroup, sheet string) {
			defer wg.Done()
			var calPartner int = 3
			var calMoney int = 34
			for {
				rows, _ := f.GetRows(sheet)
				if calMoney > len(rows) {
					break
				}
				var partnerValue string
				if tmp := strings.Split(rows[calPartner-1][0], " "); len(tmp) < 2 {
					return
				} else {
					partnerValue = tmp[1]
				}
				if rows[calMoney-1][5] == "0.00" || ((partnerString != "" && partnerString != "*") && !strings.Contains(partnerString, partnerValue)) {
					for i := calMoney - 33; i <= calMoney+6; i++ {
						f.RemoveRow(sheet, calMoney-33)
					}
					if partnerValue != "" && partnerValue != "__________________________" {
						fmt.Println(partnerValue)
					}
				} else {
					fmt.Println(partnerValue, "--written")
					calPartner += 40
					calMoney += 40
				}
			}

		}(&wg, sheet)
	}
	wg.Wait()

	output := fmt.Sprintf("%s/%s%s", path, prefixOutput, file)
	err := f.SaveAs(output)
	if err != nil {
		log.Fatal("Save file fail: ", err)
	}
	fmt.Println("-> Done.")
	fmt.Printf("-> output file: %s\n", output)
}

func readArgument() {
	flag.StringVar(&path, "path", "", "path of data")
	flag.StringVar(&file, "file", "", "file of data")
	flag.StringVar(&sheetString, "sheet", "*", "select sheet list, Example * or สงขลา,สุราษฎร์ธานี")
	flag.StringVar(&partnerString, "partner", "*", "select partner list, Example * or ดงเฮง,เกิดแก้ว")
	flag.BoolVar(&listSheet, "list-sheet", false, "list sheet")
	flag.Parse()

	if listSheet {
		f := readFile()
		defer f.Close()
		fmt.Println("Sheet list: ", f.GetSheetList())
		os.Exit(0)
	}
}

func readFile() File {
	input := fmt.Sprintf("%s/%s", path, file)
	fmt.Println("file input => ", input)
	f, err := excelize.OpenFile(input)
	if err != nil {
		log.Panic(err)
	}
	return File{
		File: f,
	}
}

func (f *File) setSheetList() {
	if sheetString == "" || sheetString == "*" {
		sheetList = f.GetSheetList()
	} else {
		sheetList = strings.Split(sheetString, ",")
		for _, sheet := range f.GetSheetList() {
			if !strings.Contains(sheetString, sheet) {
				f.DeleteSheet(sheet)
			}
		}
	}
}
