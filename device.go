package scada

import(
	"net"
	"net/http"
	"errors"
	"strings"
)

var ErrUnknownTok = errors.New("Unknown Token")

type TokenFunc func(tok string) (TokenFunc, error)

type checkOK func() bool

type Status func() string

type Device struct {
	Root TokenFunc
	Status Status
	CheckOK checkOK
	Con net.Listener
}

func (d *Device) StartSCADA() error {
	b := make([]byte, 512)
	if d.Con == nil {
		var err error
		d.Con, err = net.Listen("tcp", ":1337")
		if err != nil {
			return err
		}
	}
	for {
		c, err := d.Con.Accept()
		if err != nil {
			return err
		}
		go func() {
			for {
				_, err := c.Read(b)
				if err != nil {
					c.Close()
					return
				}
				tokFunc := d.Root
				for _, s := range strings.Fields(string(b)) {
					tokFunc, err = tokFunc(s)
					if err != nil {
						c.Write([]byte(err.Error() + "\n"))
						c.Close()
						return
					}
					if tokFunc == nil {
						c.Close()
						return
					}
				}
			}
		}()
	}
}

func (d *Device) StartHTTP(port string) error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !d.CheckOK() {
			http.Error(w, "Just send the money", http.StatusPaymentRequired)
			return
		}
		w.Write([]byte(d.Status()))
	})

	return http.ListenAndServe(port, handler)
}