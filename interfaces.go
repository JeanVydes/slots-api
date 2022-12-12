package main

type User struct {
	ID           string   `json:"id" bson:"_id"`
	Username     string   `json:"username" bson:"username"`
	Email        string   `json:"email" bson:"email"`
	Birthday     Birthday `json:"birthday" bson:"birthday"`
	CreationDate int64    `json:"creation_date" bson:"creation_date"`
	Password     []byte   `json:"password" bson:"password"`
	Balances     Balances `json:"balances" bson:"balances"`
}

type Balances struct {
	FIAT   FIAT   `json:"fiat" bson:"fiat"`
	Crypto Crypto `json:"crypto" bson:"crypto"`
}

type FIAT struct {
	USD float64 `json:"usd" bson:"usd"`
}

type Crypto struct {
	BTC float64 `json:"btc" bson:"btc"`
	ETH float64 `json:"eth" bson:"eth"`
}

type Birthday struct {
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
}

type Message struct {
	ExitCode int         `json:"exit_code"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
}

type Session struct {
	Token     string `json:"token"`
	AccountID string `json:"account_id"`
}
