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
		private $dataId;
		
		public function __construct($new_method, $new_url, $new_data) {
			$this->setMethod($new_method);
			$this->setUrl($new_url);
			$this->setData($new_data);
			$this->setHashId("");
			$this->setDataId("");
		}
		
		public function getMethod() { return $this->method; }
		public function getUrl() { return $this->url; }
		public function getData() { return $this->dataAra; }
		public function getHashId() { return $this->hashId; }
		public function getDataId() { return $this->dataId; }
		
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
		
		public function setDataId($new_dataId) {
			$this->dataId = $new_dataId;
		}
		
		private function createRedisIds($redisConn) {
			preg_match('/http:\/\/(.*)\//', $this->getUrl(), $matches);
			
			if(sizeof($matches) == 0) {
				http_response_code(500);
			}
			
			$id = $matches[1].rand();
			$this->setHashId($id.":postback");
			$this->setDataId($id.":data");
			while($redisConn->exists($this->getHashId()) && $redisConn->exists($this->getDataId())) {
				$id = $matches[1].rand();
				$this->setHashId($id.":postback");
				$this->setDataId($id.":data");
			}
		}
		
		public function storeData($redisConn) {
			$this->createRedisIds($redisConn);
			
			foreach($this->getData() as $key => $val) {
				$redisConn->hSet($this->getDataId(), $key, $val);
			}
		}
		
		public function getHashFields() {
			return array("endpoint" => $this->getUrl(), "method" => $this->getMethod(), "dataId" => $this->getDataId());
		}
	}
	
?>