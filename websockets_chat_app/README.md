Just a simple program to demonstrate how one can implement Websockets in Go albeit implemented rather crudely
To test this out, just run main.go using the following command 
go run main.go

Next, to spin up 2 actors, say Alice & Bob, install "websocat" if needed.
brew install websocat

Run the following commands from different windows/terminal sessions.
websocat ws://localhost:9000/chat -H "username: alice" -H "friends: bob" -H "Sec-WebSocket-Version: 13" -H "Sec-WebSocket-Key: AQIDBAUGBwgJCgsMDQ4PEA=="
websocat ws://localhost:9000/chat -H "username: bob" -H "friends: alice" -H "Sec-WebSocket-Version: 13" -H "Sec-WebSocket-Key: AQIDBAUGBwgJCgsMDQ4PEA=="
