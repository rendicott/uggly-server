package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/rendicott/uggly"
	"github.com/rendicott/uggly-server/pageconfig"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"errors"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 10000, "The server port")
	pageConfigFile = flag.String("pages", "pages.yml", "yaml file containing page definitions")
)

/* pageServerPage holds elements required by the protobuf definition for a page which
includes its elemntal properties
*/
type pageServerPage struct {
	name string
	divBoxes *pb.DivBoxes
	elements *pb.Elements
	response *pb.PageResponse
}
/* pageServer is a struct from which to attach the required methods for the Page Service
as defined in the protobuf definition
*/
type pageServer struct {
	pb.UnimplementedPageServer
	pages []*pageServerPage
}

/* feedServer is a struct from which to attach the required methods for the Feed Service
as defined in the protobuf definition
*/
type feedServer struct {
	pb.UnimplementedFeedServer
	pages []*pb.PageListing
}

/* GetFeed implements the Feed Service's GetFeed method as required in the protobuf definition.

It is the primary listening method for the server. It accepts a FeedRequest and then attempts to build
a FeedResponse which the client will process. 
*/
func (f feedServer) GetFeed(ctx context.Context, freq *pb.FeedRequest) (fresp *pb.FeedResponse, err error) {
	fresp = &pb.FeedResponse{}
	fresp.Pages = f.pages
	return fresp, err
}

/* GetPage implements the Page Service's GetPage method as required in the protobuf definition.

It is the primary listening method for the server. It accepts a PageRequest and then attempts to build
a PageResponse which the client will process and display on the client's pcreen. 
*/
func (s pageServer) GetPage(ctx context.Context, preq *pb.PageRequest) (presp *pb.PageResponse, err error) {
	found := false
	for _, page := range(s.pages) {
		if page.name == preq.Name {
			found = true
			page.response.Name = page.name
			return page.response, err
		}
	}
	if !found {
		err = errors.New("requested page not found")
	}
	return presp, err
}

/* newFeedServer takes the loaded pageconfig YAML and converts it to the structs
required so that the GetFeed method can adequately respond with a FeedResponse which
is primarily a list of the pages this server serves.
*/
func newFeedServer(pc *pageconfig.Pages) *feedServer {
	fServer := &feedServer{}
	for _, page := range pc.Pages {
		sListing := &pb.PageListing{}
		sListing.Name = page.Name
		fmt.Printf("Have page name: %s\n", page.Name)
		// ./server.go:82:17: first argument to append must be slice; have *uggly.Pages
		fServer.pages = append(fServer.pages, sListing)
	}
	return fServer
}


/* newPageServer takes the loaded pageconfig YAML and converts it to the structs
required so that the GetPage method can adequately respond with a PageResponse.
*/
func newPageServer(pc *pageconfig.Pages) *pageServer {
	pServer := &pageServer{}
	for i := range(pc.Pages) {
		psp := pageServerPage{}
		psp.name = pc.Pages[i].Name
		psp.divBoxes = &pb.DivBoxes{}
		for _, pbox := range pc.Pages[i].DivBoxes {
			ubox := pb.DivBox{
				Name:       pbox.Name,
				Border:     pbox.Border,
				BorderW:    pbox.BorderW,
				BorderChar: pbox.BorderChar,
				FillChar:   pbox.FillChar,
				StartX:     pbox.StartX,
				StartY:     pbox.StartY,
				Width:      pbox.Width,
				Height:     pbox.Height,
			}
			if pbox.BorderSt != nil {
				ubox.BorderSt = &pb.Style{
					Fg:   pbox.BorderSt.Fg,
					Bg:   pbox.BorderSt.Bg,
					Attr: pbox.BorderSt.Attr,
				}
			}
			if pbox.FillSt != nil {
				ubox.FillSt = &pb.Style{
					Fg:   pbox.FillSt.Fg,
					Bg:   pbox.FillSt.Bg,
					Attr: pbox.FillSt.Attr,
				}
			}
			psp.divBoxes.Boxes = append(psp.divBoxes.Boxes, &ubox)
		}
		log.Printf("have divboxes of len %d\n", len(psp.divBoxes.Boxes))
		psp.elements = &pb.Elements{}
		for _, sele := range pc.Pages[i].Elements {
			for _, pblob := range sele.TextBlobs {
				ublob := pb.TextBlob{
					Content:  pblob.Content,
					Wrap:     pblob.Wrap,
					DivNames: pblob.DivNames,
				}
				if pblob.Style != nil {
					ublob.Style = &pb.Style{
						Fg: pblob.Style.Fg,
						Bg: pblob.Style.Bg,
					}
				}
				psp.elements.TextBlobs = append(
					psp.elements.TextBlobs, &ublob)
			}
		}
		log.Printf("have textblobs of len %d\n", len(psp.elements.TextBlobs))
		// now pre-build the response
		psp.response = &pb.PageResponse{}
		psp.response.DivBoxes = &pb.DivBoxes{}
		psp.response.Elements = &pb.Elements{}
		psp.response.DivBoxes = psp.divBoxes
		psp.response.Elements = psp.elements
		pServer.pages = append(pServer.pages, &psp)
	}
	return pServer
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// parse page config
	pcPages, err := pageconfig.NewPageConfig(*pageConfigFile)
	if err != nil {
		log.Printf("error parsing page config file: '%s'\n", err.Error())
		os.Exit(1)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	f := newFeedServer(pcPages)
	pb.RegisterFeedServer(grpcServer, *f)
	s := newPageServer(pcPages)
	pb.RegisterPageServer(grpcServer, *s)
	grpcServer.Serve(lis)
	log.Println("Server listening")
}
