package factory

import (
	controller2 "process-user-transaction/internal/adapters/inbound/s3"
	"process-user-transaction/internal/core/controller"
	"process-user-transaction/internal/core/parsers"
	"process-user-transaction/internal/core/service"
)

type Factory struct {
	service controller.Service
}

func NewFactory() *Factory {
	contentParsers := parsers.NewContentParsers(
		parsers.NewTxtContentParser(),
		parsers.NewXmlContentParser(),
		parsers.NewJpgContentParser(),
	)

	d := service.NewService()

	s := controller2.NewController(d, contentParsers)
	return &Factory{service: s}
}

func (f *Factory) Start(inputDir, outputDir string) error {
	err := f.service.Handle(inputDir, outputDir)
	if err != nil {
		return err
	}
	return nil
}
