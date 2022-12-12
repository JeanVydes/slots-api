package main

import (
	"fmt"
	"time"
)

func SendPacket(client Client, packet WSPacket) {
	client.send <- []byte(fmt.Sprint(packet))
}

func WSPingPong(client Client, packetType string, data interface{}, interval time.Duration) {
	quit := make(chan bool)
	ticker := time.NewTicker(PingPongDelay)
	defer func() {
		quit <- true
		ticker.Stop()
	}()

ping:
	for {
		select {
		case <-ticker.C:
			SendPacket(client, WSPacket{
				Type: packetType,
				Data: data,
			})
		case <-quit:
			ticker.Stop()
			break ping
		}
	}
}
