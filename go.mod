module github.com/Laboratory-for-Safe-and-Secure-Systems/paho.mqtt.golang

go 1.23

toolchain go1.23.5

require (
	github.com/Laboratory-for-Safe-and-Secure-Systems/go-asl v1.0.3
	github.com/gorilla/websocket v1.5.3
	golang.org/x/net v0.27.0
	golang.org/x/sync v0.7.0
)

replace github.com/Laboratory-for-Safe-and-Secure-Systems/go-asl v1.0.3 => ../go-asl
