package main

import (
	"errors"
)

// A channel is a way to manage the packets that will be sent to the clients in a broadcast way, but not everyone socket need that packet
// Example: A user is in the rocket game, so the server will send packets about rocket game, but no about other games

// This function will register all channels (games, etc)
func newChannnelsManager() ChannelsManagerI {
	ChannelsManager = ChannelsManagerI{
		Channels: make(map[string]Channel),
	}

	go ChannelsManager.NewPacketManager()
	go ChannelsManager.ChannelsRegistration()

	return ChannelsManager
}

// Add Socket to a list to get packets, like other players data
func (cm *ChannelsManagerI) AddSocketToChannel(channelID string, client *Client) error {
	if ChannelsManager.Channels[channelID].ChannelID == "" {
		return errors.New("Channel not found")
	}

	ChannelsManager.Channels[channelID].Connections[client.userID] = client

	return nil
}

func (cm *ChannelsManagerI) RemoveSocketFromChannel(channelID string, client *Client) error {
	if ChannelsManager.Channels[channelID].ChannelID == "" {
		return errors.New("Channel not found")
	}

	delete(ChannelsManager.Channels[channelID].Connections, client.userID)

	return nil
}

// This filter and broadcast by private broadcast (selected sockets) or a public broadcast (everybody), usually for multipurpose uses
func (cm *ChannelsManagerI) NewPacketManager() {
	for {
		select {
		case ICPacket := <-InternalChannelPacketChan:
			if ChannelsManager.Channels[ICPacket.Channel].ChannelID == "" {
				return
			}

			var socketsToBroadcast map[string]*Client
			if ICPacket.BroadcastType == ChannelPublicBroadcast {
				socketsToBroadcast = ChannelsManager.Channels[ICPacket.Channel].Connections
			} else if ICPacket.BroadcastType == ChannelPrivateBroadcast {
				socketsToBroadcast = ICPacket.SocketsToBroadcast
			}

			for _, client := range socketsToBroadcast {
				SendPacket(*client, ICPacket)
			}
		}
	}
}

func (cm *ChannelsManagerI) ChannelsRegistration() {
	for {
		select {
		case channel := <-ChannelRegistrationQueue:
			ChannelsManager.Channels[channel.ChannelID] = channel
		case channel := <-ChannelUnregistrationQueue:
			delete(ChannelsManager.Channels, channel.ChannelID)
		}
	}
}

func (cm *ChannelsManagerI) GetChannel(channelID string) Channel {
	return ChannelsManager.Channels[channelID]
}

func (cm *ChannelsManagerI) GetChannelSockets(channelID string) map[string]*Client {
	return ChannelsManager.Channels[channelID].Connections
}

func (cm *ChannelsManagerI) RegisterChannelToClientLog(client *Client, channelID string) error {
	for _, channel := range client.channels {
		if channel == channelID {
			return errors.New("Channel already registered")
		}
	}

	client.channels = append(client.channels, channelID)

	return nil
}

func (cm *ChannelsManagerI) UnregisterChannelToClientLog(client *Client, channelID string) []string {
	var newChannels []string
	for _, channel := range client.channels {
		if channel != channelID {
			newChannels = append(newChannels, channel)
		}
	}

	client.channels = newChannels

	return client.channels
}
