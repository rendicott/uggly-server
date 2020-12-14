package main

import (
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

type screenerServer struct {
    pb.UnimplementedScreenerServer
    screenSets map[string]*pb.Screens
}

func (s screenerServer) GetScreens(ss *pb.ScreenSet, stream pb.Screener_GetScreensServer) (err error) {
    // serve screens for requested screenset
    for _, i := range(s.screenSets[ss.Name].Screens) {
        if err := stream.Send(i); err != nil {
            return err
        }
    }
    return err
}


func newServer() *screenerServer {
	s := &screenerServer{screenSets: make(map[string]*pb.Screens)}
    screens := pb.Screens{}
    screen1 := pb.Screen{
        Contents: "hello grpc world",
    }
    screen2 := pb.Screen{
        Contents: "happy to be here",
    }
    screens.Screens = append(screens.Screens, &screen1)
    screens.Screens = append(screens.Screens, &screen2)
    s.screenSets["one"] = &screens
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
    pb.RegisterScreenerServer(grpcServer, *s)
    grpcServer.Serve(lis)
}
