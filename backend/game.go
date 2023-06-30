package main

import (
	"bufio"
	"encoding/json"
	"log"
	"sort"

	// "errors"
	"fmt"
	// "log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	// "github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
	"github.com/gorilla/websocket"
	// socketio "github.com/googollee/go-socket.io"
	// "github.com/rs/cors"
)

// type Team struct {
// 	Name   string `json:"Name"`
// 	Tokens int    `json:"Tokens"`
// }

// type Mensaje struct {
// 	Numero int
// }

type Player struct {
	Name       string
	Tokens     int
	positionX  int
	positionY  int
	index      int
	inGame     bool
	delay      int
	freezed    bool
	tokensWon  int
	tokensLost int
	rpcWon     int
}

type PlayerPointsTable struct {
	PlayerName       string `json:"playerName"`
	PlayerPoints     int    `json:"playerPoints"`
	PlayerPositionX  int    `json:"playerPositionX"`
	PlayerPositionY  int    `json:"playerPositionY"`
	PlayerTokensWon  int    `json:"playerTokensWon"`
	PlayerTokensLost int    `json:"playerTokensLost"`
	PlayerRpcWon     int    `json:"playerRpcWon"`
}

var mutex sync.Mutex
// var updatingTableMutex sync.Mutex
var wgGroup []*sync.WaitGroup
var chGroup []chan Player
var teams []*Player
var stopGroup []chan bool
var tokensDefault int
var numero int
var addressLocal string
var addressRemoto string
var cantTeams int
var playerStart Player
var points int
var updating bool
var firstTimeReadTable bool
// var upgrader = websocket.Upgrader{}
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Permite todas las conexiones de origen (para desarrollo)
		return true
	},
}

var clients = make(map[*websocket.Conn]bool) // Mapa para rastrear los clientes conectados
type Message struct {
	Message string `json:"message"`
}

func main() {
	firstTimeReadTable=true
	// c := cors.Default()
	// handler := c.Handler(http.HandlerFunc(imprimir))
	// http.HandleFunc("/points", imprimir)
	// go http.ListenAndServe(":9000", handler)

	var wg sync.WaitGroup
	//lectura por consola del host origen
	brIn := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto del host local: ")
	puertoLocal, _ := brIn.ReadString('\n')
	puertoLocal = strings.TrimSpace(puertoLocal)
	addressLocal = fmt.Sprintf("localhost:%s", puertoLocal)

	//lectura por consola del host destino
	brIn = bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto del host remoto: ")
	puertoRemoto, _ := brIn.ReadString('\n')
	puertoRemoto = strings.TrimSpace(puertoRemoto)
	addressRemoto = fmt.Sprintf("localhost:%s", puertoRemoto)

	//lectura de nro de mensajes a recibir
	brIn = bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el numero de equipos a participar: ")
	numstr, _ := brIn.ReadString('\n')
	numstr = strings.TrimSpace(numstr)
	numero, _ = strconv.Atoi(numstr)
	cantTeams = 0

	//lectura de tokens por equipo
	brIn = bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese la cantidad de tokens por equipo: ")
	tokensstr, _ := brIn.ReadString('\n')
	tokensstr = strings.TrimSpace(tokensstr)
	tokensDefault, _ = strconv.Atoi(tokensstr)

	//creamos canal
	// chTeam = make(chan Team, 1)
	// chTeam <- player

	//habilitar el modo escucha (servidor) nodo local
	// ln, _ := net.Listen("tcp", addressLocal)
	// defer ln.Close()

	//manejo de concurrencia para aceptar conexion de clientes
	http.HandleFunc("/ws", handleWebSocket)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Servidor WebSocket iniciado en el puerto 3000")
	go http.ListenAndServe(":3000", nil)

	ln := escuchar()
	defer ln.Close()

	for cantTeams < numero-1 {
		wg.Add(1)
		conn, _ := ln.Accept()
		go manejador(conn, &wg)
	}
	wg.Wait()
	// fmt.Println("Asd")
	// fmt.Println(cantTeams)
	// fmt.Println(len(teams))
	game()
}
func imprimir(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var playersPoints []PlayerPointsTable
	for _, player := range teams {
		playersPoints = append(playersPoints, PlayerPointsTable{PlayerName: player.Name, PlayerPoints: player.Tokens, PlayerPositionX: player.positionX, PlayerPositionY: player.positionY, PlayerTokensWon: player.tokensWon, PlayerTokensLost: player.tokensLost, PlayerRpcWon: player.rpcWon})
	}
	//SERIALIZAMOS
	jsonData, err := json.Marshal(playersPoints)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Write(jsonData)
}

