GODEPS = \
	"github.com/jessevdk/go-flags"

BIN = minit

all: $(BIN)

.deps:
	GOPATH=`pwd`/.deps go get -u $(GODEPS)

$(BIN): $(BIN).go .deps
	GOPATH=`pwd`/.deps go build $<

clean:
	rm -rf .deps $(BIN)
