# 'gobesvrbase'
데이터 베이스, HTTP, 소켓, Kafka를 이용한 go언어 기반 백엔드 서버 기본 모델

## 1.기술 스택
* `Language` : go (go1.21.9)
* `Database`: mariadb
* `web framework`: gorilla
* `communication`: kafka, tcp/ip, restful api, sse

## 2.백엔드 서버 설명
**`gobesvrbase`**
1. http RestApi and SSE 기능을 제공
2. Kafka and Tcp/ip 기능을 제공
3. Database access(mariadb)
4. 프로세스간 통신은 shared memory를 이용함
5. CLI 기능을 제공(디버깅, 상태테크 등)


## 3.프로세스 구성
![process](https://github.com/user-attachments/assets/26749e12-674f-4c95-a1c8-b85fc67f2eef)

## 4.프로세스 설명
* `apimp` : 메인 프로세스, 하위프로세스의 동작을 관리
* `apisvr` : 백엔드 프로세스, http, sse, kafka, tcp/ip, database access의 기능을 수행
* `cli` : Client line interface 

## 5.패키지 설명
### 5.1) *apimp*
* 메인 패키지

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
│   ├── netapp `[tcp/ip 통신 패키지]`
│   ├── kafkatapp `[kafka message consume and produce 패키지]`
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
* `/gettest` : get test api, 응답 : "{"result":0,"data":[{"Dt":"2024-08-05 17:09:33","Val":1}],"reqdata":null}"
* `/posttest`: post test api, 응답 : "{"result":0,"data":"PostTest","reqdata":null}" 
* `/deltest` : delete test api, 응답 "{"result":0,"data":"DelTest","reqdata":null}"
* `/netcmd`  : request network command(소켓을 통한 메세지 전송)
* `/events`  : SSE, send server time every 1sec (YYYY-MM-DD HH:MM:SS)

## 8.Data flow
### 8.1) *httpapp - dbapp* 
#### - http를 통해 데이터베이스 정보를 조회하는 구조
![http-dbapp](https://github.com/user-attachments/assets/7bed880c-79f6-4341-8cac-8281b0a60794)

### 8.2) *httpapp - netapp - msgapp* (사용자가 http를 통해 외부 장치에 request)
#### - http를 통해 외부 장치에 request
#### - 외부 장치 정보를 SSE를 통해 사용자에게 제공
![dafaflow2](https://github.com/user-attachments/assets/d0481ece-77ea-4975-b67c-b3be053dfe48)

### 8.2) *httpapp - kafkaapp - netapp - msgapp*
#### - 외부 시스템을 통해 kafka 메세지 생산
#### - 수집한 kafka 메세지를 SSE를 통해 사용자에게 제공
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
* help : show help message
* system : show system infomation
* version : show application version infomation
* process : show process state infomation
* debug : change logging level
* exit : quti cli
* terminate : terminate all process
