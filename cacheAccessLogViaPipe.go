package main

import (
	"fmt"
    "log"
    "github.com/globalsign/mgo"
    "github.com/globalsign/mgo/bson"
    "time"
)

func main() {
    //Connecting to mongoDB
	session, err := mgo.Dial("localhost")
	if err != nil { log.Fatal(err) }
	defer session.Close()
    session.SetMode(mgo.Eventual, true)
	c := session.DB("local").C("accessLog")

    start := time.Now()

    p := c.Pipe([]bson.M{
        { "$group": bson.M{"_id": "$user_id", "ip_addr": bson.M{ "$addToSet": "$ip_addr" } }},
        { "$out": "cacheAccessLogOnFly" },
    }).Iter()
    err = p.Close()
	if err != nil { log.Fatal(err) }

    fmt.Println(time.Now().Sub(start))
    fmt.Println("Cache created successfully")
}
