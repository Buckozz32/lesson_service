package database

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	_ "github.com/eminetto/clean-architecture-go/migrations"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"

	migrate "github.com/eminetto/mongo-migrate"
"github.com/globalsign/mgo"
)


func Cli()  {
	if len(os.Args) == 1 {
		log.Fatal("Missing options: up or down")
	}
option := os.Args[1]

session, err := mgo.Dial(os.Getenv("MONGOBD_HOST"))
if err != nil {
	log.Fatal(err.Error())
}

defer session.Close()

db := session.DB(os.Getenv("MONGODB_DATABASE"))

migrate.SetDatabase(db)

migrate.SetMigrationsCollection("migrations")
migrate.SetLogger(log.New(os.Stdout, "INFO: ", 0))

switch option {
case "new": 
if len(os.Args) != 3 {
	log.Fatal("Should be: new description-of-migration")
}
lesson := fmt.Sprintf("./migrations/%s_%s.go", time.Now().Format(""), os.Args[2])
from, err := os.Open("./migration.go")
if err != nil {
	log.Fatal("Should be: new description-of-migration")
}
defer from.Close()
}
token, err := os.OpenFile(lesson, os.O_RDWR|os.O_CREATE, "0666")
if err != nil {
	log.Fatal(err.Error())
}
defer token.Close()

_, err = io.Copy(token, &os.File{})
if err != nil {
	log.Fatal(err.Error())
}

log.Printf("new migration created: %s\n", lesson)


_, err = migrate.Up(migrate.AllAvailable)

_, err = migrate.Down(migrate.AllAvailable)



if err != nil {
	log.Fatal(err.Error())
}

}


