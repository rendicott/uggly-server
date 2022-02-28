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
	"errors"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 10000, "The server port")
	siteConfigFile = flag.String("sites", "site.yml", "yaml file containing site definitions")
)

/* siteServerSite holds elements required by the protobuf definition for a site which
includes its elemntal properties
*/
type siteServerSite struct {
	name string
	divBoxes *pb.DivBoxes
	elements *pb.Elements
	response *pb.SiteResponse
}
/* siteServer is a struct from which to attach the required methods for the Site Service
as defined in the protobuf definition
*/
type siteServer struct {
	pb.UnimplementedSiteServer
	sites []*siteServerSite
}

/* feedServer is a struct from which to attach the required methods for the Feed Service
as defined in the protobuf definition
*/
type feedServer struct {
	pb.UnimplementedFeedServer
	sites []*pb.SiteListing
}

/* GetFeed implements the Feed Service's GetFeed method as required in the protobuf definition.

It is the primary listening method for the server. It accepts a FeedRequest and then attempts to build
a FeedResponse which the client will process. 
*/
func (f feedServer) GetFeed(ctx context.Context, freq *pb.FeedRequest) (fresp *pb.FeedResponse, err error) {
	fresp = &pb.FeedResponse{}
	fresp.Sites = f.sites
	return fresp, err
}

/* GetSite implements the Site Service's GetSite method as required in the protobuf definition.

It is the primary listening method for the server. It accepts a SiteRequest and then attempts to build
a SiteResponse which the client will process and display on the client's screen. 
*/
func (s siteServer) GetSite(ctx context.Context, sreq *pb.SiteRequest) (sresp *pb.SiteResponse, err error) {
	found := false
	for _, site := range(s.sites) {
		if site.name == sreq.Name {
			found = true
			return site.response, err
		}
	}
	if !found {
		err = errors.New("requested site not found")
	}
	return sresp, err
}

/* newFeedServer takes the loaded siteconfig YAML and converts it to the structs
required so that the GetFeed method can adequately respond with a FeedResponse which
is primarily a list of the sites this server serves.
*/
func newFeedServer(sc *siteconfig.Sites) *feedServer {
	fServer := &feedServer{}
	for _, site := range sc.Sites {
		sListing := &pb.SiteListing{}
		sListing.Name = site.Name
		fmt.Printf("Have site name: %s\n", site.Name)
		// ./server.go:82:17: first argument to append must be slice; have *uggly.Sites
		fServer.sites = append(fServer.sites, sListing)
	}
	return fServer
}


/* newSiteServer takes the loaded siteconfig YAML and converts it to the structs
required so that the GetSite method can adequately respond with a SiteResponse.
*/
func newSiteServer(sc *siteconfig.Sites) *siteServer {
	sServer := &siteServer{}
	for i := range(sc.Sites) {
		ssite := siteServerSite{}
		ssite.name = sc.Sites[i].Name
		ssite.divBoxes = &pb.DivBoxes{}
		for _, sbox := range sc.Sites[i].DivBoxes {
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
			ssite.divBoxes.Boxes = append(ssite.divBoxes.Boxes, &ubox)
		}
		log.Printf("have divboxes of len %d\n", len(ssite.divBoxes.Boxes))
		ssite.elements = &pb.Elements{}
		for _, sele := range sc.Sites[i].Elements {
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
				ssite.elements.TextBlobs = append(
					ssite.elements.TextBlobs, &ublob)
			}
		}
		log.Printf("have textblobs of len %d\n", len(ssite.elements.TextBlobs))
		// now pre-build the response
		ssite.response = &pb.SiteResponse{}
		ssite.response.DivBoxes = &pb.DivBoxes{}
		ssite.response.Elements = &pb.Elements{}
		ssite.response.DivBoxes = ssite.divBoxes
		ssite.response.Elements = ssite.elements
		sServer.sites = append(sServer.sites, &ssite)
	}
	return sServer
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// parse site config
	scSites, err := siteconfig.NewSiteConfig(*siteConfigFile)
	if err != nil {
		log.Printf("error parsing site config file: '%s'\n", err.Error())
		os.Exit(1)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	f := newFeedServer(scSites)
	pb.RegisterFeedServer(grpcServer, *f)
	s := newSiteServer(scSites)
	pb.RegisterSiteServer(grpcServer, *s)
	grpcServer.Serve(lis)
	log.Println("Server listening")
}
