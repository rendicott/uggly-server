package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/rendicott/uggly"
	"github.com/rendicott/uggly-server/siteconfig"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 10000, "The server port")
	siteConfig = flag.String("sites", "site.yml", "yaml file containing site definitions")
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

func newServer(sc *siteconfig.Sites) *feedServer {
	fServer := &feedServer{}
	fServer.divBoxes = &pb.DivBoxes{}
	log.Printf("convering config to uggly pb's for %d sites\n", len(sc.Sites))
	for _, site := range sc.Sites {
		for _, sbox := range site.DivBoxes {
			ubox := pb.DivBox{
				Name:       sbox.Name,
				Border:     sbox.Border,
				BorderW:    sbox.BorderW,
				BorderChar: sbox.BorderChar,
				FillChar:   sbox.FillChar,
				StartX:     sbox.StartX,
				StartY:     sbox.StartY,
				Width:      sbox.Width,
				Height:     sbox.Height,
			}
			if sbox.BorderSt != nil {
				ubox.BorderSt = &pb.Style{
					Fg:   sbox.BorderSt.Fg,
					Bg:   sbox.BorderSt.Bg,
					Attr: sbox.BorderSt.Attr,
				}
			}
			if sbox.FillSt != nil {
				ubox.FillSt = &pb.Style{
					Fg:   sbox.FillSt.Fg,
					Bg:   sbox.FillSt.Bg,
					Attr: sbox.FillSt.Attr,
				}
			}
			fServer.divBoxes.Boxes = append(fServer.divBoxes.Boxes, &ubox)
		}
		log.Printf("have divboxes of len %d\n", len(fServer.divBoxes.Boxes))
		fServer.elements = &pb.Elements{}
		for _, sele := range site.Elements {
			for _, sblob := range sele.TextBlobs {
				ublob := pb.TextBlob{
					Content:  sblob.Content,
					Wrap:     sblob.Wrap,
					DivNames: sblob.DivNames,
				}
				if sblob.Style != nil {
					ublob.Style = &pb.Style{
						Fg: sblob.Style.Fg,
						Bg: sblob.Style.Bg,
					}
				}
				fServer.elements.TextBlobs = append(
					fServer.elements.TextBlobs, &ublob)
			}
		}
	}
	log.Printf("have textblobs of len %d\n", len(fServer.elements.TextBlobs))
	return fServer
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// parse site config
	sites, err := siteconfig.NewSiteConfig(*siteConfig)
	if err != nil {
		log.Printf("error parsing site config file: '%s'\n", err.Error())
		os.Exit(1)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	s := newServer(sites)
	pb.RegisterFeedServer(grpcServer, *s)
	grpcServer.Serve(lis)
	log.Println("Server listening")
}
