package main

import (
	"fmt"
	"github.com/vistrcm/fun2bot/pts"
)

func main() {
	for i := 0; i < 10; i++ {
		strip, _ := pts.GetRandomImageURL()
		fmt.Println(strip)
	}

}
