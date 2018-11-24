package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"

	. "github.com/pspaces/gospace"
)

// Insere um novo usuário em um ambiente
func createUser(master *Space, username string, password string, long string, lat string, room string) {
	var _password string
	var _long string
	var _lat string
	var _room string
	var _msg string

	// Verifica se o usuário já existe
	_, err1 := master.QueryP(username, &_password, &_long, &_lat, &_room, &_msg)
	_, err2 := master.QueryP(room)

	if err1 == nil {
		fmt.Printf("O usuário %s já existe!\n", username)

	} else if err2 == nil { // A sala já existe, adiciona o usuário nessa sala
		master.Put(username, password, long, lat, room, "")

	} else if err2 != nil { // A sala não existe, então cria a sala e adiciona o usuário nela
		master.Put(room)
		master.Put(username, password, long, lat, room, "")
	}
}

// Lista os usuários
func listUser(master *Space) {
	var username string
	var password string
	var long string
	var lat string
	var room string
	var msg string

	t, _ := master.QueryAll(&username, &password, &long, &lat, &room, &msg)

	fmt.Println("\nUsuários:")
	for i := 0; i < len(t); i++ {
		fmt.Printf("%s: long: %s lat: %s, sala: %s\n", t[i].GetFieldAt(0), t[i].GetFieldAt(2),
			t[i].GetFieldAt(3), t[i].GetFieldAt(4))
	}
}

// Ligar Radar
func turnOnRadar(master *Space, selfname string) {

	var selfPassword string
	var selfLong string
	var selfLat string
	var selfRoom string
	var selfMsg string

	var username string
	var password string
	var long string
	var lat string
	var room string
	var msg string

	t1, _ := master.QueryP(selfname, &selfPassword, &selfLong, &selfLat, &selfRoom, &selfMsg)
	myLong, _ := strconv.ParseFloat(t1.GetFieldAt(2).(string), 64)
	myLat, _ := strconv.ParseFloat(t1.GetFieldAt(3).(string), 64)

	t2, _ := master.QueryAll(&username, &password, &long, &lat, &room, &msg)

	for i := 0; i < len(t2); i++ {
		long, _ := strconv.ParseFloat(t2[i].GetFieldAt(2).(string), 64)
		lat, _ := strconv.ParseFloat(t2[i].GetFieldAt(3).(string), 64)

		// Lista distância dos outros usuários que estão na mesma sala
		if selfname != t2[i].GetFieldAt(0).(string) && t1.GetFieldAt(4) == t2[i].GetFieldAt(4) {
			var distance = math.Sqrt(math.Pow((myLong-long), 2) + math.Pow((myLat-lat), 2))

			if distance > 300 && distance < 20000 {
				fmt.Printf("Usuário %s está a %s metros do usuário %s\n", t2[i].GetFieldAt(0), FloatToString(distance), selfname)
			} else if distance < 300 {
				fmt.Printf("Usuário %s está a menos de 300 metros do usuário %s\n", t2[i].GetFieldAt(0), selfname)
			} else if distance > 20000 {
				fmt.Printf("Usuário %s está a mais de 20000 metros do usuário %s\n", t2[i].GetFieldAt(0), selfname)
			}
		}
	}
}

func createRoom(master *Space, roomName string) {

	_, err1 := master.QueryP(roomName)
	if err1 == nil {
		fmt.Printf("A sala %s já existe!\n", roomName)

	} else { // A sala não existe, então cria a sala e adiciona o usuário nela
		master.Put(roomName)
	}
}

func listRoom(master *Space) {
	var _room string

	t, _ := master.QueryAll(&_room)

	fmt.Printf("\nSalas:\n")
	for i := 0; i < len(t); i++ {
		fmt.Println(t[i].GetFieldAt(0))
	}
}

