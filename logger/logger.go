package logger

import (
	"fmt"
	"os"
	"strings"

	apexLog "github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/graylog"
)

// Graylog to Apex apexLog levels
var levelsDict = map[string]apexLog.Level{
	"EMERGENCY": apexLog.FatalLevel,
	"ALERT":     apexLog.FatalLevel,
	"CRITICAL":  apexLog.FatalLevel,
	"ERROR":     apexLog.ErrorLevel,
	"WARNING":   apexLog.WarnLevel,
	"NOTICE":    apexLog.InfoLevel,
	"INFO":      apexLog.InfoLevel,
	"DEBUG":     apexLog.DebugLevel,
}

// Log entry
var Log *apexLog.Entry

// Conf properties
type Conf struct {
	Host    string
	Port    string
	Level   string
	App     string
	Version string
	Handler apexLog.Handler
}

// Init the logger.Log.
func Init(conf Conf) error {
	apexLog.SetLevel(levelsDict[strings.ToUpper(conf.Level)])
	if conf.Host != "" && conf.Handler == nil {
		handler, err := graylog.New("udp://" + conf.Host + ":" + conf.Port)
		if err != nil {
			return fmt.Errorf("failed init graylog: %s", err)
		}
		apexLog.SetHandler(handler)
	} else if conf.Handler != nil {
		apexLog.SetHandler(conf.Handler)
	} else {
		apexLog.SetHandler(cli.New(os.Stdout))
	}
	Log = apexLog.WithFields(apexLog.Fields{
		"facility": conf.App,
		"version":  conf.Version,
	})
	return nil
}
