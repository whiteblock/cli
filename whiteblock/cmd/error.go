package cmd

import (
    "fmt"
)

/**
 * Unify error messages through function calls
 */


func InvalidArgument(arg string){
    fmt.Println("\nError: invalid argument given: "+arg)
}