// server.go
package ciobs

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func startServer(exit chan int) {
	arith := new(Arith)

	s := rpc.NewServer()
	s.Register(arith)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error", e)
	}

	for {
		c, e := l.Accept()
		if e != nil {
			log.Fatal("server droped client connection:", e)
		}

		go s.ServeCodec(jsonrpc.NewServerCodec(c))
	}

	exit <- 1
}

func main() {
	exit := make(chan int)
	go startServer(exit)

	conn, e := net.Dial("tcp", "localhost:1234")
	if e != nil {
		log.Fatal("client failed to connect", e)
	}
	defer conn.Close()

	args := &Args{7, 8}
	var reply int

	c := jsonrpc.NewClient(conn)

	e = c.Call("Arith.Multiply", args, &reply)
	if e != nil {
		log.Fatal("arith error:", e)
	}
	fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)

	<-exit
}
