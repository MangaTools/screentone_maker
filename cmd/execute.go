package cmd

import (
	"runtime"

	"github.com/shadream/screentone_maker/executor"
	"github.com/spf13/cobra"
)

var execute = &cobra.Command{
	Use:   "execute",
	Short: "Данная команда запускает процесс выполнения наложение скринтона на изображение в указанной папке.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := executeSettings.Validate(); err != nil {
			return err
		}

		return executor.RunExecution(executeSettings)
	},
}

var executeSettings = executor.ExecutionSettings{}

func init() {
	execute.Flags().UintVarP(&executeSettings.DotSize, "dot_size", "d", 5,
		"Максимальный размер точки в пикселях. Значение должно быть в отрезке [2;100]")
	execute.Flags().UintVarP(&executeSettings.ClusterSize, "cluster_size", "c", 0,
		"Размер матрицы из точек(точки не равно пиксели). Необходимо для dithering. Выберите одно из значений: 0(выключено), 2, 4, 8, 16.")
	execute.Flags().UintVarP(&executeSettings.Threads, "threads", "t", uint(runtime.NumCPU()),
		"Количество одновременно обрабатываемых изображений. По умолчанию равно количеству логических ядер процессора.")
	execute.Flags().UintVarP(&executeSettings.Black, "black", "b", 1,
		"Максимальный цвет пикселя, при котором любой пиксель будет черным.")
	execute.Flags().UintVarP(&executeSettings.White, "white", "w", 255,
		"Минимальный цвет пикселя, при котором любой пиксель будет белым.")
	execute.Flags().StringVarP(&executeSettings.InputPath, "input", "i", "",
		"Путь до папки с изображениями.")

	execute.Flags().StringVarP(&executeSettings.OutPath, "out", "o", "",
		"Путь до папки с итоговыми изображениями (если папки не существует, то она создастся).\n"+
			"Если не указано, то записывается в input папку. На выходе будут .png файлы. Если в папке уже есть файлы с таким же названием, то файлы перезаписываются.")
}
