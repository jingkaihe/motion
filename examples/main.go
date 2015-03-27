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

	d.ListenAndReceive()
	defer d.Close()

	for msg := range d.FrameQueue {
		if len(msg.Gestures) > 0 {
			for _, g := range msg.Gestures {
				fmt.Println(g.Type)
			}
		}
	}
}
