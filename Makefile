TEST = ./nex,.

all: higure hoshikuzu

again: clean all

OS != uname -s

PREFIX ?= /usr

hoshikuzu: yofukashi.go nex/nex.go cmd/hoshikuzu/main.go
	go build -C cmd/hoshikuzu -o ../../hoshikuzu

higure: yofukashi.go nex/nex.go cmd/higure/main.go
	go build -C cmd/higure -o ../../higure
	
clean:
	rm -f higure hoshikuzu

test:
	YOFUKASHI_TEST_URL=nex://manatsu.town/ \
	go test -cover -coverpkg $(TEST)

cover:
	YOFUKASHI_TEST_URL=nex://manatsu.town/ \
	go test -coverprofile=cover.out -coverpkg $(TEST)
	go tool cover -html cover.out

install:
	install higure /usr/local/bin
	install hoshikuzu /usr/local/bin

push:
	got send
	git push github

fmt:
	gofmt -s -w *.go */*.go cmd/*/*.go

README.md: README.gmi INSTALL.gmi
	cat README.gmi INSTALL.gmi | sisyphus -a "." -f markdown > README.md

doc: README.md

release: push
	git push github --tags

