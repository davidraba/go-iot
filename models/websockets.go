package models

type AnalogData struct {
	Capacity    float64 `json:"capacity"`
	Battery     float64 `json:"battery"`
	Temperature float64 `json:"temp"`
}

type WebsocketData struct {
	Timestamp int64      `json:"timestamp"`
	Analog    AnalogData `json:"analog"`
}
