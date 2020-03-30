package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mingcheng/proxypool/model"
	rpc "github.com/mingcheng/proxypool/protobuf"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mingcheng/proxypool"
)

type ProxyPoolRPCServer struct {
	rpc.ProxyPoolServer
}

func (p ProxyPoolRPCServer) Add(_ context.Context, proxy *rpc.Proxy) (*empty.Empty, error) {
	proxypool.Add(model.Proxy{Proxy: *proxy})
	return &empty.Empty{}, nil
}

func (p ProxyPoolRPCServer) Random(context.Context, *empty.Empty) (*rpc.Proxy, error) {
	return &proxypool.Random().Proxy, nil
}

func (p ProxyPoolRPCServer) All(context.Context, *empty.Empty) (*rpc.Proxies, error) {
	proxies := proxypool.All()

	if len(proxies) > 0 {
		var rpcProxies = &rpc.Proxies{
			Counts: uint64(len(proxies)),
		}
		for _, v := range proxies {
			rpcProxies.Proxies = append(rpcProxies.Proxies, &v.Proxy)
		}

		return rpcProxies, nil
	} else {
		return nil, fmt.Errorf("no suitable proxy found")
	}
}

func main() {
	config := proxypool.Config{
		FetchInterval:   15 * time.Minute,
		CheckInterval:   2 * time.Minute,
		CheckConcurrent: 10,
	}

	go proxypool.Start(config)
	defer proxypool.Stop()

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer()
	rpc.RegisterProxyPoolServer(s, ProxyPoolRPCServer{})
	go s.Serve(lis)

	r := gin.Default()
	r.GET("/all", func(c *gin.Context) {
		proxies := proxypool.All()
		if len(proxies) > 0 {
			var rpcProxies = &rpc.Proxies{
				Counts: uint64(len(proxies)),
			}
			for _, v := range proxies {
				rpcProxies.Proxies = append(rpcProxies.Proxies, &v.Proxy)
			}

			m := &jsonpb.Marshaler{}
			s, _ := m.MarshalToString(rpcProxies)
			c.Header("Content-Type", "application/json")
			c.String(http.StatusOK, s)
		} else {
			c.String(http.StatusNotFound, "no suitable proxy found")
		}
	})

	r.GET("/random", func(c *gin.Context) {
		if proxies := proxypool.Random(); proxies != nil {
			m := &jsonpb.Marshaler{}
			s, _ := m.MarshalToString(&proxies.Proxy)
			c.Header("Content-Type", "application/json")
			c.String(http.StatusOK, s)
		} else {
			c.String(http.StatusNotFound, "no suitable proxy found")
		}
	})

	// Start HTTP Server
	_ = r.Run(":8080")
}
