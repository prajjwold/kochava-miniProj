# Kochava Miniproject :: "Postback Delivery"

###Contents
- [Description](#Description)
- [Instructions](#Instructions)
- [Extra Merit](#Extra-Merit)
- [Data Flow](#Data-Flow)
- App Operation
    - [Ingestion Agent](#App-Operation-Ingestion-Agent-(php))
    - [Delivery Agent](#App-Operation-Delivery-Agent-(go))
- Samples
    - [Request](#Sample-Request)
    - [Response](#Sample-Response-(Postback))

##Description 
You will be building a service to function as a small scale simulation of how we at Kochava distribute data to third parties in real time.

##Instructions
1) Provision provided Ubuntu server (see Resources - Server) with software stack required to complete project.
2) Change default redis port immediately on startup and add authentication
3) Build a php application to ingest http requests, and a golang application to deliver http responses. Use Redis to host a job queue between them.
4) Reach out to (Resources - Contact) once your project is ready to demo, or if you encounter a block during development. Pursue independent troubleshooting prior to escalating questions with contact resource.
5) Maintain development notes, provide support documentation, and commit your project (application code / stack config) to Github.

##Extra Merit
- Clean, descriptive Git commit history.
- Clean, easy-to-follow support documentation for an engineer attempting to troubleshoot your system.
- All services should be configured to run automatically, and service should remain functional after system restarts.
- High availability infrastructure considerations.
- Data integrity considerations, including safe shutdown.
- Modular code design.
- Configurable default value for unmatched url {key}s.
- Performance of system under external load.
- Minimal bandwidth utilization between ingestion and delivery servers.
- Configurable response delivery retry attempts. 
- Data validation / error handling.
- Ability to deliver POST (as well as GET) responses.
- Service monitoring / application profiling.
- Delivery volume / success / failure visualizations.
- Internal benchmarking tool.

##Data flow
1) Web request (see sample request) >
2) "Ingestion Agent" (php) >
3) "Delivery Queue" (redis)
4) "Delivery Agent" (go) >
5) Web response (see sample response)

##App Operation - Ingestion Agent (php)
1) Accept incoming http request
2) Push a "postback" object to Redis for each "data" object contained in accepted request.

##App Operation - Delivery Agent (go)
1) Continuously pull "postback" objects from Redis
2) Deliver each postback object to http endpoint:
    - Endpoint method: request.endpoint.method.
    - Endpoint url: request.endpoint.url, with {xxx} replaced with values from each request.endpoint.data.xxx element.
3) Log delivery time, response code, response time, and response body.

##Sample Request
```
    (POST) http://{server_ip}/ingest.php
    (RAW POST DATA) 
    {  
      "endpoint":{  
        "method":"GET",
        "url":"http://sample_domain_endpoint.com/data?title={mascot}&image={location}&foo={bar}"
      },
      "data":[  
        {  
          "mascot":"Gopher",
          "location":"https://blog.golang.org/gopher/gopher.png"
        }
      ]
    }
```

##Sample Response (Postback)
```
    GET http://sample_domain_endpoint.com/data?title=Gopher&image=https%3A%2F%2Fblog.golang.org%2Fgopher%2Fgopher.png&foo=
```