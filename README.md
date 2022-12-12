# Slots API

(This whole file has been explained by Chat GPT, may contain errors)

The API provides a way for developers to create and manage games that can be played by multiple users over the web. The use of websockets allows for real-time communication between the server and the clients, allowing for interactive gameplay.

The authentication protocol ensures that only authorized users can access the API and play the games. This is done by using a unique token, known as the "X-Auth-Token", that is associated with each user's session. This token is read on the server-side and used to verify the user's identity before allowing them to play a game.

To create a game, developers can use a provided template that allows them to specify the necessary details, such as the name and version of the game, the start function, and the HTTP methods for playing the game. Once the game functionality has been implemented, the games manager will take care of setting up the necessary endpoints and channels for the game to function properly.

In summary, the API provides a convenient and secure way for developers to create and manage games that can be played by multiple users in real-time over the web.

## Enviroment Variables

```
MONGO_URL=
MAIN_DATABASE_NAME=mydb
PORT=8080
```

## Interfaces

### Games

```
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
```

The GameI struct represents a game that can be managed by the API. It has several fields that contain information about the game, such as its ID, name, and version.

The ID field contains a unique identifier for the game. The Name and Version fields contain the name and version of the game, respectively.

The Start, Stop, and Reload fields contain functions that can be used to start, stop, and reload the game, respectively. The SetHTTP field contains a function that can be used to set up the necessary HTTP methods for playing the game.

The APIPath field contains the path to the game's API endpoint, which can be used by clients to interact with the game. The Data field can be used to store any additional data that the game needs to function properly.

The MaxBet and MinBet fields contain the maximum and minimum bets that are allowed for the game, respectively. These values can be used to limit the amount of money that players can bet on a game.

In summary, the GameI struct provides a way to store information about a game and its associated functions and data. It can be used to manage and interact with the game through the API.

#### How add a Game?

```
// Speicfy the game information
game = GameI{
		ID:      "1000000", // ID
		Name:    "Coin Flip", // A Name
		Version: "0.0.1beta", // A version
		Start:   something, // A Function to start the game
		Stop:    something, // A function to stop the game
		SetHTTP: something, // A function that create the game endpoints that manage the games
		Reload:  func() {}, // A Function that reload the game
		APIPath: "/coinflip/0-0-1beta", // API Path relative to /v1/games/
		Data:    nil, // Other information about the game
		MaxBet:  5000, // Max money that can be put in every single bet
		MinBet:  1, // Minimum money that can be put in every single bet
}

// Send the game to game registration queue
gameRegistrationQueue <- game
```

### User

```
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
```

The User struct represents a user in the system. It contains several fields that store information about the user, such as their ID, username, email, and birthday.

The Balances field contains information about the user's balances in different currencies, such as FIAT and crypto. The FIAT and Crypto structs contain information about the user's balances in different FIAT and crypto currencies, respectively.

The Birthday struct contains the user's birthday in timestamp format. The CreationDate field contains the date and time when the user was created in the system.

In summary, the User struct provides a way to store and manage information about users in the system, including their balances and other details. It can be used to interact with the user and their data in the system.

### Authentication

```
type Session struct {
	Token     string `json:"token"`
	AccountID string `json:"account_id"`
}
```

The Session struct represents a session in the system. It contains two fields: Token and AccountID.

The Token field contains a unique token that is associated with the session. This token can be used to identify and authenticate the user associated with the session.

The AccountID field contains the ID of the user's account that is associated with the session. This can be used to retrieve information about the user and their account in the system.

In summary, the Session struct provides a way to store and manage information about sessions in the system. It can be used to authenticate users and track their activities in the system.

### Channels

```
type Channel struct {
	ChannelID   string
	Connections map[string]*Client
}

type ChannelsManagerI struct {
	Channels map[string]Channel
}
```

A Channel in this context refers to a virtual channel that is used for communication between the server and the clients playing a particular game. It is essentially a direct connection between the server and the clients that allows for real-time interaction and communication during gameplay.

Each channel has a unique ChannelID that is used to identify it in the system. The Connections field contains a map of all the clients that are connected to the channel, with each client represented by a Client struct.

The ChannelsManagerI struct provides a way to manage the channels in the system. It contains a map of all the channels, with each channel identified by its ChannelID.

In summary, channels provide a way for the server and the clients to communicate and interact in real-time during gameplay. The Channel and ChannelsManagerI structs provide a way to manage and access these channels in the system.

### Websocket

```
type WSPacket struct {
	Type               string      `json:"type"`
	Data               interface{} `json:"data"`
	Client             *Client     `json:"-"`
	Channel            string
	BroadcastType      string
	SocketsToBroadcast map[string]*Client
}
```

The WSPacket struct represents a packet of data that is sent and received over a websocket connection. It contains several fields that provide information about the packet and its contents.

The Type field specifies the type of the packet, which can be used to identify its purpose and contents. The Data field contains the actual data that is being sent or received in the packet.

The Client field contains a reference to the Client struct that represents the client that sent or received the packet. The Channel field specifies the channel that the packet was sent on, if applicable.

The BroadcastType field specifies the type of broadcast that should be performed with the packet, if any. The SocketsToBroadcast field contains a map of all the clients that the packet should be broadcasted to, if applicable.

In summary, the WSPacket struct provides a way to package and manage data that is sent and received over a websocket connection. It contains information about the packet and its contents, as well as details about the client and the channel it was sent on.



```
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
```
The PingPongDelay variable specifies the interval at which the server should send ping messages to the clients to keep the websocket connection alive. This ensures that the connection remains active and responsive even if no data is being sent over the websocket.

The ChannelPublicBroadcast and ChannelPrivateBroadcast variables specify the types of broadcasts that can be performed on a channel. A public broadcast will be sent to all clients connected to the channel, while a private broadcast will only be sent to a selected group of clients.

The Ping, ChannelAccessDenied, ChannelAccessGranted, RequestChannelAccess, GamePing, and GameInteraction variables contain the string values of different types of packets that can be sent over the websocket. For example, the Ping value can be used to identify ping packets, while the GameInteraction value can be used to identify packets containing game interactions.

In summary, the variables provide constants and values that can be used to manage and identify the different types of packets and broadcasts that can be performed over the websocket.
