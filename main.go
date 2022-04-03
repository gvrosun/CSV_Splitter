package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/go-gota/gota/dataframe"
	"io/ioutil"
	"math"
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

func generateRange(min int, max int) []int {
	var seq []int
	for i := min; i < max; i++ {
		seq = append(seq, i)
	}
	return seq
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
	totalChunks := 1.0
	chunkRow := 0

	if option == "Number of Rows" {
		if numOfRow < numOf {
			ch <- "Given Number of Rows exceeds csv rows"
			return
		}
		remainder := numOfRow % numOf
		chunkRow = numOf
		totalChunks = math.Floor(float64(numOfRow / numOf))
		if remainder != 0 {
			totalChunks += 1
		}
	}
	if option == "Number of Chunks" {
		remainder := numOfRow % numOf
		chunkRow = numOfRow / numOf
		totalRows := math.Floor(float64(chunkRow))
		if remainder != 0 {
			chunkRow += 1
		}
		if totalRows == 0 {
			ch <- "Too large number of chunks"
			return
		}
		totalChunks = float64(numOf)
	}

	fmt.Println(totalChunks, chunkRow)
	start := 0
	end := chunkRow
	for i := 1.0; i <= totalChunks; i += 1.0 {
		newDF := df.Subset(generateRange(start, end))
		fileChunkName := ""
		if folder {
			folderChunkName := fmt.Sprintf("Chunk%d", int(i))
			err := os.Mkdir(".\\"+folderChunkName, 0755)
			if err != nil {
				ch <- "Error creating folder"
			}
			fileChunkName = ".\\" + folderChunkName + "\\" + strings.ReplaceAll(fileName, ".csv", "") + "_" + folderChunkName + ".csv"
		} else {
			fileChunkName = strings.ReplaceAll(fileName, ".csv", "") + "_" + fmt.Sprintf("Chunk%d.csv", int(i))
		}
		f, err := os.Create(fileChunkName)
		if err != nil {
			ch <- "Error creating chunks"
			return
		}

		err = newDF.WriteCSV(f)
		if err != nil {
			ch <- "Error creating chunks"
			return
		}
		progress.SetValue(i / totalChunks)
		start = end
		if (end + chunkRow) < numOfRow {
			end += chunkRow
		} else {
			end = numOfRow
		}
	}
	ch <- fmt.Sprintf("Success!\nRows: %d\nTotal chunks: %d", numOfRow, int(totalChunks))
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
