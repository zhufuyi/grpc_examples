package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/resolver"
)

type server struct {
	Name    string `json:"name"`    // rpc server name
	Addr    string `json:"addr"`    // rpc server address
	Version string `json:"version"` // rpc server version
	Weight  int64  `json:"weight"`  // rpc server weight
}

func buildPrefix(info server) string {
	if info.Version == "" {
		return fmt.Sprintf("/%s/", info.Name)
	}
	return fmt.Sprintf("/%s/%s/", info.Name, info.Version)
}

func buildRegPath(info server) string {
	return fmt.Sprintf("%s%s", buildPrefix(info), info.Addr)
}

func ParseValue(value []byte) (server, error) {
	info := server{}
	if err := json.Unmarshal(value, &info); err != nil {
		return info, err
	}
	return info, nil
}

func splitPath(path string) (server, error) {
	info := server{}
	strs := strings.Split(path, "/")
	if len(strs) == 0 {
		return info, errors.New("invalid path")
	}
	info.Addr = strs[len(strs)-1]
	return info, nil
}

// exist helper function
func exist(l []resolver.Address, addr resolver.Address) bool {
	for i := range l {
		if l[i].Addr == addr.Addr {
			return true
		}
	}
	return false
}

// remove helper function
func remove(s []resolver.Address, addr resolver.Address) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr.Addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}

func buildResolverUrl(serverName string) string {
	return schema + ":///" + serverName
}
