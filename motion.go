package motion

import (
	"golang.org/x/net/websocket"
	"log"
)

type TrackingDate struct {
	CurrentFrameRate float32        `json:"currentFrameRate"`
	Id               float32        `json:"id"`
	R                [][]float32    `json:"r"`
	S                float32        `json:"s"`
	T                []float32      `json:"t"`
	Timestamp        int            `json:"timestamp"`
	Gestures         []Gesture      `json:"gestures"`
	Hands            []Hand         `json:"hands"`
	InteractionBox   InteractionBox `json:"interactionBox"`
	Pointables       []Pointable    `json:"pointables"`
}

type Gesture struct {
	Center        []float32 `json:"center"`
	Direction     []float32 `json:"direction"`
	Duration      int       `json:"duration"`
	HandIds       []int     `json:"handIds"`
	Id            int       `json:"id"`
	Normal        []float32 `json:"normal"`
	PointableIds  []int     `json:"pointableIds"`
	Position      []float32 `json:"position"`
	Progress      float32   `json:"progress"`
	Radius        float32   `json:"radius"`
	Speed         float32   `json:"speed"`
	StartPosition []float32 `json:"startPosition"`
	State         string    `json:"state"`
	Type          string    `json:"type"`
}

type Hand struct {
	armBasis               []float32   `json:"armBasis"`
	armWidth               float32     `json:"armWidth"`
	Confidence             float32     `json:"confidence"`
	Direction              []float32   `json:"direction"`
	Elbow                  []float32   `json:"elbow"`
	GrabStrength           float32     `json:"grabStrength"`
	Id                     int         `json:"id"`
	PalmNormal             []float32   `json:"palmNormal"`
	PalmPosition           []float32   `json:"palmPosition"`
	PalmVelocity           []float32   `json:"palmVelocity"`
	PinchStrength          float32     `json:"pinchStrength"`
	R                      [][]float32 `json:"r"`
	S                      float32     `json:"s"`
	SphereCenter           []float32   `json:"sphereCenter"`
	SphereRadius           float32     `json:"sphereRadius"`
	StabilizedPalmPosition []float32   `json:"stabilizedPalmPosition"`
	T                      []float32   `json:"t"`
	TimeVisible            float32     `json:"timeVisible"`
	Type                   string      `json:"type"`
	Wrist                  []float32   `json:"wrist"`
}

type InteractionBox struct {
	Center []float32 `json:"center"`
	Size   []float32 `json:"size"`
}

type Pointable struct {
	Base                  []float32 `json:"bases"`
	BtipPosition          []float32 `json:"btipPosition"`
	CarpPosition          []float32 `json:"carpPosition"`
	DipPosition           []float32 `json:"dipPosition"`
	Direction             []float32 `json:"direction"`
	Extended              bool      `json:"extended"`
	HandId                int       `json:"handId"`
	Id                    int       `json:"id"`
	Length                float32   `json:"length"`
	McpPosition           []float32 `json:"mcpPosition"`
	PipPosition           []float32 `json:"pipPosition"`
	StabilizedTipPosition []float32 `json:"stabilizedTipPosition"`
	TimeVisible           float32   `json:"timeVisible"`
	TipPosition           []float32 `json:"tipPosition"`
	TipVelocity           []float32 `json:"tipVelocity"`
	Tool                  bool      `json:"tool"`
	TouchDistance         float32   `json:"touchDistance"`
	TouchZone             string    `json:"touchZone"`
	Type                  int       `json:"type"`
	Width                 float32   `json:"width"`
}

type Device struct {
	Ws         *websocket.Conn
	FrameQueue chan TrackingDate
}

func NewDevice() (*Device, error) {
	d := Device{}
	d.FrameQueue = make(chan TrackingDate)
	conn, err := websocket.Dial("ws://127.0.0.1:6437/v3.json", "", "http://localhost/")
	if err != nil {
		return &d, err
	}

	d.Ws = conn
	return &d, nil
}

func (d *Device) ListenAndReceive() error {

	enableGestures := struct {
		enableGestures bool `json:"enableGestures"`
	}{true}

	if err := websocket.JSON.Send(d.Ws, &enableGestures); err != nil {
		return err
	}

	backgroundMessage := struct {
		background bool `json:"background"`
	}{true}

	if err := websocket.JSON.Send(d.Ws, &backgroundMessage); err != nil {
		return err
	}

	go d.Receive()

	return nil
}

func (d *Device) Receive() {
	var data TrackingDate
	for {
		if err := websocket.JSON.Receive(d.Ws, &data); err != nil {
			log.Println(err)
			continue
		} else {
			d.FrameQueue <- data
		}
	}
	close(d.FrameQueue)
}

func (d *Device) Close() {
	d.Ws.Close()
}
