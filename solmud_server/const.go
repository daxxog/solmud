// const.go
package main

const (
	// Network constants
	LISTEN_HOST = "127.0.0.1" // Localhost only during development
	LISTEN_PORT = 43594       // Standard RuneScape 317 game port

	// Login handshake opcodes (317 revision)
	OPCODE_HANDSHAKE = 14
	OPCODE_NEW_LOGIN = 16 // Standard fresh login (most common)
	OPCODE_RECONNECT = 18 // Resume existing session (rarely used in basic 317 servers)

	// Login response codes
	LOGIN_RESPONSE_OK                    = 2
	LOGIN_RESPONSE_INVALID_CREDENTIALS   = 3
	LOGIN_RESPONSE_ACCOUNT_DISABLED      = 4
	LOGIN_RESPONSE_ALREADY_LOGGED_IN     = 5
	LOGIN_RESPONSE_GAME_UPDATED          = 6
	LOGIN_RESPONSE_WORLD_FULL            = 7
	LOGIN_RESPONSE_LOGIN_SERVER_REJECTED = 8
	LOGIN_RESPONSE_LOAD_ERROR            = 9
	LOGIN_RESPONSE_TOO_MANY_CONNECTIONS  = 10
	LOGIN_RESPONSE_INVALID_LOGIN_SERVER  = 11
	LOGIN_RESPONSE_BAD_SESSION_ID        = 12
)
