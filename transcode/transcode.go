package transcode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/xfrr/goffmpeg/transcoder"
)

func isExisted(fp string) bool {
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		return false
	}
	return true
}

func isValidPath(fp string) (ret bool) {
	if !filepath.IsAbs(fp) {
		panic("not absolute path: " + fp)
	}

	if !isExisted(fp) {
		panic("path is not existed: " + fp)
	}

	return true
}

func getOutputFileName(outputPath string, title string) string {
	isValidPath(outputPath)

	currentTime := time.Now()
	outputFile := filepath.Join(outputPath, title+"_"+currentTime.Format("2020-01-01")+".flv")
	cnt := 0
	for {
		if isExisted(outputFile) {
			fmt.Println("already existed: ", outputFile)
			cnt++
			outputFile = filepath.Join(outputPath, title+"_"+currentTime.Format("2020-01-01")+"_"+strconv.Itoa(cnt)+".flv")
		} else {
			fmt.Println("new file: ", outputFile)
			break
		}
	}

	return outputFile
}

func transcode(inputURL string, outputFile string, durationString string) {
	// Create new instance of transcoder
	trans := new(transcoder.Transcoder)

	// Initialize transcoder passing the input file path and output file path
	err := trans.Initialize(inputURL, outputFile)
	if err != nil {
		panic(err)
	}

	trans.MediaFile().SetDuration(durationString)
	trans.MediaFile().SetOutputPath(outputFile)

	// Start transcoder process without checking progress
	done := trans.Run(false)

	// This channel is used to wait for the process to end
	err = <-done
	if err != nil {
		panic(err)
	}
}

// Run run butteryfly with parameters
func Run(input string, output string, duration int, title string) (outputFile string) {
	inputURL := input
	outputFile = getOutputFileName(output, title)
	//durationString := strconv.Itoa(duration * 60)
	durationString := strconv.Itoa(duration)

	transcode(inputURL, outputFile, durationString)

	return
}

// Preset struct for preset
type Preset struct {
	Title      string `json:"title"`
	InputURL   string `json:"inputURL"`
	OutputPath string `json:"outputPath"`
	Duration   string `json:"duration"`
	Token      string `json:"token,omitempty"`
}

func parsePreset(preset string) (string, string, string, string) {
	isValidPath(preset)

	file, _ := ioutil.ReadFile(preset)
	data := Preset{}

	_ = json.Unmarshal([]byte(file), &data)

	if !isExisted(data.OutputPath) {
		err := os.Mkdir(data.OutputPath, 0755)
		fmt.Println("create output directory:", data.OutputPath)
		if err != nil {
			panic(err)
		}
	}

	outputFile := getOutputFileName(data.OutputPath, data.Title)
	Token := data.Token
	fmt.Println(Token)
	return data.InputURL, outputFile, data.Duration, data.Token
}

// RunPreset run butterfly with preset
func RunPreset(preset string) (string, string) {
	inputURL, outputFile, durationString, token := parsePreset(preset)

	transcode(inputURL, outputFile, durationString)
	return outputFile, token
}
