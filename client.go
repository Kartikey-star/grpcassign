package main

import (
	"context"
	"log"
	"time"

	"github.com/kartikey-star/grpcassign/user"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8089", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := user.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	entry := user.User{
		Firstname:    "Kartikey",
		Lastname:     "Mamgain",
		Emailaddress: "mamgainkartikey@gmail.com",
	}
	// select seat for user
	item := &user.UserRequest{User: &entry, From: "London", Section: "A", To: "France", Seat: "50"}
	res, err := c.Create(ctx, item)
	if err != nil {
		log.Fatalf("could not create item: %v", err)
	}
	log.Printf("User Receipt generated: %v", res.GetUser())

	// Read an item
	resRead, err := c.Get(ctx, &user.ReadUserRequest{User: &entry})
	if err != nil {
		log.Fatalf("could not read item: %v", err)
	}
	log.Printf("User Receipt read: %v", resRead.GetUser())

	// Update an item
	update := &user.UpdateUserRequest{User: &entry, From: "London", Section: "A", To: "France", Seat: "53"}
	resUpdate, err := c.Update(ctx, update)
	if err != nil {
		log.Fatalf("could not update item: %v", err)
	}
	log.Printf("Seat updated for user : %v%v", resUpdate.GetUser().Emailaddress, resUpdate.GetSeat())

	// Delete user booking
	resDelete, err := c.Delete(ctx, &user.DeleteUserRequest{User: &entry})
	if err != nil {
		log.Fatalf("could not delete item: %v", err)
	}
	log.Printf("Item deleted: %v", resDelete.GetUser())

	// make multiple entries in memory
	item = &user.UserRequest{User: &entry, From: "London", Section: "B", To: "France", Seat: "50"}
	c.Create(ctx, item)

	entry.Emailaddress = "dummy@gmail.com"
	item = &user.UserRequest{User: &entry, From: "London", Section: "A", To: "France", Seat: "54"}
	c.Create(ctx, item)

	resList, err := c.List(ctx, &user.ListUserRequest{Section: "A"})
	if err != nil {
		log.Fatalf("could not delete item: %v", err)
	}
	log.Printf("Item list: %v", resList.UserReceiptlist)

}
