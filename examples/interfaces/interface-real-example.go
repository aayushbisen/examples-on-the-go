package main

import "fmt"

// Why interfaces are needed in programming?
// tight coupling breaks as you scale
//
// Problem
// Your code depends directly on a specific implementation, so any change spreads everywhere.
//
// example of tight coupling
//

type EmailService struct{}

func (e EmailService) Send(msg string) {
	fmt.Println("Sending EMAIL:", msg)
}

// embedding email sturct in Notification struct
type Notification1 struct {
	email EmailService
}

func (n Notification1) Notify1(msg string) {
	// email message to notify
	n.email.Send(msg)
}

// What Goes Wrong When Requirements Grow?
// New requirement:
// “Also support SMS and Push notifications”
// Now we modify existing code to accomodate sms and push
type SMSService struct{}

func (s SMSService) Send(msg string) {
	fmt.Println("Sending SMS:", msg)
}

type Notification struct {
	email EmailService
	sms   SMSService
	push  PushService //push service has to be added here
}

func (n Notification) Notify(msg string, channel string) {
	if channel == "email" {
		n.email.Send(msg)
	} else if channel == "push" {
		n.push.Send(msg)
	} else if channel == "sms" {
		n.sms.Send(msg)
	}
}

// same is the story with push
type PushService struct{}

func (s PushService) Send(msg string) {
	fmt.Println("Sending Push:", msg)
}

// this is bad as we have more use cases
// Same Example Using Interface (Loose Coupling)
//
//
// define behavior

type Notifier interface {
	Send(msg string)
}

// If the type has the same method as interface it implictly implements it
type EmailService struct{}

func (e EmailService) Send(msg string) {
	fmt.Println("Sending EMAIL:", msg)
}

type SMSService struct{}

func (s SMSService) Send(msg string) {
	fmt.Println("Sending SMS:", msg)
}

// main logic

type Notification struct {
	service Notifier
}

func (n Notification) Notify(msg string) {
	// any service is satisfied here as long as it satisfies the interface
	n.service.Send(msg)
}


// what changed above?
// Now you can do:
//

n1 := Notification{service: EmailService{}}
n1.Notify("Hello")

n2 := Notification{service: SMSService{}}
n2.Notify("Hello")

// Scaling Now Becomes Easy
//
// Add new feature:
//

type PushService struct{}

func (p PushService) Send(msg string) {
    fmt.Println("Sending PUSH:", msg)
}

n3 := Notification{service: PushService{}}
n3.Notify("Hello")


// Tight coupling is bad at scale because:

// Changes spread everywhere
// Code becomes rigid and fragile
// Adding features requires rewriting existing logic

// Interfaces fix this by:

//  Separating what is done from how it is done
