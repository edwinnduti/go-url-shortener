# go-url-shortener
A URL-SHORTENER API that uses Golang and MongoDB and works similarly like bit.ly .

![License: MIT](https://img.shields.io/badge/Language-Golang-blue.svg)
[![Build Status](https://travis-ci.com/edwinnduti/go-url-shortener.svg?branch=master)](https://travis-ci.com/edwinnduti/go-url-shortener)
![License: MIT](https://img.shields.io/badge/Database-MongoDB-lightgreen.svg)

A REST API url-shortener made in Golang.
### Requirements


<ul>
<li>GOLANG</li>
<li>POSTMAN</li>
<li>MONGO COMPASS GUI</li>
</ul>


You can run it locally and online:

##### Locally

```bash
$ git clone https://github.com/edwinnduti/go-url-shortener.git
$ cd go-url-shortener
$ go install
$ export MONGOURI=mongodb://localhost:27017
$ sudo service mongod start
$ go run main.go
```

Available locally:

| function              |   path                    |   method  |
|   ----                |   ----                    |   ----    |
| Create shorturl       |   /			    		|	POST    |
| Get single url	    |   /{id}		            |	GET     |
| Get redirected        |   /{urlid}             	|	GET     |
| Delete single url 	|   /{id}           		|	DELETE  |
| update single user	|   /{id}           		|	UPDATE  |
| Get longurl    	    |   /expand                 |   GET     |   


##### Online
```
$ curl -X POST -H "Content-Type:application/json" -d {"longurl": "<enter longurl e.g https://google.com/search?=Skygardener>"} https://localhost:8045/
```

Have Fun!
