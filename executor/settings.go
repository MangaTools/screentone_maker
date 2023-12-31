package executor

import (
	"runtime"

	validation "github.com/go-ozzo/ozzo-validation"
)

type ExecutionSettings struct {
	DotSize uint `json:"dot_size"`

	InputPath string `json:"input_path"`
	OutPath   string `json:"out_path"`
	Recursive bool   `json:"recursive"`

	Black uint `json:"black"`
	White uint `json:"white"`

	Threads uint `json:"threads"`
}

func (s ExecutionSettings) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.DotSize, validation.Min(uint(2)), validation.Max(uint(100))),
		validation.Field(&s.InputPath, validation.Required),
		validation.Field(&s.Threads, validation.Min(uint(1)), validation.Max(uint(runtime.NumCPU()))),
		validation.Field(&s.Black, validation.Min(uint(1)), validation.Max(s.White-1)),
		validation.Field(&s.White, validation.Min(uint(s.Black+1)), validation.Max(uint(255))),
	)
}

func (s ExecutionSettings) ValidatePipe() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.DotSize, validation.Min(uint(2)), validation.Max(uint(100))),
		validation.Field(&s.Black, validation.Min(uint(1)), validation.Max(s.White-1)),
		validation.Field(&s.White, validation.Min(uint(s.Black+1)), validation.Max(uint(255))),
	)
}
