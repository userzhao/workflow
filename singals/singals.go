package singals

import (
	"os/signal"
	"syscall"
	"fmt"
	"os"
)

var S = make(chan os.Signal)

func ListenSingal() {
	signal.Notify(S, syscall.SIGUSR2)
	for {
		s := <-S
		fmt.Println("get signal:", s)
	}
}