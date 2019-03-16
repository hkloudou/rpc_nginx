package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"time"

	pt "github.com/hkloudou/rpc_nginx/proto"
	"google.golang.org/grpc"
)

var (
	grpcURL  string
	certPath string
	keyPath  string
	name     string
	apiKey   string
)

var conn *grpc.ClientConn

func init() {
	var err error
	conn, err = grpc.Dial("localhost", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
}

func main() {
	flag.StringVar(&grpcURL, "url", "http://localhost", "grpc URL")
	flag.StringVar(&certPath, "cert", "", "certPath")
	flag.StringVar(&keyPath, "key", "", "keyPath")
	flag.StringVar(&name, "name", "test", "Name")
	flag.StringVar(&apiKey, "apikey", "grpcnginx", "apiKey")
	c := pt.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	items := &pt.MultSSLSetRequest{}
	items.Apikey = ""
	items.Item = make([]*pt.SSLSetRequest, 0)

	if cert, err := ioutil.ReadFile(certPath); err != nil {
		log.Fatal(err)
	} else if key, err := ioutil.ReadFile(keyPath); err != nil {
		log.Fatal(err)
	} else {
		items.Item = append(items.Item, &pt.SSLSetRequest{
			Directory: "",
			CertName:  name,
			KeyName:   name,
			Cert:      cert,
			Key:       key,
		})
	}

	if _, err := c.MultSSLSet(ctx, items); err != nil {
		log.Fatal(err)
	}
}
