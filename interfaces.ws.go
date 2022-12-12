package main

import (
	"time"
)

type WSPacket struct {
	Type               string      `json:"type"`
	Data               interface{} `json:"data"`
	Client             *Client     `json:"-"`
	Channel            string
	BroadcastType      string
	SocketsToBroadcast map[string]*Client
}

var (
	PingPongDelay           = time.Second * 5
	ChannelPublicBroadcast  = "public_broadcast"
	ChannelPrivateBroadcast = "selected_broadcast"
)

var (
	Ping                 = "ping"
	ChannelAccessDenied  = "channel_access_denied"
	ChannelAccessGranted = "channel_access_granted"
	RequestChannelAccess = "request_channel_access"
	GamePing             = "game_ping"
	GameInteraction      = "game_interaction"
)
