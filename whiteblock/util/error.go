package util

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

/**
 * Unify error messages through function calls
 */

func CheckArguments(cmd *cobra.Command,args []string,min int,max int){
    if len(args) < min {
        fmt.Println(cmd.UsageString())
        PrintStringError("Missing arguments.")
        os.Exit(1)
    }
    if max != -1 &&  len(args) > max {
        fmt.Println(cmd.UsageString())
        PrintStringError("Too many arguments.")
        os.Exit(1)
    }
}

func InvalidArgument(arg string){
    PrintStringError(fmt.Sprintf("Invalid argument given: %s.",arg))
}

func InvalidInteger(name string,value string,fatal bool){
    PrintStringError(fmt.Sprintf("Invalid integer, given \"%s\" for %s.",value,name))
    if fatal {
        os.Exit(1)
    }
}

func ClientNotSupported(client string){
    PrintStringError(fmt.Sprintf("This function is not supported for %s.",client))
    os.Exit(1)
}

func PrintErrorFatal(err error){
    PrintError(err)
    os.Exit(1)
}

func PrintError(err error){
    PrintStringError(err.Error())
}

func PrintStringError(err string){
    fmt.Printf("\033[31mError:\033[0m %s\n",err)
}