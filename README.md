# syslogsrv4dev

Simple log server works with syslog, Additional Udp listener can be added. 
The server doesn't require installation. The data is stored in SQLLite database

## Run application

1) run the appropriate executable file (based on OS)
    ```./syslogsrv4dev```
    Notes: by default the wep app is available at 5050, the syslog service is available at 7570 and udp listener at 9090

## Build

1) change directory:
```cd ./app```
2) build react application
```npm run build```
3) build go application
```go build```
4)  add web application into executable file
```rice append --exec ${name}x32```

## Udp listener
The regext can be defined in conf.json.
Example:
```"(?s)^(?P<date>.*?)\\|(?P<levelId>.*?)\\|(?P<facilityid>.*?)\\|(?P<message>.*?)\\|(?P<host>.*?)\\|(?P<app>.*?)$"```
The following groups in regexp will be used for log entry:
- date as Date
- levelId as severity
- facilityid as Facility
- message as Message
- host as host
- application as AppName
