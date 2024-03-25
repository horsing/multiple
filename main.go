package main

import "fmt"

func modify(slice []int) {
	slice[0] = 10
}

type person struct {
	name string
	age  int
}

func day10() {
	p := person{
		name: "Alice",
		age:  25,
	}
	changePerson(p)
	fmt.Println(p.name, p.age)
}

func changePerson(p person) {
	p.name = "Bob"
	p.age = 30
}

func main() {
	a := []int{1, 2, 3}
	modify(a)
	fmt.Println(a[0])
	day10()

	for pos, char := range "简体\x80中文" { // \x80 is an illegal UTF-8 encoding
		fmt.Printf("character %#U starts at byte position %d\n", char, pos)
	}

}