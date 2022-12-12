package main

var (
	ChannelsManager ChannelsManagerI

	InternalChannelPacketChan  = make(chan WSPacket)
	ChannelRegistrationQueue   = make(chan Channel)
	ChannelUnregistrationQueue = make(chan Channel)
)

type Channel struct {
	ChannelID   string
	Connections map[string]*Client
}

type ChannelsManagerI struct {
	Channels map[string]Channel
}
