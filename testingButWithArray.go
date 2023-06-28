package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Player struct {
	Name      string
	Tokens    int
	positionX int
	positionY int
	index     int
	inGame    bool
}

var wgGroup []*sync.WaitGroup
var chGroup []chan Player
var teams []*Player
var stopGroup []chan bool
var matriz [][]int

func main() {
	team1 := Player{Name: "Luis", Tokens: 3, positionX: 1, positionY: 1, index: 1, inGame: true}
	team2 := Player{Name: "Juan", Tokens: 3, positionX: 1, positionY: 2, index: 2, inGame: true}
	team3 := Player{Name: "Josue", Tokens: 3, positionX: 1, positionY: 3, index: 3, inGame: true}
	team4 := Player{Name: "Mario", Tokens: 3, positionX: 1, positionY: 4, index: 4, inGame: true}
	generatePlayer(team1)
	generatePlayer(team2)
	generatePlayer(team3)
	generatePlayer(team4)
	// numTeams := len(teams)
	// matriz = make([][]int, numTeams)
	// for i := range matriz {
	// 	matriz[i] = make([]int, 10)
	// 	matriz[i][0] = i
	// }
	// for _, fila := range matriz {
	// 	for _, valor := range fila {
	// 		fmt.Printf("%d ", valor)
	// 	}
	// 	fmt.Println()
	// }

	var wg sync.WaitGroup
	for ind := range teams {
		wg.Add(1)
		go manage(teams[ind], wgGroup[ind], chGroup[ind], stopGroup[ind], &wg)
	}
	wg.Wait()

}

func generatePlayer(player Player) {
	wgGroup = append(wgGroup, &sync.WaitGroup{})
	ch := make(chan Player)
	chGroup = append(chGroup, ch)
	stop := make(chan bool)
	stopGroup = append(stopGroup, stop)
	teams = append(teams, &player)
}

func manage(player *Player, wg *sync.WaitGroup, chPlayer chan Player, stop chan bool, mainWg *sync.WaitGroup) {
	defer mainWg.Done()
	// for i := 0; i < 20; i++ {
	for !isGameFinished() {
		if player.inGame {

			wg.Add(1)
			// go rps(player, wg, chPlayer, stop)
			go move(player, wg, chPlayer, stop)
			// if player.positionX == 10 {
			// 	aleatorio := rand.Intn(2)
			// 	if aleatorio == 0 {
			// 		player.positionY -= 1
			// 	} else {
			// 		player.positionY += 1
			// 	}
			// }
			// <-stop
			playerInfo := <-chPlayer
			<-stop
			fmt.Printf("%s %d %d %d\n", playerInfo.Name, playerInfo.Tokens, playerInfo.positionX, playerInfo.positionY)
		}
	}
}

func move(player *Player, wg *sync.WaitGroup, chPlayer chan Player, stop chan bool) {
	defer wg.Done()
	wg.Add(1)
	time.Sleep(1 * time.Second)
	if player.positionX == 10 {
		aleatorio := rand.Intn(2)
		var multiplier int
		if aleatorio == 0 {
			multiplier = -1
		} else {
			multiplier = 1
		}
		if aleatorio == 0 {
			player.positionY += 1 * multiplier
			player.positionX -= 1
		} else {
			player.positionY += 1 * multiplier
			player.positionX -= 1
		}
		if player.positionY < 1 {
			player.positionY = len(teams)
		} else if player.positionY > len(teams) {
			player.positionY %= len(teams)
		}
		for !teams[player.positionY-1].inGame && player.positionY != player.index {
			player.positionY += 1 * multiplier
			if player.positionY < 1 {
				player.positionY = len(teams)
			} else if player.positionY > len(teams) {
				player.positionY %= len(teams)
			}
		}

	} else if player.positionX == 1 && player.positionY != player.index {
		// player.Tokens += 1
		// teams[player.positionY-1].Tokens -= 1
		player.positionX = 1
		player.positionY = player.index
	} else {
		if player.positionY != player.index {
			player.positionX -= 1
		} else {
			player.positionX += 1
		}
	}
	go collisions(player, wg, chPlayer, stop)

	// validateColition()
	// stopGroup[player.index-1]<-true
	chPlayer <- *player
}

