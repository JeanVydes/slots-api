package main

import (
	"fmt"
	"time"
)

var (
	RequestChannelAccessChan = make(chan WSPacket)
	GamePingChan             = make(chan WSPacket)
	GameInteractionChan      = make(chan WSPacket)
)

type PacketsManagerI struct{}

func newPacketsManager() PacketsManagerI {
	ChannelsManager = newChannnelsManager()
	return PacketsManagerI{}
}

func (pm *PacketsManagerI) HandleWSPackets() {
	for {
		select {
		case packet := <-RequestChannelAccessChan:
			data := packet.Data.(map[string]interface{})
			if data == nil {
				SendPacket(*packet.Client, WSPacket{
					Type: ChannelAccessDenied,
					Data: Map{
						"reason": "channel_id_key_not_included",
					},
				})
				continue
			}

			channelID := data["channel_id"].(string)
			channel := ChannelsManager.GetChannel(channelID)

			if channel.ChannelID == "" {
				SendPacket(*packet.Client, WSPacket{
					Type: ChannelAccessDenied,
					Data: map[string]interface{}{
						"reason": "channel_not_found",
					},
				})

				continue
			}

			if channel.Connections[packet.Client.userID] != nil {
				SendPacket(*packet.Client, WSPacket{
					Type: ChannelAccessDenied,
					Data: map[string]interface{}{
						"reason": "already_connected",
					},
				})

				continue
			}

			err := ChannelsManager.AddSocketToChannel(channel.ChannelID, packet.Client)
			if err != nil {
				SendPacket(*packet.Client, WSPacket{
					Type: ChannelAccessDenied,
					Data: Map{
						"reason": "channel_does_not_exist",
					},
				})

				continue
			}

			_ = ChannelsManager.RegisterChannelToClientLog(packet.Client, channel.ChannelID)

			SendPacket(*packet.Client, WSPacket{
				Type: ChannelAccessGranted,
				Data: Map{
					"channel_id": channel.ChannelID,
					"socket_id":  packet.Client.userID,
				},
			})

			go WSPingPong(*packet.Client, GamePing, WSPacket{
				Type: GamePing,
				Data: Map{
					"channel_id": channel.ChannelID,
				},
			}, 2*time.Second)
		case packet := <-GamePingChan:
			fmt.Println(packet)
		case packet := <-GameInteractionChan:
			fmt.Println(packet)
		}
	}
}
