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
	Short: "Данная команда запускает процесс выполнения наложение скринтона на изображение в указанной папке.",
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
	Short: "Данная команда запускает обработку одного файла из stdin потока и выводит в stdout.",
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
		"Максимальный размер точки в пикселях. Значение должно быть в отрезке [2;100]")
	executeCmd.PersistentFlags().UintVarP(&executeSettings.ClusterSize, "cluster_size", "c", 0,
		"Размер матрицы из точек(точки не равно пиксели). Необходимо для dithering. Выберите одно из значений: 0(выключено), 2, 4, 8, 16.")
	executeCmd.PersistentFlags().UintVarP(&executeSettings.Black, "black", "b", 1,
		"Максимальный цвет пикселя, при котором любой пиксель будет черным.")
	executeCmd.PersistentFlags().UintVarP(&executeSettings.White, "white", "w", 255,
		"Минимальный цвет пикселя, при котором любой пиксель будет белым.")

	executeCmd.Flags().UintVarP(&executeSettings.Threads, "threads", "t", uint(runtime.NumCPU()),
		"Количество одновременно обрабатываемых изображений. По умолчанию равно количеству логических ядер процессора.")
	executeCmd.Flags().StringVarP(&executeSettings.InputPath, "input", "i", "",
		"Путь до папки с изображениями.")
	executeCmd.Flags().BoolVarP(&executeSettings.Recursive, "recursive", "r", false,
		"Вложенный поиск файлов.")
	executeCmd.Flags().StringVarP(&executeSettings.OutPath, "out", "o", "",
		"Путь до папки с итоговыми изображениями (если папки не существует, то она создастся).\n"+
			"Если не указано, то записывается в input папку. На выходе будут .png файлы. Если в папке уже есть файлы с таким же названием, то файлы перезаписываются.")
}
