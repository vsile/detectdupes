package main

import (
	"fmt"
    "log"
    "github.com/globalsign/mgo"
    "github.com/globalsign/mgo/bson"
    "net/http"
    "time"
    "strconv"
)

type accessLog struct {
    User_id int
    Ip_addr string
    Ts      time.Time 
}

func main() {
    //Connecting to mongoDB
	session, err := mgo.Dial("localhost")
	if err != nil {	log.Fatal(err) }
	defer session.Close()
	c := session.DB("local").C("accessLog")

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {     //localhost:12345/?a=2&b=4
        q := r.URL.Query()
        idA, errA := strconv.Atoi(q.Get("a"))
        idB, errB := strconv.Atoi(q.Get("b"))
        if errA != nil || errB != nil { fmt.Fprint(w, "UserID format is not allowed"); return }

        dataA, dataB := []accessLog{}, []accessLog{}
        errA = c.Find(bson.M{"user_id": idA}).All(&dataA)   //3.5s - We have >10M records in DB, and ~10K records for idA and idB...
        errB = c.Find(bson.M{"user_id": idB}).All(&dataB)   //3.5s
        if errA != nil || errB != nil { fmt.Fprint(w, err.Error()); return }
        
        //Create a map with unique ip-addresses for the first UserId              
        m := map[string]bool{}
        for _, v := range dataA { m[v.Ip_addr] = true }
        startMapLength := len(m)

        //Delete keys from map
        for _, v := range dataB { delete(m, v.Ip_addr) }
        endMapLength := len(m)

        //fmt.Println(m)

        w.Header().Set("Content-Type", "application/json")
        if startMapLength - endMapLength > 1 {  //Return true if map length decreased more than 1
            fmt.Fprint(w, "{\"dupes\":true}")
        } else {
            fmt.Fprint(w, "{\"dupes\":false}")
        }
    })

    fmt.Println("Service is running")
    err = http.ListenAndServe(":12345", nil)
    if err != nil { log.Fatal(err) }
}
