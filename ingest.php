<?php

	require_once("postback.php");
	
	// If no data received, return error code
	if(!isset($_POST["endpoint"], $_POST["data"])) { http_response_code(400); return; }

	// Static constant key for Redis list of requests
    $rqstList = "requests";
	
	// Create a new Redis client and connect, return an error code if unable to connect
	$redis = new Redis();
	if(!$redis->connect('127.0.0.1', 7331)) {
		http_response_code(500);
	}
	
	// Authorize the connection to the Redis server
	$redis->auth("staaben-miniProj");
	// Store time from Redis server for when the request was received
	$recTime = $redis->time()[0];
	
	// Get POST data
	$endObj = $_POST["endpoint"];
	$dataAra = $_POST["data"];
	
	// Create a "postback" object from the given POST data
	$postback = new Postback($endObj["method"], $endObj["url"], $dataAra, $recTime);
	// Create and store ID for Redis key
	$postback->createRedisId($redis);
	// Get the array containing the fields for a Redis Hash object
	$hashFields = $postback->getHashFields();
	$rKey = $postback->getHashId();
	
	// Create the Redis Hash object on the Redis server, return an error code if an error
	// occurs while store the Hash object
	$resp = $redis->hMSet($rKey, $hashFields);
	if(!$resp) {
		http_response_code(500);
	}

	$resp = $redis->rPush($rqstList, $rKey);
	if(!$resp) {
	    http_response_code(500);
    }

?>