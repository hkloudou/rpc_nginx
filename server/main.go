package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	nginx "github.com/hkloudou/rpc_nginx/nginx"
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
var endpoint = "unix:///var/run/docker.sock"

// server is used to implement helloworld.GreeterServer.
type server struct{}

var puberr error
var c *docker.Client

func init() {
	if client, err := NewDockerClient(endpoint); err != nil {
		puberr = grpc.Errorf(1005, "NewDockerClient error:%s", err)
	} else {
		c = client
	}
}

//获得AccessKey
func (s *server) MultSSLSet(ctx context.Context, in *nginx.MultSSLSetRequest) (*nginx.SSLSetReply, error) {
	if in.GetApikey() != apikey {
		return nil, grpc.Errorf(1001, "apikey not equal")
	}
	for _, item := range in.GetItem() {
		p := path.Join(nginxSslPath, item.GetDirectory())
		os.MkdirAll(p, 0777)
		if item.GetCertName() == "" || item.GetKeyName() == "" {
			return nil, grpc.Errorf(1002, "cert key must not be empty")
		}
		pathCert := path.Join(p, strings.Replace(item.GetCertName(), "..", "", -1)+".crt")
		pathKey := path.Join(p, strings.Replace(item.GetKeyName(), "..", "", -1)+".key")
		pathcaCrt := path.Join(p, strings.Replace(item.GetKeyName(), "..", "", -1)+"_ca.key")
		log.Println("pathCert", pathCert)
		log.Println("pathKey", pathKey)
		log.Println("pathcaCrt", pathcaCrt)
		if !strings.HasPrefix(pathCert, nginxSslPath) || !strings.HasPrefix(pathKey, nginxSslPath) || !strings.HasPrefix(pathcaKey, nginxSslPath) {
			return nil, grpc.Errorf(1002, "path can not include ../")
		}
		if err := ioutil.WriteFile(pathCert, item.GetCert(), 0655); err != nil {
			return nil, grpc.Errorf(1003, "error write cert:%s", err)
		}
		if err := ioutil.WriteFile(pathKey, item.GetKey(), 0655); err != nil {
			return nil, grpc.Errorf(1004, "error write key:%s", err)
		}

		if len(item.GetCa()) > 0 {
			if err := ioutil.WriteFile(pathcaCrt, item.GetCa(), 0655); err != nil {
				return nil, grpc.Errorf(1005, "error write ca:%s", err)
			}
		}
	}

	killOpts := docker.KillContainerOptions{
		ID:     signContainer,
		Signal: docker.SIGHUP,
	}
	if puberr != nil {
		return nil, puberr
	} else if err := c.KillContainer(killOpts); err != nil {
		return nil, grpc.Errorf(1006, "Error sending signal to container: %s", err)
	}
	log.Println("ok")
	/*
		1.保存文件
		2.通知nginx
	*/
	return &nginx.SSLSetReply{Ok: true}, nil
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

	if os.Getenv("ENDPOINT") != "" {
		endpoint = os.Getenv("ENDPOINT")
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
	nginx.RegisterGreeterServer(s, &server{})
	log.Println("["+_appName_+"]", "grpc listen:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[%s] failed to serve: %v", _appName_, err)
	}
}
