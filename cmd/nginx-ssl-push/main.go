package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"time"

	nginx "github.com/hkloudou/rpc_nginx/nginx"
	"gitlab.me/font/gtls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	grpcURL  string
	certPath string
	keyPath  string
	caPath   string
	name     string
	apiKey   string
)

var conn *grpc.ClientConn

func init() {
	flag.StringVar(&grpcURL, "url", "", "grpc URL")
	flag.StringVar(&certPath, "cert", "", "certPath")
	flag.StringVar(&keyPath, "key", "", "keyPath")
	flag.StringVar(&caPath, "ca", "", "caPath")
	flag.StringVar(&name, "name", "test", "Name")
	flag.StringVar(&apiKey, "apikey", "grpcnginx", "apiKey")
	flag.Parse()
}

func main() {
	if grpcURL != "" {
		log.Println("new conn:", grpcURL)
		var err error
		conn, err = grpc.Dial(grpcURL)
		if err != nil {
			log.Fatalf("notls did not connect: %v", err)
			time.Sleep(5 * time.Second)
			panic("e")
		}
	} else {
		conn = gtls.InitSSLGrpcWithEnv(nginx.GRPCServerURL, "NGINX_DNS_NAME", "nginx")
		log.Println("new tls conn:", grpcURL)
	}

	c := nginx.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	items := &nginx.MultSSLSetRequest{}
	items.Apikey = apiKey
	items.Item = make([]*nginx.SSLSetRequest, 0)
	log.Println("apikey:", items.Apikey)
	log.Println("ct read:", certPath)
	log.Println("ky read:", keyPath)
	if cert, err := ioutil.ReadFile(certPath); err != nil {
		log.Fatal(err)
	} else if key, err := ioutil.ReadFile(keyPath); err != nil {
		log.Fatal(err)
	} else {
		var ca []byte
		var err error
		if len(caPath) > 0 {
			log.Println("ca read:", caPath)
			if ca, err = ioutil.ReadFile(caPath); err != nil {
				log.Fatal(err)
			}
		}
		items.Item = append(items.Item, &nginx.SSLSetRequest{
			Directory: "",
			CertName:  name,
			KeyName:   name,
			Cert:      cert,
			Key:       key,
			Ca:        ca,
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
