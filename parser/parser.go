package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// TelegramParser parser for telegram
type TelegramParser struct {
	FileHelper

	InfoFile string
	InfoData TelegramInfo
}

// TelegramInfo telegram info
type TelegramInfo struct {
	Token string     `json:"token"`
	Chats []ChatInfo `json:"chat"`
}

// ChatInfo chat info
type ChatInfo struct {
	Title string `json:"title"`
	ID    string `json:"id"`
}

// PrintData print data
func (p *TelegramParser) PrintData() {
	jsonData, _ := json.MarshalIndent(p.InfoData, "", "\t")
	fmt.Println("TelegramInfo")
	fmt.Println(string(jsonData))
}

// Parse parse telegram json
func (p *TelegramParser) Parse() {
	if absPresetPath, valid := p.validate(p.InfoFile); valid {
		file, _ := ioutil.ReadFile(absPresetPath)
		_ = json.Unmarshal([]byte(file), &p.InfoData)
	}

	return
}

// AddItem add telegram info
func (d *TelegramInfo) AddItem(title string, ID string) {
	newChat := ChatInfo{
		Title: title,
		ID:    ID,
	}
	d.Chats = append(d.Chats, newChat)
}

// SaveTelegramInfo save telegram data to json
func (p *TelegramParser) SaveTelegramInfo() {
	absPresetPath := p.getAbsPath(p.InfoFile)
	if !p.isExisted(absPresetPath, false) {
		fmt.Println("telegram info is not submitted")
		return
	}

	data, _ := json.MarshalIndent(p.InfoData, "", "\t")
	_ = ioutil.WriteFile(p.InfoFile, data, 0644)
}

// PresetParser parser for preset
type PresetParser struct {
	FileHelper

	InfoFile string
	InfoData PresetInfo
}

// PresetInfo preset info
type PresetInfo struct {
	Title      string `json:"title"`
	InputURL   string `json:"inputURL"`
	OutputPath string `json:"outputPath"`
	OutputFile string `json:"outputFile"`
	Duration   int    `json:"duration"`
	Video      bool   `json:"video"`
}

// PrintData print data
func (p *PresetParser) PrintData() {
	fmt.Println("PresetInfo")
	jsonData, _ := json.MarshalIndent(p.InfoData, "", "\t")
	fmt.Println(string(jsonData))
}

// ParseParameter parse parameter
func (p *PresetParser) ParseParameter(input string, output string, duration int, title string, isVideo bool) {
	p.InfoData = PresetInfo{
		InputURL:   input,
		OutputPath: output,
		Duration:   duration,
		Title:      title,
		Video:      isVideo,
	}
	p.InfoData.OutputFile = p.getOutputFileName(p.InfoData)

}

// Parse parse telegram json
func (p *PresetParser) Parse() {
	if absPresetPath, valid := p.validate(p.InfoFile); valid {
		file, _ := ioutil.ReadFile(absPresetPath)
		_ = json.Unmarshal([]byte(file), &p.InfoData)
	}

	p.InfoData.OutputFile = p.getOutputFileName(p.InfoData)
	return
}

// FileHelper helper for file and directory
type FileHelper struct{}

func (h *FileHelper) validate(file string) (string, bool) {
	absPresetPath := h.getAbsPath(file)
	if !h.isExisted(absPresetPath, false) {
		fmt.Println("telegram info is not submitted")
		return "", false
	}

	return absPresetPath, true
}

func (h *FileHelper) isExisted(fp string, isMandatory bool) (ret bool) {
	ret = true
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		if isMandatory {
			h.createDirectory(fp)
		} else {
			ret = false
		}
	}
	return ret
}

func (h *FileHelper) createDirectory(path string) {
	err := os.Mkdir(path, 0755)
	if err != nil {
		panic(err)
	}
	fmt.Println("create output directory: ", path)
}

func (h *FileHelper) getAbsPath(fp string) string {
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

func (h *FileHelper) getNewFilename(absPath string, title string, isVideo bool) string {
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

func (h *FileHelper) getOutputFileName(pData PresetInfo) string {
	absPath := h.getAbsPath(pData.OutputPath)
	h.isExisted(absPath, true)

	outputFile := h.getNewFilename(absPath, pData.Title, pData.Video)
	for {
		if h.isExisted(outputFile, false) {
			fmt.Println("already existed: ", outputFile)
			outputFile = h.getNewFilename(absPath, pData.Title, pData.Video)
		} else {
			break
		}
	}

	return outputFile
}
