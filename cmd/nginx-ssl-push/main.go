package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"time"

	pt "github.com/hkloudou/rpc_nginx/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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
	flag.StringVar(&grpcURL, "url", "localhost:80", "grpc URL")
	flag.StringVar(&certPath, "cert", "", "certPath")
	flag.StringVar(&keyPath, "key", "", "keyPath")
	flag.StringVar(&name, "name", "test", "Name")
	flag.StringVar(&apiKey, "apikey", "grpcnginx", "apiKey")
	flag.Parse()
	var err error
	conn, err = grpc.Dial(grpcURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Println("suc connect")
}

func main() {
	c := pt.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	items := &pt.MultSSLSetRequest{}
	items.Apikey = apiKey
	items.Item = make([]*pt.SSLSetRequest, 0)
	log.Println("read:", certPath)
	log.Println("read:", keyPath)
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
		//
		if actual, ok := status.FromError(err); ok {
			log.Println("actual", "code", actual.Code(), "err:", actual.Message())
		} else {
			log.Fatal("not actual", err)
		}
	}
	log.Println("ok")
}
