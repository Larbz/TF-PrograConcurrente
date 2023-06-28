package main

import (
	"fmt"
	"sync"
	"time"
)

type Player struct {
	Name   string
	Tokens int
}

func main() {
	// var teams []*Team;
	var wg sync.WaitGroup
	var wg2 sync.WaitGroup
	var wg1 sync.WaitGroup
	c1 := make(chan Player)
	c2 := make(chan Player)
	team := Player{Name: "Luis", Tokens: 30}
	team2 := Player{Name: "Juan", Tokens: 150}
	stop1 := make(chan bool)
	stop2 := make(chan bool)
	wg.Add(2)
	go mov(&team, &wg1, c1, stop1, &wg)
	go mov(&team2, &wg2, c2, stop2, &wg)
	wg.Wait()

}

func mov(player *Player, wg *sync.WaitGroup, chPlayer chan Player, stop chan bool, mainWg *sync.WaitGroup) {
	defer mainWg.Done()
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go rps(player, wg, chPlayer, stop)
		<-stop
		playerInfo := <-chPlayer
		fmt.Println(playerInfo.Name, playerInfo.Tokens)
	}

}

func rps(player *Player, wg *sync.WaitGroup, chPlayer chan Player, stop chan bool) {
	defer wg.Done()
	if player.Tokens == 20 {
		time.Sleep(2 * time.Second)
	}
	stop <- true
	player.Tokens -= 1
	chPlayer <- *player

}
