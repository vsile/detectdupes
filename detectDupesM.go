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
        err = c.Find(bson.M{"user_id": idA}).All(&dataA)   //~50ms - We have >10M records in DB, and ~10K records for idA and idB...
    	if err != nil {
            fmt.Fprintf(w, "UserID %v %v", idA, err)
            return
        }
        err = c.Find(bson.M{"user_id": idB}).All(&dataB)
    	if err != nil {
            fmt.Fprintf(w, "UserID %v %v", idA, err)
            return
        }
        
        //Create a map with unique ip-addresses for the first UserId              
        m := map[string]bool{}
        for _, v := range dataA { m[v.Ip_addr] = true }

        //Count matches
        matches := 0
        for _, v := range dataB {
            if _, ok := m[v.Ip_addr]; ok {
                matches++
            }
        }

        w.Header().Set("Content-Type", "application/json")
        if matches > 1 {
            fmt.Fprint(w, "{\"dupes\":true}")
        } else {
            fmt.Fprint(w, "{\"dupes\":false}")
        }
    })

    fmt.Println("Service is running")
    err = http.ListenAndServe(":12345", nil)
    if err != nil { log.Fatal(err) }
}
