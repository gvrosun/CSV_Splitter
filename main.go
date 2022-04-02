package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/go-gota/gota/dataframe"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func readCSV(fileName string, ch chan string) (int, dataframe.DataFrame) {
	csvFile, err := os.Open(fileName)
	if err != nil {
		ch <- "Unable to read file"
	}
	df := dataframe.ReadCSV(csvFile)
	return df.Nrow(), df
}

func processSplitting(fileName string, option string, numOf int, folder bool, progress *widget.ProgressBar, ch chan string) {
	fmt.Println(fileName, option, numOf, folder)
	if option == "" {
		ch <- "Please select the chunk type"
		return
	}
	if numOf == 0 {
		ch <- "Please enter " + option + " in Enter Value field"
		return
	}
	numOfRow, df := readCSV(fileName, ch)
	fmt.Println(numOfRow, df)
	for i := 0.0; i <= 1.0; i += 0.1 {
		time.Sleep(time.Millisecond * 250)
		progress.SetValue(i)
	}
	ch <- fmt.Sprintf("Success!\nRows: %d\nTotal chunks: %d", numOfRow, 0)
}

func setFileName(fileWidget *widget.Label, ch chan string) {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		ch <- "Permission Error"
	}

	for _, file := range files {
		if !file.IsDir() && strings.Contains(file.Name(), ".csv") {
			ch <- file.Name()
			return
		}
	}
	fileWidget.Text = "Not Found"
	ch <- "Error: CSV file not found"
	//*fileName = "Not Found"
}

func main() {
	a := app.New()
	w := a.NewWindow("Split CSV")
	w1 := a.NewWindow("Result")
	w.CenterOnScreen()
	w.Resize(fyne.Size{Width: 400, Height: 200})
	entryOption := ""
	isFolder := false
	ch := make(chan string)
	CSVFileName := ""

	fileName := widget.NewLabel("Detecting...")
	go setFileName(fileName, ch)
	go func() {
		result := <-ch
		if strings.Contains(result, "Error") {
			time.Sleep(time.Millisecond * 1000)
			w.Close()
			w1.SetContent(widget.NewLabel(result))
			w1.CenterOnScreen()
			w1.Show()
		}
		CSVFileName = result
		fileName.Text = result
	}()
	optionValue := widget.NewEntry()
	check := widget.NewCheck("Create folder for each chunk", func(value bool) {
		isFolder = value
	})
	combo := widget.NewSelect([]string{"Number of Rows", "Number of Chunks"}, func(value string) {
		entryOption = value
	})
	progress := widget.NewProgressBar()

	// Hide init widgets
	progress.Hidden = true

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "CSV FileName: ", Widget: fileName},
			{Text: "Chunk by", Widget: combo},
			{Text: "Enter Value", Widget: optionValue},
			{Text: "", Widget: check},
		},
		OnSubmit: func() { // optional, handle form submission
			progress.Hidden = false
			val, _ := strconv.Atoi(optionValue.Text)
			go processSplitting(CSVFileName, entryOption, val, isFolder, progress, ch)
			result := <-ch
			w.Close()
			w1.SetContent(widget.NewLabel(result))
			w1.CenterOnScreen()
			w1.Show()
		},
		OnCancel: func() {
			w.Close()
		},
		SubmitText: "Split Now",
		CancelText: "Close",
	}

	w.SetContent(container.NewVBox(form, progress))
	w.ShowAndRun()
}
