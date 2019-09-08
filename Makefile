PREFIX := /usr/local/bin

scada: cmd/scada/main.go
	cd cmd/scada/ &&\
	go build -v 

install: scada
	cp cmd/scada/scada $(PREFIX)/scada
	cp rc.scada /etc/rc.d/scada

