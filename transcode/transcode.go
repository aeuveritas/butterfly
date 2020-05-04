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

func isExisted(fp string, isMandatory bool) (ret bool) {
	ret = true
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		if isMandatory {
			createDirectory(fp)
		} else {
			ret = false
		}
	}
	return ret
}

func createDirectory(path string) {
	err := os.Mkdir(path, 0755)
	if err != nil {
		panic(err)
	}
	fmt.Println("create output directory: ", path)
}

func getAbsPath(fp string) string {
	absPath := fp
	if !filepath.IsAbs(fp) {
		_absPath, err := filepath.Abs(fp)
		if err != nil {
			panic("path is invalid: " + fp)
		}
		absPath = _absPath
	}

	return absPath
}

func getNewFilename(absPath string, title string, isVideo bool) string {
	currentTime := time.Now()
	timeString := currentTime.Format(time.RFC3339)
	var format string
	if isVideo {
		format = ".mp4"
	} else {
		format = ".mp3"
	}
	filename := title + "_" + timeString + format
	outputFile := filepath.Join(absPath, filename)

	return outputFile
}

func getOutputFileName(outputPath string, title string, isVideo bool) string {
	absPath := getAbsPath(outputPath)
	isExisted(absPath, true)

	outputFile := getNewFilename(absPath, title, isVideo)
	for {
		if isExisted(outputFile, false) {
			fmt.Println("already existed: ", outputFile)
			outputFile = getNewFilename(absPath, title, isVideo)
		} else {
			break
		}
	}

	return outputFile
}

// Run transcode media stream to file
func Run(inputURL string, outputFile string, durationString string, isVideo bool, isDebug bool) {
	// Create new instance of transcoder
	trans := new(transcoder.Transcoder)

	// Initialize transcoder passing the input url and output file path
	err := trans.Initialize(inputURL, outputFile)
	if err != nil {
		panic(err)
	}

	trans.MediaFile().SetDuration(durationString)
	trans.MediaFile().SetOutputPath(outputFile)
	if isVideo {
		trans.MediaFile().SetVideoCodec("libx264")
	} else {
		trans.MediaFile().SetSkipVideo(true)
	}

	// Start transcoder process without checking progress
	fmt.Println("new file: ", outputFile)
	done := trans.Run(false)

	// This channel is used to wait for the process to end
	err = <-done
	if err != nil {
		panic(err)
	}
}

// Preset struct for preset
type Preset struct {
	Title      string `json:"title"`
	InputURL   string `json:"inputURL"`
	OutputPath string `json:"outputPath"`
	Duration   int    `json:"duration"`
	Token      string `json:"token,omitempty"`
	Video      bool   `json:"video"`
}

func getDurationString(duration int, isDebug bool) string {
	var weight int = 1
	if !isDebug {
		weight = 60
	}
	return strconv.Itoa(duration * weight)
}

// ParsePreset parse preset
func ParsePreset(preset string, isDebug bool) (string, string, string, string, bool) {
	absPresetPath := getAbsPath(preset)
	if !isExisted(absPresetPath, false) {
		panic("wrong preset file: " + absPresetPath)
	}

	file, _ := ioutil.ReadFile(absPresetPath)
	data := Preset{}
	_ = json.Unmarshal([]byte(file), &data)

	outputFile := getOutputFileName(data.OutputPath, data.Title, data.Video)

	return data.InputURL, outputFile, getDurationString(data.Duration, isDebug), data.Token, data.Video
}

// ParseParameter parse parameter
func ParseParameter(input string, output string, duration int, title string, isVideo bool) (string, string, string) {
	inputURL := input
	outputFile := getOutputFileName(output, title, isVideo)

	return inputURL, outputFile, getDurationString(duration, false)
}
