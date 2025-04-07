package geo

import (
	"time"
)

type EventType int
type FlareClass int

const (
	XRAY_FLUX_CHANGED EventType = iota
	SOLAR_WIND_SPEED_CHANGED
	GROUND_FIELD_CHANGED
)

const (
	XRAY_FLARE_A FlareClass = iota
	XRAY_FLARE_B
	XRAY_FLARE_C
	XRAY_FLARE_M
	XRAY_FLARE_X
)

func (xf FlareClass) String() string {
	return [...]string{"Class A", "Class B", "Class C", "Class M", "Class X"}[xf]
}

func (et EventType) String() string {
	return [...]string{"Xray Flux Changed", "Solar Wind Speed Changed", "Ground Field Changed"}[et]
}

type GeoEvent struct {
	Time        time.Time
	TimeStart   time.Time
	TimeEnd     time.Time
	Event       EventType
	Class       FlareClass
	Value       float32
	Processed   bool
	Description string
}

type GeoEventReport struct {
	Events []GeoEvent
}
