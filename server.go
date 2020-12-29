package main

import (
    "context"
    "net"
    "log"
    "google.golang.org/grpc"
    "flag"
    "fmt"
    pb "github.com/rendicott/uggly"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 10000, "The server port")
)

type feedServer struct {
    pb.UnimplementedFeedServer
    divBoxes *pb.DivBoxes
}

func (s feedServer) GetDivs(ctx context.Context, ss *pb.FeedRequest) (boxes *pb.DivBoxes, err error) {
    // serve divs for feedrequest
    return s.divBoxes, err
}


func newServer() *feedServer {
	s := &feedServer{}
    s.divBoxes = &pb.DivBoxes{}
    box1 := pb.DivBox{
        Border: true,
        BorderW: 1,
        BorderChar: []rune("+")[0],
        FillChar: []rune(" ")[0],
        StartX: 8,
        StartY: 8,
        Width: 40,
        Height: 8,
    }
    box2 := pb.DivBox{
        Border: true,
        BorderW: 1,
        BorderChar: []rune("-")[0],
        FillChar: []rune("*")[0],
        StartX: 8,
        StartY: 30,
        Width: 10,
        Height: 6,
    }
    s.divBoxes.Boxes = append(s.divBoxes.Boxes, &box1)
    s.divBoxes.Boxes = append(s.divBoxes.Boxes, &box2)
	return s
}


func main() {
    flag.Parse()
    lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    var opts []grpc.ServerOption
    grpcServer := grpc.NewServer(opts...)
    s := newServer()
    pb.RegisterFeedServer(grpcServer, *s)
    grpcServer.Serve(lis)
    log.Println("Server listening")
}
