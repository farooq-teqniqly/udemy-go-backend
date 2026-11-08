package main

import (
	"fmt"
)

type Contact struct {
	Id    int
	Name  string
	Email string
	Phone string
}

var contacts []Contact
var contactIndexByName map[string]int
var nextId = 1

func init() {
	contacts = make([]Contact, 0)
	contactIndexByName = make(map[string]int)
}

func addContact(name, email, phone string) {
	if _, exists := contactIndexByName[name]; exists {
		fmt.Printf("Contact already exists: %s\n", name)
		return
	}

	newContact := Contact{
		Id:    nextId,
		Name:  name,
		Email: email,
		Phone: phone,
	}

	contactIndexByName[name] = len(contacts)
	contacts = append(contacts, newContact)
	nextId++

	fmt.Printf("Contact added: %+v\n", newContact)
}

func findContactByName(name string) *Contact {
	index, exists := contactIndexByName[name]

	if !exists {
		return nil
	}

	return &contacts[index]
}

func listContacts() {
	if len(contacts) == 0 {
		fmt.Println("No contacts")
	}

	for _, contact := range contacts {
		fmt.Printf("%+v\n", contact)
	}
}

func main() {
	listContacts()

	addContact("Farooq", "farooq@foo.com", "+18335553456")
	addContact("Bob", "bob@foo.com", "+183355537654")
	addContact("Jane", "jane@foo.com", "+183355537659")
	addContact("Alice", "alice@foo.com", "+183355538000")

	bob := findContactByName("Bob")
	fmt.Printf("%s\n", bob.Name)

	ted := findContactByName("Ted")
	fmt.Printf("%+v\n", ted)

	alice := findContactByName("Alice")
	fmt.Printf("%+v\n", alice)

	listContacts()
}
