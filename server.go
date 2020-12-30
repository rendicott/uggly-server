package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/rendicott/uggly"
	"google.golang.org/grpc"
	"log"
	"net"
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
    elements *pb.Elements
}

func (s feedServer) GetFeed(ctx context.Context, freq *pb.FeedRequest) (fresp *pb.FeedResponse, err error) {
    fresp = &pb.FeedResponse{}
    fresp.DivBoxes = &pb.DivBoxes{}
    fresp.Elements = &pb.Elements{}
    log.Printf("attaching feedserver divboxes of len %d to "+
        "resp.DivBoxes with mem address %v",
        len(s.divBoxes.Boxes), &fresp.DivBoxes)
    fresp.DivBoxes = s.divBoxes
    log.Printf("attaching feedserver elements of len %d to "+
        "resp.Elements with mem address %v",
        len(s.elements.TextBlobs), &fresp.Elements)
    fresp.Elements = s.elements
	return fresp, err
}

func newServer() *feedServer {
	s := &feedServer{}
	s.divBoxes = &pb.DivBoxes{}
	box1 := pb.DivBox{
        Name:       "big",
		Border:     true,
		BorderW:    1,
		BorderChar: []rune("+")[0],
		FillChar:   []rune(" ")[0],
		StartX:     8,
		StartY:     8,
		Width:      40,
		Height:     8,
	}
	box2 := pb.DivBox{
        Name:       "little",
		Border:     true,
		BorderW:    1,
		BorderChar: []rune("-")[0],
		FillChar:   []rune("*")[0],
		StartX:     8,
		StartY:     30,
		Width:      20,
		Height:     6,
	}
	s.divBoxes.Boxes = append(s.divBoxes.Boxes, &box1)
	s.divBoxes.Boxes = append(s.divBoxes.Boxes, &box2)
    s.elements = &pb.Elements{}
    tb := pb.TextBlob{
        Content: "hello world, how are ya now? Good n you?",
        Wrap: true,
        DivNames: []string{"little"},
    }
    s.elements.TextBlobs = append(s.elements.TextBlobs, &tb)
    log.Printf("have textblobs of len %d\n", len(s.elements.TextBlobs))
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
