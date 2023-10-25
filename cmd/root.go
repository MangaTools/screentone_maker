package cmd

import (
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "screentone_maker",
	Short: "Эта программа накладывает скринтон на изображение. Входными изображениями могут быть файлы формата jpeg, jpg и png",
}

func init() {
	root.DisableSuggestions = true
	root.CompletionOptions.DisableDefaultCmd = true
	root.AddCommand(executeCmd)
}

func Execute() {
	root.Execute()
}