func deleteRoom(master *Space, room string) {
	var _username string
	var _password string
	var _long string
	var _lat string
	var _msg string

	// Exclui a sala da tupla de salas
	master.GetP(room)

	// Adquire todos os usuários que estão na sala indicada por parâmetro
	t, _ := master.QueryAll(&_username, &_password, &_long, &_lat, room, &_msg)

	for i := 0; i < len(t); i++ {
		// Remove os usuários da sala indicada, excluindo eles da tupla principal
		master.GetP(t[i].GetFieldAt(0).(string), t[i].GetFieldAt(1).(string),
			t[i].GetFieldAt(2).(string), t[i].GetFieldAt(3).(string), room, t[i].GetFieldAt(5).(string))

		// Insere os usuários novamente, mas em nenhuma sala
		master.Put(t[i].GetFieldAt(0).(string), t[i].GetFieldAt(1).(string),
			t[i].GetFieldAt(2).(string), t[i].GetFieldAt(3).(string), "nenhuma", "")
	}
}

func enterRoom(master *Space, username string, room string) {
	var _password string
	var _long string
	var _lat string
	var _room string
	var _msg string

	// Procura o usuário indicado e armazena os dados em t
	t, _ := master.QueryP(username, &_password, &_long, &_lat, &_room, &_msg)

	// Se o usuário já não estiver na sala indicada (por ex, entrar na sala1 já estando nela)
	if room != t.GetFieldAt(4) {
		// Remove o usuário da tupla principal
		master.GetP(username, &_password, &_long, &_lat, room, &_msg)

		// Insere o usuário na tupla principal com a nova sala
		master.Put(username, t.GetFieldAt(1).(string), t.GetFieldAt(2).(string), t.GetFieldAt(3).(string), room, "")

		fmt.Printf("Usuário %s saiu da sala %s e entrou na sala %s\n", username, t.GetFieldAt(4).(string), room)
	}
}

func exitRoom(master *Space, username string) {
	var _password string
	var _long string
	var _lat string
	var _room string
	var _msg string

	// Procura o usuário indicado, deleta e armazena os dados em t
	t, _ := master.GetP(username, &_password, &_long, &_lat, &_room, &_msg)

	// Recria o usuário, mas em nenhuma sala
	master.Put(username, t.GetFieldAt(1).(string), t.GetFieldAt(2).(string),
		t.GetFieldAt(3).(string), "nenhuma", "")
}

func sendMessage(master *Space, username string, willSend string) {
	var _username string
	var _password string
	var _long string
	var _lat string
	var _room string
	var _msg string

	// Procura o usuário indicado e armazena os dados em t
	t, _ := master.QueryP(username, &_password, &_long, &_lat, &_room, &_msg)

	// Armazena todos os usuários que estão NA MESMA SALA do usuário que vai enviar a mensagem
	t1, _ := master.GetAll(&_username, &_password, &_long, &_lat, t.GetFieldAt(4), &_msg)

	// A operação anterior remove inclusive o usuário que está enviando a mensagem
	// portanto é necessário readicioná-lo
	master.Put(username, t.GetFieldAt(1), t.GetFieldAt(2), t.GetFieldAt(3), t.GetFieldAt(4), t.GetFieldAt(5))

	for i := 0; i < len(t1); i++ {
		if username != t1[i].GetFieldAt(0).(string) {
			// Insere os usuários novamente na tupla principal e com a mensagem que foi enviada
			master.Put(t1[i].GetFieldAt(0).(string), t1[i].GetFieldAt(1).(string),
				t1[i].GetFieldAt(2).(string), t1[i].GetFieldAt(3).(string),
				t1[i].GetFieldAt(4).(string), willSend)
		}
	}

	// Adquire todos os usuários para extrair a mensagem que foi recebida
	t2, _ := master.QueryAll(&_username, &_password, &_long, &_lat, t.GetFieldAt(4).(string), &_msg)

	for i := 0; i < len(t2); i++ {
		fmt.Printf("\nUsuário %s recebeu: %s\n", t2[i].GetFieldAt(0), t2[i].GetFieldAt(5))
	}
}

