package main

import (
	redis "gopkg.in/redis.v5"
	"regexp"
	"net/url"
	"net/http"
	"io"
	"os"
	"log"
	"strconv"
	"time"
	"strings"
)

// Global level variables
var rqstList string = "requests"
var procList string = "processing"
var logFilename string = "kochava.log"

// Postback Object
// Struct to contain data relevant for logging information about the "postback" object
// as well as working with "postback" objects
type PostbackLog struct {
	deliveryTime int
	responseCode int
	responseBody string
	responseTime int64
	postbackId string
	requestMethod string
}

func main() {
	redisClient, err := redisConnect("127.0.0.1:7331", "staaben-miniProj", 0)
	logFile := setupLogger(redisClient)
	defer logFile.Close()

	if(err != nil) {
		log.Panicf("%v\n", err)
		panicExit(redisClient)
	} else {
		log.Println("Connected to Redis server.")
		for {
			poppedReq, err := getRequest(redisClient)
			if err != nil {
				log.Panicf("%v\n", err)
				panicExit(redisClient)
			}

			go beginProcessing(redisClient, poppedReq)
		}
	}
}

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

// Open the file to log all information to
func setupLogger(conn *redis.Client) *os.File {
	file, err := os.OpenFile(logFilename, os.O_CREATE | os.O_RDWR | os.O_APPEND, 0666)
	if err != nil {
		log.Panic("Couldn't open log file.")
		panicExit(conn)
	}

	log.SetOutput(file)

	return file
}

// Go routine function to begin processing the passed in "postback" object from the Redis server
// Retrieves request data from server, sets it in the Postback struct for the request
// Logs delivery time and response time from request
//
func beginProcessing(conn *redis.Client, hashId string) {
	defer endProcessing(conn, hashId)
	pb := new(PostbackLog)
	pb.postbackId = hashId

	log.Printf("Processing %v\n", hashId)

	reqData, err := getRequestData(hashId, conn)
	if err != nil {
		log.Printf("%v\nUnable to process %v.\n", err, hashId)
		return
	}
	pb.requestMethod = reqData["method"]

	fmtdReq := formatRequest(reqData, conn)
	recTime, _ := strconv.Atoi(reqData["receivedTime"])
	pb.deliveryTime = (int(conn.Time().Val().Unix()) - recTime)

	response, respTime, err := getResponse(fmtdReq, pb)
	if(err != nil) {
		log.Panicf("%v\n", err)
		panicExit(conn)
	}
	pb.responseTime = respTime

	pb.SetLogData(response, conn)

	pb.Log()
}

// Get the request ID from the server
func getRequest(conn *redis.Client) (string, error) {
	cmd := conn.BRPopLPush(rqstList, procList, 0)
	if err := cmd.Err(); err != nil {
		log.Panicf("%v\n", err)
		panicExit(conn)
	}

	return cmd.Result()
}

// Get the request data from the given request ID
func getRequestData(hashId string, conn *redis.Client) (map[string]string, error) {
	return conn.HGetAll(hashId).Result()
}

// Format the URL contained within the request
func formatRequest(data map[string]string, conn *redis.Client) string {
	result := data["endpoint"]
	method := data["method"]

	switch method {
		case "GET":
			for key, val := range data {
				cleanKey := strings.Replace(key, "data:", "", -1)
				pattern := regexp.MustCompile("{" + cleanKey + "}")
				result = pattern.ReplaceAllString(result, url.QueryEscape(val))
			}
			break
		default:
			log.Panic("Invalid request method.")
			panicExit(conn)
	}

	return result
}

// Get the HTTP response from the endpoint provided
func getResponse(url string, pb *PostbackLog) (*http.Response, int64, error) {
	var response *http.Response
	var err error
	var responseTime int64
	var respEnd int64

	respStart := int64(time.Now().Unix())

	switch pb.requestMethod {
		case "GET":
			response, err = http.Get(url)
			respEnd = int64(time.Now().Unix())
			break
	}
	responseTime = (respEnd - respStart)

	return response, responseTime, err
}

// Log all information in the Postback struct in the logging file (default of kochava.log)
func (p *PostbackLog) Log() {
	log.Printf("%v:\n\tDelivery time: %v\n\tResponse time: %v\n\tResponse code: %v\nBEGIN RESPONSE BODY:\n%v\nEND RESPONSE BODY",
		p.postbackId, p.deliveryTime, p.responseTime, p.responseCode, p.responseBody)
}

// Sets all pertinent information in the Postback object using the provided HTTP response
func (p *PostbackLog) SetLogData(resp *http.Response, conn *redis.Client) {
	body := make([]byte, resp.ContentLength)
	defer resp.Body.Close()

	bytesRead, err := resp.Body.Read(body)
	if err != io.EOF {
		log.Panicf("%v\n", err)
		panicExit(conn)
	}

	p.responseBody = string(body[:bytesRead])

	p.responseCode = resp.StatusCode
}

// Log exiting message and shutdown the Redis server
func panicExit(conn *redis.Client) {
	log.Println("Exiting...")
	shutdownCmd := conn.ShutdownSave()
	log.Println(shutdownCmd.Result())
	os.Exit(-1)
}

// Clean-up after successful processing of a request
func endProcessing(conn *redis.Client, hashId string) {
	cmd := conn.LRem(procList, 0, hashId)
	rmRes, err := cmd.Result()
	if err != nil {
		log.Printf("%v\n", err)
	} else if !(rmRes >= 1) {
		log.Printf("%s\n", "No values found in the \"processing\" list to remove.")
	}

	cmd = conn.Del(hashId)
}