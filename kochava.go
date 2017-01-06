package main

import (
	redis "gopkg.in/redis.v5"
	"fmt"
)

// Create and return a *redis.Client connected to the default
func redisConnect(addr string, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password,
		DB: db,
	})

	_, err := client.Ping().Result()

	return client, err
}

func main() {
	_, err := redisConnect("127.0.0.1:7331", "staaben-miniProj", 0)
	if(err != nil) {
		fmt.Println(err)
	}


}