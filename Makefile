
.PHONY: all

BIN=ratelimit
all: $(BIN)

$(BIN): *.go
	go build -o $@ *.go

clean:
	@rm -fv $(BIN)