func changeCoordinates(master *Space, username string, long string, lat string) {
	var _password string
	var _long string
	var _lat string
	var _room string
	var _msg string

	// Exclui o usuário indicado por parâmetro e armazena seus dados em t
	t, _ := master.GetP(username, &_password, &_long, &_lat, &_room, &_msg)

	master.Put(username, t.GetFieldAt(1), long, lat, t.GetFieldAt(4), t.GetFieldAt(5))

	t1, _ := master.QueryP(username, &_password, &_long, &_lat, &_room, &_msg)

	fmt.Printf("\nNovas coordenadas do usuário %s: long: %s lat: %s\n\n", username,
		t1.GetFieldAt(2).(string), t1.GetFieldAt(3))
}

func main() {

	//Cria uma Tuple Space
	master := NewSpace("master")

	var option string

	var username string
	var password string
	var long string
	var lat string
	var room string
	var msg string

	for {

		fmt.Printf("\n\n1 - Criar usuário\n")
		fmt.Printf("2 - Listar usuários\n")
		fmt.Printf("3 - Ligar Radar\n")
		fmt.Printf("4 - Criar Sala\n")
		fmt.Printf("5 - Excluir Sala\n")
		fmt.Printf("6 - Listar Salas\n")
		fmt.Printf("7 - Entrar em Sala\n")
		fmt.Printf("8 - Sair da Sala\n")
		fmt.Printf("9 - Enviar mensagem\n")
		fmt.Printf("10 - Mudar coordenadas\n")

		fmt.Printf("Opção: ")
		fmt.Scanf("%s", &option)

		if option == "1" {
			clearScreen()
			fmt.Printf("Nome do usuário: ")
			fmt.Scanf("%s", &username)

			fmt.Printf("Senha do usuário: ")
			fmt.Scanf("%s", &password)

			fmt.Printf("Longitude do usuário em metros: ")
			fmt.Scanf("%s", &long)

			fmt.Printf("Latitude usuário em metros: ")
			fmt.Scanf("%s", &lat)

			listRoom(&master)
			fmt.Printf("Em qual sala deseja entrar?: ")
			fmt.Scanf("%s", &room)

			createUser(&master, username, password, long, lat, room)

		} else if option == "2" {
			clearScreen()
			listUser(&master)

		} else if option == "3" {
			clearScreen()
			listUser(&master)
			fmt.Printf("\nLigar radar de qual usuário?: ")
			fmt.Scanf("%s", &username)

			turnOnRadar(&master, username)

		} else if option == "4" {
			clearScreen()
			fmt.Printf("Qual o nome da sala? ")
			fmt.Scanf("%s", &room)

			createRoom(&master, room)

		} else if option == "5" {
			clearScreen()
			listRoom(&master)
			fmt.Printf("Excluir qual sala? ")
			fmt.Scanf("%s", &room)

			deleteRoom(&master, room)

		} else if option == "6" {
			clearScreen()
			listRoom(&master)

		} else if option == "7" {
			clearScreen()
			listUser(&master)
			fmt.Printf("\nColocar qual usuário em uma sala? ")
			fmt.Scanf("%s", &username)

			listRoom(&master)
			fmt.Printf("Colocar o usuário %s em qual sala? ", username)
			fmt.Scanf("%s", &room)

		} else if option == "8" {
			clearScreen()
			listUser(&master)
			fmt.Printf("\nTirar qual usuário da sala? ")
			fmt.Scanf("%s", &username)

			exitRoom(&master, username)

		} else if option == "9" {
			clearScreen()
			listUser(&master)
			fmt.Printf("\nQual usuário irá enviar uma mensagem? ")
			fmt.Scanf("%s", &username)

			fmt.Printf("\nInsira a mensagem: ")
			fmt.Scanf("%s", &msg)

			sendMessage(&master, username, msg)

		} else if option == "10" {
			clearScreen()
			listUser(&master)
			fmt.Printf("\nMudar coordenadas de qual usuário?: ")
			fmt.Scanf("%s", &username)

			fmt.Printf("Longitude do usuário em metros: ")
			fmt.Scanf("%s", &long)

			fmt.Printf("Latitude usuário em metros: ")
			fmt.Scanf("%s", &lat)

			changeCoordinates(&master, username, long, lat)
		}
	}
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// FloatToString converts a float to string
func FloatToString(inputNum float64) string {
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}
