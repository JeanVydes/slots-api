package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type CoinFlipGame struct {
	Game GameI
}

func (this *CoinFlipGame) Register() {
	this.Game = GameI{
		ID:      "1000000",
		Name:    "Coin Flip",
		Version: "0.0.1beta",
		Start:   this.Start,
		Stop:    this.Stop,
		SetHTTP: this.SetHTTP,
		Reload:  func() {},
		APIPath: "/coinflip/0-0-1beta",
		Data:    nil,
		MaxBet:  5000,
		MinBet:  1,
	}

	gameRegistrationQueue <- this.Game
}

func (rg *CoinFlipGame) Start() {}

func (rg *CoinFlipGame) Stop() {}

func (rg *CoinFlipGame) Reload() {}

func (rg *CoinFlipGame) SetHTTP() {
	GamesRouter.GET(rg.Game.APIPath, func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"resource": rg.Game,
		})
	})

	GamesRouter.POST(rg.Game.APIPath, TokenMiddlware(), func(ctx *gin.Context) {
		side, sideQuery := ctx.GetQuery("side")
		betAmount, betQuery := ctx.GetQuery("bet_amount")

		if side == "" || !sideQuery {
			ctx.JSON(400, Message{
				ExitCode: 1,
				Message:  "missing.side.query",
				Data:     nil,
			})

			return
		}

		if betAmount == "" || !betQuery {
			ctx.JSON(400, Message{
				ExitCode: 1,
				Message:  "missing.bet_amount.query",
				Data:     nil,
			})

			return
		}

		if side != "heads" && side != "tails" {
			ctx.JSON(400, Message{
				ExitCode: 1,
				Message:  "invalid.side.query",
				Data:     nil,
			})

			return
		}

		betAmountFloat, err := strconv.ParseFloat(betAmount, 64)
		if err != nil {
			ctx.JSON(400, Message{
				ExitCode: 1,
				Message:  "invalid.bet_amount.query",
				Data:     nil,
			})

			return
		}

		if betAmountFloat <= rg.Game.MinBet || betAmountFloat > rg.Game.MaxBet {
			ctx.JSON(400, Message{
				ExitCode: 1,
				Message:  "bet.amount.higher.1.lower.5000",
				Data:     nil,
			})

			return
		}

		userID, existsUserID := ctx.Get("accountID")
		if !existsUserID {
			ctx.JSON(400, Message{
				ExitCode: 1,
				Message:  "reset.session",
				Data:     nil,
			})

			return
		}

		user, found := GetUserByID(userID.(string))
		if !found {
			ctx.JSON(400, Message{
				ExitCode: 1,
				Message:  "reset.session",
				Data:     nil,
			})

			return
		}

		if user.Balances.FIAT.USD < betAmountFloat {
			ctx.JSON(400, Message{
				ExitCode: 1,
				Message:  "not.enough.founds",
				Data:     nil,
			})

			return
		}

		multiplier := 0.45
		boolGenerator := NewBoolGenerator()

		// A random boolean generator
		// TRUE  = HEADS
		// FALSE = TAILS
		result := boolGenerator.Bool()
		var readableResult string
		if result {
			readableResult = "heads"
		} else {
			readableResult = "tails"
		}

		if (result && side == "heads") || (!result && side == "tails") {
			profit := betAmountFloat * multiplier
			newBalances := user.Balances
			newBalances.FIAT.USD += profit

			err := UpdateUserBalance(user.ID, newBalances)
			if err != nil {
				ctx.JSON(500, Message{
					ExitCode: 1,
					Message:  "database.error",
					Data:     nil,
				})
			}

			ctx.JSON(200, Message{
				ExitCode: 0,
				Message:  "success",
				Data: Map{
					"side_selected": side,
					"side_result":   readableResult,
					"initial_bet":   betAmountFloat,
					"profit":        profit,
					"multiplier":    multiplier,
					"total":         profit + betAmountFloat,
					"new_balances":  newBalances,
				},
			})

			return
		}

		newBalances := user.Balances
		newBalances.FIAT.USD -= betAmountFloat

		err = UpdateUserBalance(user.ID, newBalances)
		if err != nil {
			ctx.JSON(500, Message{
				ExitCode: 1,
				Message:  "database.error",
				Data:     nil,
			})
		}

		ctx.JSON(200, Message{
			ExitCode: 0,
			Message:  "success",
			Data: Map{
				"side_selected": side,
				"side_result":   readableResult,
				"initial_bet":   betAmountFloat,
				"profit":        0,
				"total":         -betAmountFloat,
				"new_balances":  newBalances,
			},
		})
	})
}
