package main
import ( 
	"net" 
	"fmt" 
	"time" 
	"encoding/json"
)


type State struct {
    CurFloor string
    HeisNr string 	 
    Astring string
}

var m State

func main() {
//	p []byte;
//	buf := Write
	state1 := State{"11", "22", "heihei"}
	//data := make(map
	
	conn, err := net.Dial("udp", "129.241.187.136:35000")
	
	if err != nil {
                fmt.Println("Could not resolve udp address or connect to it  on " )
                fmt.Println(err)
                return
        }

	b, err1 := json.Marshal(state1)
        if err1 != nil {
	        fmt.Println("error with JSON")
	        fmt.Println(err1)
        }
	err3 := json.Unmarshal(b, &m)
	if err3 != nil {
                        fmt.Println("Error")
                        fmt.Println(err3)
                        return
                }
	
	fmt.Println(m.HeisNr)


        fmt.Println("About to write to connection")

        for {
		
                time.Sleep(1000*time.Millisecond)
                n, err := conn.Write(b)
                if err != nil {
                        fmt.Println("error writing data to server")
                        fmt.Println(err)
                        return
                }

                if n > 0 {
                        fmt.Println("Wrote ",n, " bytes to server at ")
                }
        }
}
