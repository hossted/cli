/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"

	"github.com/hossted/cli/cmd"
)

func main() {
	fmt.Println(cmd.BUILDTIME)
	cmd.Execute()
}
