package transcode

import (
	"fmt"
	"strconv"

	"github.com/aeuveritas/butterfly/parser"
	"github.com/xfrr/goffmpeg/transcoder"
)

// Run transcode media stream to file
func Run(data parser.PresetInfo, isDebug bool) {
	// Create new instance of transcoder
	trans := new(transcoder.Transcoder)

	// Initialize transcoder passing the input url and output file path
	err := trans.Initialize(data.InputURL, data.OutputFile)
	if err != nil {
		panic(err)
	}

	durationString := getDurationString(data.Duration, isDebug)
	trans.MediaFile().SetDuration(durationString)
	trans.MediaFile().SetOutputPath(data.OutputFile)
	if data.Video {
		trans.MediaFile().SetVideoCodec("libx264")
	} else {
		trans.MediaFile().SetSkipVideo(true)
	}

	// Start transcoder process without checking progress
	fmt.Println("new file: ", data.OutputFile)
	done := trans.Run(false)

	// This channel is used to wait for the process to end
	err = <-done
	if err != nil {
		panic(err)
	}
}

func getDurationString(duration int, isDebug bool) string {
	var weight int = 1
	if !isDebug {
		weight = 60
	}
	return strconv.Itoa(duration * weight)
}
