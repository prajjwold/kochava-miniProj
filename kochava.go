package main

import (
	redis "gopkg.in/redis.v5"
	"fmt"
	"regexp"
	"net/url"
	"net/http"
	"io"
)

var rqstList string = "requests"
var procList string = "processing"

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

func beginProcessing(conn *redis.Client, hashId string) {
	defer endProcessing(conn, hashId)

	fmt.Println("Processing...")

	reqData, err := getRequestData(hashId, conn)
	if err != nil {
		fmt.Errorf("%s\n", err)
	}

	fmtdReq := formatRequest(reqData)
	fmt.Println(fmtdReq)
	resp, err := http.Get(fmtdReq)
	if(err != nil) {
		panic(err)
	}

	logResponse(resp)
}

func getRequest(conn *redis.Client) (string, error) {
	cmd := conn.BRPopLPush(rqstList, procList, 0)
	if err := cmd.Err(); err != nil {
		panic(err)
	}

	return cmd.Result()
}

func getRequestData(hashId string, conn *redis.Client) (map[string]string, error) {
	return conn.HGetAll(hashId).Result()
}

func formatRequest(data map[string]string) string {
	result := data["endpoint"]

	for key, val := range data {
		pattern := regexp.MustCompile("{" + key + "}")
		result = pattern.ReplaceAllString(result, url.QueryEscape(val))
	}

	return result
}

func logResponse(resp *http.Response) {
	body := make([]byte, resp.ContentLength)

	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	bytesRead, err := resp.Body.Read(body)
	if err != io.EOF {
		panic(err)
	}

	fmt.Println(string(body[:bytesRead]))
}

func endProcessing(conn *redis.Client, hashId string) {
	cmd := conn.LRem(procList, 0, hashId)
	rmRes, err := cmd.Result()
	if err != nil {
		fmt.Errorf("%s\n", err)
	} else if !(rmRes >= 1) {
		fmt.Errorf("%s\n", "No values found in the \"processing\" list to remove.")
	}

	cmd = conn.Del(hashId)
}

func main() {
	redisClient, err := redisConnect("127.0.0.1:7331", "staaben-miniProj", 0)

	if(err != nil) {
		fmt.Println(err)
	} else {
		fmt.Println("Connected to Redis server.")
		for {
			poppedReq, err := getRequest(redisClient)
			if err != nil {
				fmt.Errorf("%s\n", err)
			}

			go beginProcessing(redisClient, poppedReq)


		}
	}
}