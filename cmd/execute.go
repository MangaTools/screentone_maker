package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/shadream/screentone_maker/executor"
	"github.com/spf13/cobra"
)

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "This command starts the process of overlaying a screentone on the images in the specified folder.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := executeSettings.Validate(); err != nil {
			return err
		}

		executor := executor.NewExecutor(executeSettings)

		return executor.ExecuteFolder(executeSettings.InputPath, executeSettings.OutPath, executeSettings.Recursive)
	},
}

var executePipeCmd = &cobra.Command{
	Use:   "pipe",
	Short: "This command starts processing of overlaying a screentone of a single file from the stdin stream and outputs to stdout stream.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := executeSettings.ValidatePipe(); err != nil {
			return err
		}

		executor := executor.NewExecutor(executeSettings)

		data, err := executor.Execute(os.Stdin)
		if err != nil {
			return err
		}

		fmt.Fprint(os.Stdout, data)

		return nil
	},
}

var executeSettings = executor.ExecutionSettings{}

func init() {
	executeCmd.AddCommand(executePipeCmd)

	executeCmd.PersistentFlags().UintVarP(&executeSettings.DotSize, "dot_size", "d", 5,
		"Max dot size in pixels. Value must be in range [2;100]. This flag is responsible for how big your screentone will be.")
	executeCmd.PersistentFlags().UintVarP(&executeSettings.Black, "black", "b", 1,
		"The maximum pixel color at which any pixel will be black. Flag only includes values less than the specified value.")
	executeCmd.PersistentFlags().UintVarP(&executeSettings.White, "white", "w", 255,
		"The minimum pixel color at which any pixel will be white. Flag includes values greater or equal than the specified value.")

	executeCmd.Flags().UintVarP(&executeSettings.Threads, "threads", "t", uint(runtime.NumCPU()),
		"Number of simultaneously processed images. The default value is equal to the number of logical processor cores.")
	executeCmd.Flags().StringVarP(&executeSettings.InputPath, "input", "i", "",
		"The path to the folder with the images.")
	executeCmd.Flags().BoolVarP(&executeSettings.Recursive, "recursive", "r", false,
		"Recursive image search.")
	executeCmd.Flags().StringVarP(&executeSettings.OutPath, "out", "o", "",
		`Path to the folder with final images (if the folder does not exist, it will be created).
	If not specified, it is written to the input folder. The output will be .png files. If there are already files with the same name in the folder, the files are overwritten.`)
}
