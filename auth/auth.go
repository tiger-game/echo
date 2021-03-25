package main

import (
	"fmt"
	"math/rand"
	"time"
)

const set = `0123456789qwertyuiopasdfghjklzxcvbnm:/\+-!#@%QWERTYUIOPASDFGHJKLZXCVBNM<>.`

func main() {
	data := make([]byte, 16)
	RandomPrivateKey(data)
	fmt.Println(string(data))
}

func RandomPrivateKey(data []byte) {
	rand.Seed(time.Now().UnixNano())
	for i := range data {
		data[i] = set[rand.Intn(len(set))]
	}
}
