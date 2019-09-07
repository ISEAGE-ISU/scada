package security

import(
	"os"
	"errors"
	"io/ioutil"

	"github.com/ISEAGE-ISU/scada"
)

func Root(tok string) (scada.TokenFunc, error) {
	switch tok{
	case "start":
		return nil, Start()
	case "stop":
		return nil, Stop()
	case "power":
		os.Remove(onFile)
		os.Exit(0)
	case "check":
		f, err := os.Open(freeFile)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err 
		}
		return nil, errors.New(string(b) + " GB used")
	}
	return nil, scada.ErrUnknownTok
}
