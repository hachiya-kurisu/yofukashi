# よふかし／yofukashi

nocturnal software for the small web

![yofukashi in action (gif)](yofukashi.gif)

## 日暮／higure

nex server. only active at night.

```
$ higure # serves up /var/nex
$ higure -r ~/nex # serves up nex from your homedir
$ higure -a # keeps the server open around the clock
$ higure -lat 35 -lon 139 # latitude/longitude for calculating dawn/dusk
```

## 星屑／hoshikuzu

nex client.

```
$ hoshikuzu nex://manatsu.town/
$ hoshikuzu nex://manatsu.town/sky.jpg
```

## author

[蜂谷栗栖](//blekksprut.net/)
## installation

### go

```
$ go install blekksprut.net/yofukashi/cmd/higure@latest
$ go install blekksprut.net/yofukashi/cmd/hoshikuzu@latest
```
