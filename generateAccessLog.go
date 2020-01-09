package main

import (
	"fmt"
    "log"
    "github.com/globalsign/mgo"
    "github.com/globalsign/mgo/bson"
    "math/rand"
    "time"
    "strconv"
)

func main() {
    //Подключаемся к mongoDB
	session, err := mgo.Dial("localhost")
	if err != nil {	log.Fatal(err) }
	defer session.Close()
	//session.SetMode(mgo.Monotonic, true)	//http://stackoverflow.com/questions/38572332/compare-consistency-models-used-in-mgo
	c := session.DB("local").C("accessLog")

    rand.Seed(time.Now().UnixNano())    //Рандомизируем генератор случайных чисел

    m := map[int][]string{}   //Объявляем карту, в которой каждому ключу (userId) соответствует массив с одним или двумя элементами (IP)
    
    for i := 0; i < 10000000; i++ {
        userId := rand.Intn(1016)+1 //Максимальное количество уникальных пользователей 1016 (1 - 1016)
        userIP := ""

        if len(m[userId]) < 2 {
            thirdOctet := strconv.Itoa(rand.Intn(10))
            fourthOctet := strconv.Itoa(rand.Intn(254)+1)
            userIP = "192.168."+thirdOctet+"."+fourthOctet
            m[userId] = append(m[userId], userIP)
        } else {
            userIP = m[userId][rand.Intn(2)]    //Один из двух IP пользователя
        }
        err := c.Insert(bson.M{"user_id": userId, "ip_addr": userIP, "ts": time.Now()})
        if err != nil { log.Fatal(err) }
    }

    fmt.Println("Готово!")
}
