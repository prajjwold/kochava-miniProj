<?php

	class Postback {
		private $method;
		private $url;
		private $dataAra;
		
		function __construct($new_method, $new_url, $new_data) {
			$this->setMethod($new_method);
			$this->setUrl($new_url);
			$this->setData($new_data);
		}
		
		public function getMethod() { return $this->$method; }
		public function getUrl() { return $this->$url; }
		public function getData() { return $this->$data; }
		
		public function setMethod($new_method) {
			$this->$method = $new_method;
		}
		
		public function setUrl($new_url) {
			$this->$url = $new_url;
		}
		
		public function setData($new_data) {
			$this->$data = $new_data;
		}
		
		public function asRedisHashCmd() {
			
		}
	}
	
?>