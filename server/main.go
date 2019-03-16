package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strings"

	pt "github.com/hkloudou/rpc_nginx/proto"
	"google.golang.org/grpc"
)

var port = ":5000"
var _version_ = "debug"
var _branch_ = ""
var _commitId_ = ""
var _buildTime_ = ""
var _appName_ = ""
var nginxSslPath = "/etc/nginx/certs/"
var apikey = "grpcnginx"
var signContainer = "nginx-proxy"

// server is used to implement helloworld.GreeterServer.
type server struct{}

//获得AccessKey
func (s *server) MultSSLSet(ctx context.Context, in *pt.MultSSLSetRequest) (*pt.SSLSetReply, error) {
	if in.GetApikey() != apikey {
		return nil, grpc.Errorf(1001, "apikey not equal")
	}
	for _, item := range in.GetItem() {
		p := path.Join(nginxSslPath, item.GetDirectory())
		os.MkdirAll(p, 0777)
		pathCert := path.Join(p, strings.Replace(item.GetCertName(), "..", "", -1))
		pathKey := path.Join(p, strings.Replace(item.GetKeyName(), "..", "", -1))
		if !strings.HasPrefix(nginxSslPath, pathCert) || !strings.HasPrefix(nginxSslPath, pathKey) {
			return nil, grpc.Errorf(1001, "path can not include ../")
		}
		ioutil.WriteFile(pathCert, item.GetCert(), 0655)
		ioutil.WriteFile(pathKey, item.GetKey(), 0655)
	}

	/*
		1.保存文件
		2.通知nginx
	*/

	return &pt.SSLSetReply{Ok: true}, nil
}

func init() {
	if os.Getenv("GRPC_PORT") != "" {
		port = os.Getenv("GRPC_PORT")
	}

	if os.Getenv("APIKEY") != "" {
		apikey = os.Getenv("APIKEY")
	}

	if os.Getenv("SIGN_CONTAINER") != "" {
		signContainer = os.Getenv("SIGN_CONTAINER")
	}

	if os.Getenv("NGINX_SSL_PATH") != "" {
		nginxSslPath = os.Getenv("NGINX_SSL_PATH")
	}

	log.Println("["+_appName_+"]", "init ...")
	log.Println("["+_appName_+"]", "version", _version_)
	log.Println("["+_appName_+"]", "branch", _branch_)
	log.Println("["+_appName_+"]", "commit id", _commitId_)
	log.Println("["+_appName_+"]", "build time", _buildTime_)
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("[%s] failed to listen: %v", _appName_, err)
	}
	s := grpc.NewServer()
	pt.RegisterGreeterServer(s, &server{})
	log.Println("["+_appName_+"]", "grpc listen:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[%s] failed to serve: %v", _appName_, err)
	}
}
