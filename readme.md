# 'gobesvrbase'
## 1.Description
**`gobesvrbase` is a backend server process**
1. It provides http rest api, http sse
2. It provides tcp/ip commucation, database access(mariadb)
3. It provides shared memory between process, linux signal handling, CLI process
4. **It is just backend server base process**


## 2.Process Desc
* `apimp` : API server Manage Process
* `apisvr` : backend prcoess (It perform http, db access and network communication)
* `cli` : Client line interface (It provides debug command and information)

## 3.Package Desc
### 3.1) *apimp*
* It exist only main package

### 3.1) *apisvr*
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
│   └── objdb `[Object management package]`
└── comm `[tcp/ip wrapper package]`
```
### 3.2) *cli*
* **cli package tree**
```
├── main
└── cmd `[command processing package]`
```


## 4.Usage
-clone and make

## 5.test api and sse
* `/gettest` : get test api, It resonse "{"result":0,"data":[{"Dt":"2024-08-05 17:09:33","Val":1}],"reqdata":null}"
* `/posttest`: post test api, It response {"result":0,"data":"PostTest","reqdata":null}" 
* `/deltest` : delete test api, It response {"result":0,"data":"DelTest","reqdata":null}
* `/netcmd`  : request network command
* `/events`  : send server time every 1sec (YYYY-MM-DD HH:MM:SS)

## 6.Data flow
### 6.1) *httpapp - dbapp*
![http-dbapp](https://github.com/user-attachments/assets/7bed880c-79f6-4341-8cac-8281b0a60794)

### 6.2) *httpapp - netapp - msgapp*
![http-netapp](https://github.com/user-attachments/assets/ccd3a2fd-0341-44ff-a642-70fa66aca3ab)
