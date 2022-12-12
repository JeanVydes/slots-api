package main

import (
	"log"
	"github.com/gin-gonic/gin"
)

type GameI struct {
	ID      string      `json:"id"`
	Name	string      `json:"name"`
	Version string      `json:"version"`
	Start   func()      `json:"-"`
	Stop    func()      `json:"-"`
	Reload  func()      `json:"-"`
	SetHTTP func()      `json:"-"`
	APIPath string      `json:"api_path"`
	Data    interface{} `json:"data"`
	MaxBet  float64     `json:"max_bet"`
	MinBet  float64     `json:"min_bet"`
}

type Games struct {
	Games map[string]GameI
}

var gameRegistrationQueue = make(chan GameI)

func (g *Games) New() {
	g.SetHTTP()

	for {
		select {
		case game := <-gameRegistrationQueue:
			g.Add(game)

			if game.Start != nil {
				go game.Start()
				go game.SetHTTP()
			}

			gameChannel := Channel{
				ChannelID:   game.ID,
				Connections: make(map[string]*Client),
			}

			log.Println("New game registered: ", game.ID)
			ChannelRegistrationQueue <- gameChannel
		}
	}
}

func (g *Games) SetHTTP() {
	GamesRouter.GET("/all", func(c *gin.Context) {
		var games []GameI
		for _, game := range g.Games {
			games = append(games, GameI{
				ID:      game.ID,
				Name:    game.Name,
				Version: game.Version,
				APIPath: game.APIPath,
				Data:    game.Data,
				MaxBet:  game.MaxBet,
				MinBet:  game.MinBet,
			})
		}

		c.JSON(200, Message{
			ExitCode: 0,
			Message:  "OK",
			Data:     games,
		})
	})
}

func (g *Games) Add(game GameI) bool {
	g.Games[game.ID] = game
	return true
}

func (g *Games) Stop(gameID string) error {
	for _, game := range g.Games {
		if game.ID == gameID {
			if game.Stop != nil {
				game.Stop()
			}

			break
		}
	}

	return nil
}

func (g *Games) Get(gameID string) *GameI {
	game := g.Games[gameID]
	if game.ID == "" {
		return nil
	}

	return &game
}

func (g *GameI) Send(packetType string, data interface{}, broadcastType string, broadcastPrivateSockets map[string]*Client) {
	InternalChannelPacketChan <- WSPacket{
		Type:               packetType,
		Data:               data,
		Channel:            g.ID,
		BroadcastType:      broadcastType,
		SocketsToBroadcast: broadcastPrivateSockets,
	}
}