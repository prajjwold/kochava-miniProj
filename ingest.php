<?php

	require_once("postback.php");
	
	//check for POST data before beginning
	
	$redis = new Redis();
	
	$redis->connect('127.0.0.1', 7331);
	$redis->auth("staaben-miniProj");
	$recTime = $redis->time();
	
	$endObj = $_POST["endpoint"];
	$dataAra = $_POST["data"];
	
	$pb = new Postback($endObj["method"], $endObj["url"], $dataAra);
	$pb->storeData($redis);
	$hashFields = $pb->getHashFields();
	$rKey = $pb->getHashId();
	
	$resp = $redis->hMSet($rKey, $hashFields);
	//check return from hMSet here
	
?>