func escuchar() net.Listener {
	conn, _ := net.Listen("tcp", addressLocal)
	return conn
}

func manejador(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	br := bufio.NewReader(conn)
	msgJson, _ := br.ReadString('\n')

	//deserializar
	json.Unmarshal([]byte(msgJson), &playerStart)
	fmt.Println("Mensaje recibido: ")
	playerStart.positionX = 1
	playerStart.index = len(teams) + 1
	playerStart.positionY = len(teams) + 1
	playerStart.inGame = true
	playerStart.delay = 1
	playerStart.freezed = false
	playerStart.Tokens = tokensDefault
	playerStart.tokensWon = 0
	playerStart.tokensLost = 0
	playerStart.rpcWon = 0
	generatePlayer(playerStart)
	fmt.Println(playerStart)

	// teams = append(teams, &playerStart)
	cantTeams += 1
	fmt.Println(len(teams))
	// for _, t := range teams {
	// fmt.Println(t.Name)
	// fmt.Println(t.Tokens)
	// }

}

func enviar(num int) {
	conn, _ := net.Dial("tcp", addressRemoto)
	defer conn.Close()

	playerStart.Tokens = num
	//Serializar
	arrBytesMsg, _ := json.MarshalIndent(playerStart, "", " ")
	jsonStrMsg := string(arrBytesMsg)

	fmt.Println("Mensaje enviado: ")
	fmt.Println(jsonStrMsg)

	fmt.Fprintf(conn, jsonStrMsg)
}

func game() {
	var wg sync.WaitGroup

	for ind := range teams {
		// fmt.Println(player.Name, ind)
		wg.Add(1)
		go manage(teams[ind], wgGroup[ind], chGroup[ind], stopGroup[ind], &wg)
	}
	wg.Wait()
	whoWon()
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
	for !isGameFinished() {
		if player.inGame {
			wg.Add(1)
			go move(player, wg, chPlayer, stop)
			// playerInfo := <-chPlayer
			<-chPlayer
			<-stop
			// fmt.Printf("%s %d %d %d\n", playerInfo.Name, playerInfo.Tokens, playerInfo.positionX, playerInfo.positionY)
		}
		// sendUpdates()
		// fmt.Println("ASD")
	}
}

func move(player *Player, wg *sync.WaitGroup, chPlayer chan Player, stop chan bool) {
	defer wg.Done()
	wg.Add(1)
	time.Sleep(1 * time.Second)
	time.Sleep(time.Duration(player.delay) * time.Second)
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
		for !teams[player.positionY-1].inGame && player.positionY != player.index { //SE VERIFICA QUE CAIGAMOS EN UNA FILA UNICAMENTE DE JUGADORES QUE ESTEN VIVOS
			player.positionY += 1 * multiplier
			if player.positionY < 1 {
				player.positionY = len(teams)
			} else if player.positionY > len(teams) {
				player.positionY %= len(teams)
			}
		}

	} else if player.positionX == 1 && player.positionY != player.index && teams[player.positionY-1].inGame {
		//
		mutex.Lock()
		player.Tokens += 1
		teams[player.positionY-1].Tokens -= 1
		player.tokensWon += 1
		teams[player.positionY-1].tokensLost += 1
		updating=true
		fmt.Printf("%s obtuvo 1 token de %s\n", player.Name, teams[player.positionY-1].Name)
		if teams[player.positionY-1].Tokens == 0 {
			teams[player.positionY-1].inGame = false
			fmt.Printf("%s fue eliminado del juego\n", teams[player.positionY-1].Name)
		}
		player.positionX = 1
		player.positionY = player.index
		mutex.Unlock()

	} else {
		if player.positionY != player.index {
			player.positionX -= 1
		} else {
			player.positionX += 1
		}
	}
	go collisions(player, wg, chPlayer, stop)
	fmt.Printf("%s %d %d %d\n", player.Name, player.Tokens, player.positionX, player.positionY)

	player.delay = 0
	player.freezed = false
	// validateColition()
	// stopGroup[player.index-1]<-true
	chPlayer <- *player
}

