package log

import (
	kitlog "github.com/go-kit/kit/log"
	"github.com/hyperledger/burrow/logging/loggers"
	logging_types "github.com/hyperledger/burrow/logging/types"
)

func NewBurrowInfoTraceLogger() logging_types.InfoTraceLogger {
	return loggers.NewInfoTraceLogger(nil, kitlog.NewLogfmtLogger(kitlog.StdlibWriter{}))
}
