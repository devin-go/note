package test

import (
	"encoding/gob"
	"github.com/BurntSushi/toml"
	etcd "github.com/go-kratos/kratos/contrib/config/etcd/v2"
	"github.com/go-kratos/kratos/v2/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"strings"
	"testing"
	"time"
)

func printConfig(t *testing.T, name string, cf config.Config) {
	ngs := &NGServer{}
	v := cf.Value(name)
	err := v.Scan(ngs)
	t.Logf("%+v, error:%s", ngs, err)
}

type NGServer struct {
	Addr             []string `toml:"addr"`
	NgAddr           []string `toml:"ng_addr"`
	HttpAddr         string   `toml:"http_addr"`
	Insecure         bool     `toml:"insecure" json:"insecure"`
	HttpRoutePaths   []string `toml:"http_route_paths"`
	HttpPathRewrites []string `toml:"http_path_rewrites"`
	TcpPort          int      `toml:"tcp_port" json:"tcp_port"`
}

func TestKratosConfig(t *testing.T) {
	gob.Register(&NGServer{})
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	sc, err := etcd.New(client, etcd.WithPath("/youdu/config/server/"), etcd.WithPrefix(true))
	if err != nil {
		t.Fatal(err)
	}
	cf := config.New(
		config.WithSource(sc),
		config.WithDecoder(func(value *config.KeyValue, m map[string]interface{}) error {
			ngs := &NGServer{}
			k := strings.ReplaceAll(value.Key, "/youdu/config/server/", "")
			m[k] = ngs
			err := toml.Unmarshal(value.Value, ngs)
			return err
		}))

	cf.Load()
	printConfig(t, "ng-apps", cf)
	printConfig(t, "ng-auth", cf)
	cf.Watch("ng-apps", func(s string, value config.Value) {
		t.Logf("s:%s, value:%+v", s, value.Load())
	})

	time.Sleep(30 * time.Second)
	printConfig(t, "ng-apps", cf)
}
