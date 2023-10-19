package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"google.golang.org/genproto/googleapis/cloud/bigquery/migration/v2"
)

// func perform collection which need migration
func performMigration(client *mongo.Client) error {
	Collection := client.Database("//localhost:51.158.37.242/32").Collection("lesson")
	return Collection

}





func main()  {
	//downloading parametrs connect for Mongo

ClientOptions := options.Client().ApplyURI("mongodb+srv://zhdanovnikitk4:<password>@cluster0.gfgs4fr.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(*&options.Client().ServerAPIOptions)

// Create connection for Mongo
client, err := mongo.Connect(context.Background(), ClientOptions)

if err != nil {
	log.Fatal(err)
}

// check connection successfully

err = client.Ping(context.Background(), nil)
if err != nil {
	log.Fatal(err)
}

fmt.Println("Successfully connection in DataBase")

// Migration

err = performMigration(client)
 if err != nil {
  log.Fatal(err)
}

//Closed connection

err = client.Disconnect(context.Background())
 if err != nil {
  log.Fatal(err)
 }
 fmt.Println("Closed connection in Mongo")



	
// create context with timeout

context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
log.Fatal(cancel)



update := bson.M{
	"$set": bson.M{
	 "newField": "value",
	},

   
   _, err == collection.UpdateMany(ctx, bson.M{}, update) 
 if err != nil {
  return err
 }

 return nil
}

}
