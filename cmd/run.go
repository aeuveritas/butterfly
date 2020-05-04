/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/aeuveritas/butterfly/transcode"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run -i [INPUT_URL] -o [OUTPUT_PATH] -d [DURATION]",
	Short: "transcode media stream to file",
	Long:  "butterfly run -i INPUT_URL -o OUTPUT_DIRECTORY -d DURATION_IN_MIN -t TITLE",

	Run: func(cmd *cobra.Command, args []string) {
		inputURL, outputFile, durationString := transcode.ParseParameter(input, output, duration, title, video)

		transcode.Run(inputURL, outputFile, durationString, video, false)
	},
}

var input string
var output string
var duration int
var title string
var video bool

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.Flags().StringVarP(&input, "input", "i", "", "URL for input stream (required)")
	runCmd.Flags().StringVarP(&output, "output", "o", "", "directory for output (required)")
	runCmd.Flags().IntVarP(&duration, "duration", "d", 0, "duration in minutes (required)")
	runCmd.Flags().StringVarP(&title, "title", "t", "noname", "title for output file")
	runCmd.Flags().BoolVarP(&video, "video", "v", false, "audio or video (default: audio)")
	runCmd.MarkFlagRequired("input")
	runCmd.MarkFlagRequired("output")
	runCmd.MarkFlagRequired("duration")
}
