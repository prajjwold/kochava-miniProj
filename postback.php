<?php

	/*
		Postback PHP class to encapsulate functionality related to storing objects
		in a Redis database
		
		Author: Corbin Staaben
	*/
	class Postback {
		private $method;
		private $url;
		private $dataAra;
		private $hashId;
		private $receivedTime;
		
		public function __construct($new_method, $new_url, $new_data, $recTime) {
			$this->setMethod($new_method);
			$this->setUrl($new_url);
			$this->setData($new_data);
			$this->setHashId("");
			$this->setReceivedTime($recTime);
		}
		
		public function getMethod() { return $this->method; }
		public function getUrl() { return $this->url; }
		public function getData() { return $this->dataAra; }
		public function getHashId() { return $this->hashId; }
		public function getReceivedTime() { return $this->receivedTime; }
		
		public function setMethod($new_method) {
			$this->method = $new_method;
		}
		
		public function setUrl($new_url) {
			$this->url = $new_url;
		}
		
		public function setData($new_data) {
			$this->dataAra = $new_data;
		}
		
		public function setHashId($new_hashId) {
			$this->hashId = $new_hashId;
		}
		
		public function setReceivedTime($new_time) {
			$this->receivedTime = $new_time;
		}
		
		/*
		 * 	Create ID for Redis Hash object by using a regex to extract the base URL 
		 * 	and appending a random integer and the string ":postback"
		 *  
		 * 	Returns true if there is a match and no errors occur; false if there are 
		 * 	no matches from the regex
		*/
		public function createRedisId($redisConn) {
			preg_match('/http:\/\/(.*)\//', $this->getUrl(), $matches);
			
			if(sizeof($matches) == 0) {
				return FALSE;
			}
			
			$id = $matches[1].rand();
			$this->setHashId($id.":postback");
			while($redisConn->exists($this->getHashId())) {
				$id = $matches[1].rand();
				$this->setHashId($id.":postback");
			}
			
			return TRUE;
		}
		
		/*
		 *	Creates an associative array with the raw endpoint (including the received method), 
		 *  all data objects, and the time the request was received
		 *
		 *	returns the created array
		*/
		public function getHashFields() {
			$end = $this->getUrl();
			$m = $this->getMethod();
			$recTime = $this->getReceivedTime();
			$data = $this->getData();
			$result = array("endpoint" => $end, "method" => $m, "receivedTime" => $recTime);
			
			foreach($data as $ara) {
				foreach($ara as $key => $val) {
					$result[$key] = $val;
				}
			}
			
			return $result;
		}
	}
	
?>