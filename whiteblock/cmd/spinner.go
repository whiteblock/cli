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
        //states := getState2()
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


func getState1() []string{
    states := []string{}
    base := "ᕦ(▀̿ ̿ -▀̿ ̿ )つ/̵͇̿̿/’̿’̿ ̿ ̿̿ ̿̿ ̿̿"
    for x := 0; x < 120; x++ {
        add := base
        for y := 0; y < x; y++ {
            add += " "
        }
        add += "*"
        states = append(states,add)
    }
    return states
}

func getState2() []string {
    size := 80
    states := []string{}
    base := "'̿'\\̵͇̿̿\\з=( ͡° ͜ʖ ͡°)=ε/̵͇̿̿/’̿’̿"
    for x := 0; x < size; x++ {
        add := base
        for y := 0; y < x; y++ {
            add += " "
        }
        for y := 0; y < size; y++ {
            if y == x {
                add = "*" + add
            }else{
                add = " " + add
            }
        }
        add += "*"
        states = append(states,add)
    }
    return states
}
