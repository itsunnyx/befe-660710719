package main

import (
	"fmt"
)

// var email String = "sittipornchaiya_n@su.ac.th"

func main(){
	// var name string = "Nanthawat"
	var age int = 21

	email := "sittipornchaiya_n@su.ac.th"
	gpa := 4.00

	firstname, lastname := "Nanthawat", "Sittipornchaiyakhun"

	fmt.Printf("Name %s %s, age %d, email %s, gpa %.2f\n", firstname, lastname, age, email, gpa)
}