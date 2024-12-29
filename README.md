# よふかし／yofukashi

nocturnal software for the small web

## 日暮／higure

nex server. only active at night.

```
% higure # serves up /var/gemini
% higure -r ~/nex # serves up nex from your homedir
% higure -a # keeps the server open around the clock
% higure -lat 35 -lon 135 # use lat/lon to calculate sunrise/sunset
```

### stations served by higure

悲劇駅:
[nex://higeki.jp/](nex://higeki.jp/)

真夏駅:
[nex://manatsu.town/](nex://manatsu.town/)
## installation

### go

```
go install blekksprut.net/yofukashi/cmd/higure@latest
```

