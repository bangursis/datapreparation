package logger

import "go.uber.org/zap"

func InitZap(logfile string) (*zap.Logger, error) {
	conf := zap.NewProductionConfig()

	conf.OutputPaths = []string{logfile, "stdout"}

	logger, err := conf.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
