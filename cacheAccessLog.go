package main

import (
	"fmt"
    "log"
    "github.com/globalsign/mgo"
    "github.com/globalsign/mgo/bson"
    "time"
    //"sync"
)

type accessLog struct {
    User_id int
    Ip_addr string
}

func main() {
    //Connecting to mongoDB
	session, err := mgo.Dial("localhost")
	if err != nil {	log.Fatal(err) }
	defer session.Close()
    session.SetMode(mgo.Eventual, true)
	c := session.DB("local").C("accessLog")
    c_ := session.DB("local").C("cacheAccessLog")

    start := time.Now()

    l := 1016   //Number of Users
    //Loop over the all UserIds
    //var wg sync.WaitGroup
    //wg.Add(l)
    for i := 1; i < l+1; i++ {
        //go func(i int) {
            data := []accessLog{}
            err = c.Find(bson.M{"user_id": i}).All(&data)
            if err != nil {	log.Fatal(err) }

            m := map[string]bool{}
            for _, v := range data {
                m[v.Ip_addr] = true
            }
            ips := []string{}
            for k := range m {
                ips = append(ips, k)
            }
            _, err = c_.Upsert(bson.M{"user_id": i}, bson.M{"$set": bson.M{"ip_addr": ips}})
            if err != nil {	log.Fatal(err) }
            //wg.Done()
        //}(i)
    }
    //wg.Wait()

    fmt.Println(time.Now().Sub(start))

    fmt.Println("Cache created successfully")
}
