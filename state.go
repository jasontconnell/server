package server

type ServerState map[string]interface{}

func NewServerState() ServerState {
	return ServerState(make(map[string]interface{}))
}
