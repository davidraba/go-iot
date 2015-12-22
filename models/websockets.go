package models

type SiloData struct {
	Distance    float64 `json:"distance"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

type AnalogData struct {
	Percentage    float64 `json:"percentage"`
	Capacity      float64 `json:"capacity"`
	WeightCurrent float64 `json:"weight"`
	VolumeCurrent float64 `json:"volume"`
}

type WebsocketData struct {
	Timestamp int64      `json:"timestamp"`
	Analog    AnalogData `json:"analog"`
}
