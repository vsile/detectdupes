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

type cacheAccessLog struct {
    User_id int `_id`
    Ip_addr []string
}

func main() {
    //Connecting to mongoDB
	session, err := mgo.Dial("localhost")
	if err != nil {
        log.Fatal(err)
    }
	defer session.Close()
	c := session.DB("local").C("cacheAccessLogOnFly")

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {     //localhost:12345/?a=2&b=4
        start := time.Now()
        q := r.URL.Query()
        idA, errA := strconv.Atoi(q.Get("a"))
        idB, errB := strconv.Atoi(q.Get("b"))
        if errA != nil || errB != nil {
            fmt.Fprint(w, "UserID format is not allowed")
            return
        }

        dataA, dataB := cacheAccessLog{}, cacheAccessLog{}
        err = c.Find(bson.M{"_id": idA}).One(&dataA)
    	if err != nil {
            fmt.Fprintf(w, "UserID %v %v", idA, err)
            return
        }
        err = c.Find(bson.M{"_id": idB}).One(&dataB)
    	if err != nil {
            fmt.Fprintf(w, "UserID %v %v", idB, err)
            return
        }
        
        //Create a map with unique ip-addresses for the first UserId              
        m := map[string]bool{}
        for _, v := range dataA.Ip_addr {
            m[v] = true
        }

        //Count matches
        matches := 0
        for _, v := range dataB.Ip_addr {
            if _, ok := m[v]; ok {
                matches++
            }
        }
        
        w.Header().Set("Content-Type", "application/json")
        if matches > 1 {
            fmt.Fprint(w, "{\"dupes\":true}")
        } else {
            fmt.Fprint(w, "{\"dupes\":false}")
        }
        fmt.Println(time.Now().Sub(start))
    })

    fmt.Println("Service is running")
    err = http.ListenAndServe(":12345", nil)
    if err != nil {
        log.Fatal(err)
    }
}
