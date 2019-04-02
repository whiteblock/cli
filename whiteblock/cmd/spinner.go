package cmd

import (
    "fmt"
    "time"
)


type Spinner struct {
    txt     string
    die     bool
}

func (this *Spinner) Run(ms int){
    go func(){
        states := []string{"/","-","\\","|","/","-","\\","|"}
        
        i := 0
        this.die = false
        for {
            if this.die {
                this.die = false
                fmt.Println("\033[K\r")
                return
            }
            time.Sleep(time.Duration(ms)*time.Millisecond)
            fmt.Printf("\033[K%s %s\r",this.txt,states[i])
            i++
            i = i % len(states) 
            if this.die {
                this.die = false
                fmt.Printf("\033[K\r")
                return
            }
        }
    }()
    
}

func (this *Spinner) SetText(txt string){
    this.txt = txt
}

func (this *Spinner) Kill() {
    this.die = true
}