func collisions(player *Player, wg *sync.WaitGroup, chPlayer chan Player, stop chan bool) {
	defer wg.Done()
	mutex.Lock()
	for ind := range teams {
		if ind != player.index-1 && !teams[ind].freezed && !player.freezed {
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
					// player.Tokens += 1
					// teams[ind].Tokens -= 1
					player.rpcWon += 1
					fmt.Printf("%s lose against %s\n", teams[ind].Name, player.Name)
					// if teams[ind].Tokens == 0{
					// teams[ind].inGame = false
					teams[ind].positionX = 1
					teams[ind].positionY = teams[ind].index
					teams[ind].delay = 5
					teams[ind].freezed = true
					// fmt.Printf("%s fue eliminado del juego\n", teams[ind].Name)
					// }else{
					// 	teams[ind].positionX=1
					// 	teams[ind].positionY=teams[ind].index
					// }

				} else {
					// player.Tokens -= 1
					// teams[ind].Tokens += 1
					fmt.Printf("%s lose against %s\n", player.Name, teams[ind].Name)
					// if player.Tokens == 0 {
					// player.inGame = false
					teams[ind].rpcWon += 1
					player.positionX = 1
					player.positionY = player.index
					player.delay = 5
					player.freezed = true
					// fmt.Printf("%s fue eliminado del juego\n", player.Name)
					// } else {
					// player.positionX = 1
					// player.positionY = player.index
					// }

				}
				// }
				break
			}
		}
	}
	mutex.Unlock()

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

func whoWon() {
	var winner string
	for _, player := range teams {
		if player.inGame == true {
			winner = player.Name
			break
		}
	}
	fmt.Printf("El ganador es %s\n", winner)
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

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar la conexión:", err)
		return
	}
	defer conn.Close()

	// Agrega el cliente al mapa de clientes conectados
	clients[conn] = true

	for !isGameFinished() {
		// Simulamos un cambio en la base de datos cada 5 segundos
		if updating || firstTimeReadTable {
			firstTimeReadTable=false
			mutex.Lock()
			data := obtenerDatosActualizados() // Obtén los datos actualizados de la base de datos
			mutex.Unlock()

			// Envía los datos actualizados al frontend
			if err := conn.WriteJSON(data); err != nil {
				log.Println("Error al enviar datos JSON:", err)
				break
			}
			updating=false
		}
	}

	// Elimina el cliente del mapa de clientes conectados al salir del ciclo
	delete(clients, conn)
}

func sendUpdates() {
	// for {
	// Simulamos cambios en la base de datos cada 10 segundos
	// time.Sleep(10 * time.Second)

	// Lógica para obtener los datos actualizados de la base de datos
	data := obtenerDatosActualizados()

	// Envía los datos actualizados a todos los clientes conectados
	broadcastData(data)
	// }
}

func obtenerDatosActualizados() interface{} {
	// Lógica para obtener los datos actualizados de la base de datos
	// Aquí puedes realizar consultas a la base de datos y retornar los resultados
	
	// return points
	var arrResponse []PlayerPointsTable
	for _,player:=range teams{
		arrResponse = append(arrResponse, PlayerPointsTable{
			PlayerName: player.Name,
			PlayerPoints: player.Tokens,
			PlayerPositionX: player.positionX,
			PlayerPositionY: player.positionY,
			PlayerTokensWon: player.tokensWon,
			PlayerTokensLost: player.tokensLost,
			PlayerRpcWon: player.rpcWon,
		})
	}
	sort.Slice(arrResponse, func(i, j int) bool {
		return arrResponse[i].PlayerPoints < arrResponse[j].PlayerPoints
	})
	return arrResponse
}

func broadcastData(data interface{}) {
	// Recorre todos los clientes conectados y envía los datos actualizados
	for client := range clients {
		if err := client.WriteJSON(data); err != nil {
			log.Println("Error al enviar datos JSON a un cliente:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
