package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
)

type Player struct {
	conn net.Conn
	room string
	id   int
}
type Room struct {
	players []Player
	NextID  int
}

var rooms map[string]*Room

/*
w := bufio.NewWriter(conn)
				w.Write(data)
				w.Flush()
*/

//r := bufio.NewReader(conn)

/*
fullBuf := make([]byte, 1+48+Obj2Size)
		readLen, err := r.Read(fullBuf)
		if err != nil {

        }
*/

/*
ln, _ = net.Listen("tcp", ":8081")
var err error
			conn, err = ln.Accept()
            race2(conn, false)
*/

/*
ln.Close()
		conn.SetDeadline(time.Now())
*/

func main() {
	rooms = make(map[string]*Room)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	port = ":" + port
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer listen.Close()
	fmt.Printf("Listening on %s\n", port)
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal("Error Acceptiong: %s\n", err.Error())
		}
		r := bufio.NewReader(conn)

		room, _, _ := r.ReadLine()
		if _, exists := rooms[string(room)]; !exists {
			fmt.Printf("Creating Room: %s\n", string(room))
			rooms[string(room)] = &Room{make([]Player, 0), 0}
		}
		fmt.Printf("Adding Player to: %s\n", string(room))
		p := Player{conn, string(room), rooms[string(room)].NextID}
		rooms[string(room)].NextID++
		rooms[string(room)].players = append(rooms[string(room)].players, p)
		go handlePlayer(p)
	}
}

var levelByteLen = 1
var idByteLen = 4
var obj1ByteLen = 0x30
var obj2BtyeLen = 0x20

func handlePlayer(p Player) {
	fmt.Printf("Handling Player: %d\n", p.id)
	r := bufio.NewReader(p.conn)
	byteID := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteID, uint32(p.id))
	for {
		fullBuf := make([]byte, levelByteLen+obj1ByteLen+obj2BtyeLen)

		_, err := r.Read(fullBuf)
		if err != nil {
			fmt.Printf("Err: %s\n", err)
			fmt.Printf("Connection Closed\n")
			//remove player from room
			fmt.Printf("Removing Player: %d\n", p.id)
			for i := 0; i < len(rooms[p.room].players); i++ {
				if p.id == rooms[p.room].players[i].id {
					list := rooms[p.room].players
					list[len(list)-1], list[i] = list[i], list[len(list)-1]
					rooms[p.room].players = list[:len(list)-1]
					break
				}
			}
			if len(rooms[p.room].players) == 0 {
				fmt.Printf("Deleting Room %s\n", p.room)
				delete(rooms, p.room)
			}
			//if room empty delete room

			return
		}
		players := rooms[p.room].players
		for i := 0; i < len(players); i++ {
			if !(players[i].id == p.id) {
				w := bufio.NewWriter(players[i].conn)
				data := make([]byte, 0)
				data = append(data, byteID...)
				data = append(data, fullBuf...)
				w.Write(data)
				w.Flush()
			}
		}
	}
}
