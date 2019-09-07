package generator

import(
	"os"
	"github.com/ISEAGE-ISU/scada"
)

func tokSet(level string) (scada.TokenFunc, error) {
	return nil, Set(level)
}

func Root(tok string) (scada.TokenFunc, error) {
	switch tok{
	case "start":
		return nil, Start()
	case "stop":
		return nil, Stop()
	case "set":
		return tokSet, nil
	case "power":
		os.Remove(onFile)
		os.Exit(0)
	}
	return nil, scada.ErrUnknownTok
}