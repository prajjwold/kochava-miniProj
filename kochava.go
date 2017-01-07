package main

import (
	redis "gopkg.in/redis.v5"
	"fmt"
)

var rqstList string = "requests";

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
	redisClient, err := redisConnect("127.0.0.1:7331", "staaben-miniProj", 0)

	if(err != nil) {
		fmt.Println(err)
	} else {
		//for {
			popResp := redisClient.BLPop(0, rqstList)
			if err = popResp.Err(); err != nil {
				panic(err)
			}
			postback, err := popResp.Result()
			if err != nil {
				panic(err)
			}

			scanResp := redisClient.HScan(postback[0], 0, "", 0)
			if err = scanResp.Err(); err != nil {
				panic(err)
			}
			keys, cursor, err := scanResp.Result()
			if err != nil {
				panic(err)
			}

			fmt.Println(keys, cursor)
		//}

		pushResp := redisClient.RPush(rqstList, keys)
		if err = pushResp.Err(); err != nil {
			panic(err)
		}
	}
}