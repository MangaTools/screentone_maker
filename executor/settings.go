package executor

import (
	"runtime"

	validation "github.com/go-ozzo/ozzo-validation"
)

type ExecutionSettings struct {
	DotSize     uint `json:"dot_size"`
	ClusterSize uint `json:"cluster_size"`

	InputPath string `json:"input_path"`
	OutPath   string `json:"out_path"`

	Threads uint `json:"threads"`
}

func (s ExecutionSettings) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.DotSize, validation.Min(uint(2)), validation.Max(uint(100))),
		validation.Field(&s.ClusterSize, validation.In(uint(0), uint(2), uint(4), uint(8), uint(16))),
		validation.Field(&s.InputPath, validation.Required),
		validation.Field(&s.Threads, validation.Min(uint(1)), validation.Max(uint(runtime.NumCPU()))),
	)
}
