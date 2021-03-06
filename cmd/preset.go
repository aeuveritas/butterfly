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
	"github.com/aeuveritas/butterfly/parser"
	"github.com/aeuveritas/butterfly/transcode"
	"github.com/spf13/cobra"
)

// presetCmd represents the preset command
var presetCmd = &cobra.Command{
	Use:   "preset [JSON_FOR_PRESET]",
	Short: "transcode media stream to file with preset",
	Long:  "butterfly preset JSON_FOR_PRESET",

	Run: func(cmd *cobra.Command, args []string) {
		runPreset(args)
	},
}

var preset string

func init() {
	rootCmd.AddCommand(presetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// presetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// presetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runPreset(args []string) {
	if len(args) != 1 {
		panic("one input json absolute path for preset required")
	} else {
		isDebug := false
		var tParser parser.TelegramParser = parser.TelegramParser{
			InfoFile: "./telegram.json",
		}
		tParser.Parse()
		var pParser parser.PresetParser = parser.PresetParser{
			InfoFile: args[0],
		}
		pParser.Parse()

		bot := notification.GetTelegramBot(tParser.InfoData.Token, isDebug)
		chatID := notification.GetChatID(&tParser.InfoData, pParser.InfoData)
		if chatID != "" {
			tParser.SaveTelegramInfo()
		}

		notification.SendText(bot, chatID, pParser.InfoData.OutputFile, isDebug)
		transcode.Run(pParser.InfoData, isDebug)
		if video {
			notification.SendVideo(bot, chatID, pParser.InfoData.OutputFile, isDebug)
		} else {
			notification.SendAudio(bot, chatID, pParser.InfoData.OutputFile, isDebug)
		}

	}
}
