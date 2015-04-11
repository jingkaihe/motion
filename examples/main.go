package main

import (
	"fmt"
	"github.com/jaxi/motion"
	"log"
)

func main() {
	d, err := motion.NewDevice()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to leapmotion")

	d.ListenAndReceive(true)
	defer d.Close()

	for frame := range d.FrameQueue {
		if len(frame.Gestures) > 0 {
			for _, g := range frame.Gestures {
				fmt.Println(g.Type)
			}
		}
	}
}