func collisions(player *Player, wg *sync.WaitGroup, chPlayer chan Player, stop chan bool) {
	defer wg.Done()
	// for i := 0; i < len(teams)-1; i++ {
	// 	for j := i + 1; j < len(teams); j++ {
	// 		if teams[i].positionY == teams[j].positionY {
	// 			if teams[i].positionX-teams[j].positionX >= -1 && teams[i].positionX-teams[j].positionX <= 1 {
	// 				time.Sleep(10*time.Second)
	// 			}

	// 		}
	// 	}
	// }
	for ind := range teams {
		if ind != player.index-1 {
			if player.positionY == teams[ind].positionY && (teams[ind].positionX-player.positionX == 1 || teams[ind].positionX-player.positionX == -1 || teams[ind].positionX-player.positionX == 0) && teams[ind].inGame {
				// time.Sleep(5 * time.Second)
				// player.Tokens += 1
				// teams[ind].Tokens -= 1
				// if teams[ind].Tokens <= 0 && teams[ind].inGame {
					result := playRPS()
					for result == "Tie" {
						result = playRPS()
					}
					if result == "Win" {
						player.Tokens += 1
						teams[ind].Tokens -= 1
						fmt.Printf("%s lose against %s\n", teams[ind].Name, player.Name)
						if teams[ind].Tokens == 0{
							teams[ind].inGame = false
							teams[ind].positionX=1
							teams[ind].positionY=teams[ind].index
							fmt.Printf("%s fue eliminado del juego", teams[ind].Name)
						}else{
							teams[ind].positionX=1
							teams[ind].positionY=teams[ind].index
						}

					} else {
						player.Tokens -= 1
						teams[ind].Tokens += 1
						fmt.Printf("%s lose against %s\n", player.Name, teams[ind].Name)
						if player.Tokens==0{
							player.inGame=false
							player.positionX=1
							player.positionY=player.index
							fmt.Printf("%s fue eliminado del juego", player.Name)
						}else{
							player.positionX=1
							player.positionY=player.index
						}

					}
				// }
			}
		}
	}

	stop <- true
	chPlayer <- *player

}

func isGameFinished() bool {
	playersLost := 0
	for _, player := range teams {
		if player.inGame == false {
			playersLost += 1
		}
	}
	var isFinished bool
	if len(teams)-playersLost == 1 {
		isFinished = true
	} else {
		isFinished = false
	}
	return isFinished
}

func playRPS() string {
	rpsOptions := []string{"Rock", "Paper", "Scissors"}
	playerChoice := rpsOptions[rand.Intn(len(rpsOptions))]
	otherPlayerChoice := rpsOptions[rand.Intn(len(rpsOptions))]

	fmt.Printf("You chose %s. The other player chose %s.\n", playerChoice, otherPlayerChoice)

	switch {
	case playerChoice == otherPlayerChoice:
		return "Tie"
	case playerChoice == "Rock" && otherPlayerChoice == "Scissors":
		return "Win"
	case playerChoice == "Paper" && otherPlayerChoice == "Rock":
		return "Win"
	case playerChoice == "Scissors" && otherPlayerChoice == "Paper":
		return "Win"
	default:
		return "Loss"
	}
}

// func rps(player1 *Player, player2 *Player, wg1 *sync.WaitGroup, wg2 *sync.WaitGroup, chPlayer1 chan Player, chPlayer2 chan Player, stop1 chan bool, stop2 chan bool) {
// 	defer wg1.Done()
// 	defer wg2.Done()

// 	time.Sleep(2 * time.Second)

// 	player1.positionX=1
// 	player1.positionY=player1.index

// 	stop1 <- true
// 	stop2 <- true
// 	// player1.Tokens -= 1
// 	// player2.Tokens -= 1
// 	// chPlayer1 <- *player1
// 	// chPlayer2 <- *player2
// 	// if player.Tokens == 20 {
// 	// 	time.Sleep(2 * time.Second)
// 	// }
// 	// stop <- true
// 	// player.Tokens -= 1
// 	// chPlayer <- *player

// }

// func validateColition() {
// 	// defer wg.Done()
// 	for i := 0; i < len(teams)-1; i++ {
// 		for j := i + 1; j < len(teams); j++ {
// 			if teams[i].positionY == teams[j].positionY {

// 				// fmt.Printf("%d %d %s %s\n",teams[i].positionY,teams[j].positionY,teams[i].Name,teams[j].Name)
// 				if teams[i].positionX-teams[j].positionX >= -1 && teams[i].positionX-teams[j].positionX <= 1 {
// 					wgGroup[i].Add(1)
// 					wgGroup[j].Add(1)
// 					go rps(teams[i], teams[j], wgGroup[i], wgGroup[j], chGroup[i], chGroup[j], stopGroup[i], stopGroup[j])
// 					<-stopGroup[i]
// 					<-stopGroup[j]
// 					// fmt.Printf("Colision entre %s y %s\n",teams[i].Name,teams[j].Name)
// 					// fmt.Printf("%s %d %d\n",teams[i].Name,teams[i].positionX,teams[i].positionY)
// 					// fmt.Printf("%s %d %d\n",teams[j].Name,teams[j].positionX,teams[j].positionY)

// 				}

// 			}
// 		}
// 	}
// }
