# よふかし／yofukashi

nocturnal software for the small web

## 日暮／higure

nex server. only active at night.

```
% higure # serves up /var/gemini
% higure -r ~/nex # serves up nex from your homedir
% higure -a # keeps the server open around the clock
% higure -lat 35 # latitude for calculating dawn/dusk
```

## 星屑／hoshikuzu

nex client.

```
% hoshikuzu nex://manatsu.town/
% hoshikuzu nex://manatsu.town/sky.jpg
```

## installation

### go

```
% go install blekksprut.net/yofukashi/cmd/higure@latest
```

