package server

type AppState map[string]interface{}

func NewAppState() AppState {
	return AppState(make(map[string]interface{}))
}
