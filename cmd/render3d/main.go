package main

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	fmt.Println("RuneScape 2 Build 317 - Ebitengine v2 Test")

	const screen_width = 512
	const screen_height = 384

	ebiten.SetWindowSize(screen_width, screen_height)
	ebiten.SetWindowTitle("RS2 - Ebitengine v2 Test")

	start := time.Now()

	fmt.Println("Running for 3 seconds...")

	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Second)
		fmt.Printf("Tick %d\n", i+1)
		ebiten.Update()
	}

	fmt.Println("Test complete!")

	fmt.Printf("Ran for %v\n", time.Since(start))
}
