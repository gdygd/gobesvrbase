# 'gobesvrbase'
## 1.tech stack
* `Language` : go (go1.21.9)
* `Database`: mariadb
* `web framework`: gorilla
* `communication`: kafka, tcp/ip, restful api, sse
* `kernel`: shared memory, signal
* `version control`: github



## 2.Description
**`gobesvrbase` is a backend server process**
1. It provides http rest api and http sse
2. It provides kafka and tcp/ip commucation
3. It provides database access(mariadb)
4. It provides shared memory between process, linux signal handling and CLI process
5. **It is just backend server base process**

## 3.Process structure
![process](https://github.com/user-attachments/assets/26749e12-674f-4c95-a1c8-b85fc67f2eef)

## 4.Process Desc
* `apimp` : API server Manage Process, mornitor 'apisvr', exe process
* `apisvr` : backend prcoess (It perform http, db access and network communication(tcp and kafka))
* `cli` : Client line interface (It provides debug command and information)

## 5.Package Desc
### 5.1) *apimp*
* It exist only main package

### 5.1) *apisvr*
* **apisvr package tree**
```
├── main
├── app
│   ├── appmodel `[global model package]`
│   ├── dbapp `[Database interface package]`
│   │   └── mdb `[mariadb access package]`
│   ├── httpapp `[http api, sse package]`
│   ├── msgapp `[network message processing package]`
│   ├── netapp `[network communication, httpapp-netapp and netapp-msgapp channel management package]`
│   ├── kafkatapp `[kafka message consume and produce]`
│   └── objdb `[Object management package]`
└── comm `[tcp/ip wrapper package]`
```
### 5.2) *cli*
* **cli package tree**
```
├── main
└── cmd `[command processing package]`
```

## 6.Usage
##### 1) clone repository
* https://github.com/gdygd/gobesvrbase.git
##### 2) compile source
* $sh make.sh
##### 3) Execute process
* move bin directory : $cd ./bin
* execute process : apimp
##### 4) Test Tcp Server process
* ref : https://github.com/gdygd/tcpserver

## 7.test api and sse
* `/gettest` : get test api, It resonse "{"result":0,"data":[{"Dt":"2024-08-05 17:09:33","Val":1}],"reqdata":null}"
* `/posttest`: post test api, It response {"result":0,"data":"PostTest","reqdata":null}" 
* `/deltest` : delete test api, It response {"result":0,"data":"DelTest","reqdata":null}
* `/netcmd`  : request network command
* `/events`  : send server time every 1sec (YYYY-MM-DD HH:MM:SS)

## 8.Data flow
### 8.1) *httpapp - dbapp*
![http-dbapp](https://github.com/user-attachments/assets/7bed880c-79f6-4341-8cac-8281b0a60794)

### 8.2) *httpapp - netapp - msgapp*
![dafaflow2](https://github.com/user-attachments/assets/d0481ece-77ea-4975-b67c-b3be053dfe48)

### 8.2) *httpapp - kafkaapp - netapp - msgapp*
![kafkaapp](https://github.com/user-attachments/assets/19fb84c0-7b36-47de-8f26-f4679e85b54e)

## 9.CLI Usage
* command list
```
 * help
 * system
 * version
 * process
 * debug
 * exit
 * termiante
```
* help : show command help message
* system : show system infomation
* version : show application version infomation
* process : show process state infomation
* debug : change logging level
* exit : quti cli
* terminate : terminate all process

### 9.1) CLI UI
![cli](https://github.com/user-attachments/assets/db07270b-56af-427d-89df-610e3b12662f)