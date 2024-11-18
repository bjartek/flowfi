package bot

import (
	"fmt"

	"go.uber.org/zap"
)

type ZapBotLogger struct {
	Logger *zap.Logger
}

func (zbl ZapBotLogger) Println(v ...interface{}) {
	zbl.Logger.Debug(fmt.Sprint(v...))
}

func (zbl ZapBotLogger) Printf(format string, v ...interface{}) {
	zbl.Logger.Debug(fmt.Sprintf(format, v...))
}
