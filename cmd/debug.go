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
	"github.com/aeuveritas/butterfly/notification"
	"github.com/aeuveritas/butterfly/transcode"
	"github.com/spf13/cobra"
)

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:   "debug [JSON_FOR_PRESET]",
	Short: "transcode media stream to file with preset (debug)",
	Long:  "butterfly debug JSON_FOR_PRESET",

	Run: func(cmd *cobra.Command, args []string) {
		runDebug(args)
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// debugCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// debugCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runDebug(args []string) {
	if len(args) != 1 {
		panic("one input json absolute path for preset required")
	} else {
		isDebug := true
		inputURL, outputFile, durationString, token, video := transcode.ParsePreset(args[0], isDebug)
		chatID, bot := notification.GetTelegramObject(token, isDebug)

		notification.SendText(bot, chatID, outputFile, isDebug)
		transcode.Run(inputURL, outputFile, durationString, video, isDebug)
		if video {
			notification.SendVideo(bot, chatID, outputFile, isDebug)
		} else {
			notification.SendAudio(bot, chatID, outputFile, isDebug)
		}
	}
}