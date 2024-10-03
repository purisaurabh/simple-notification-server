# Simple Notification Server

A Simple Notification Server Written in Go.
This was written to understand websockets.

## License

This code is licensed under the MIT License, for more details, checkout the [LICENSE](LICENSE)
file

# Architectur of the notification server

1. Server :

   - Server has two api
     - /subscribe
     - /listener

2. Client :
   - Client call the subscribe api
   - subscribe api generate the auth token
   - using the auth token call the listener api
   - listener api upgraded the websocket
   - websocket gives response to the client

# Note :

- Working on frontend
