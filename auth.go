//go:generate godocdown -o README.md

package api

import (
	"context"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"reflect"
)

type Authorizer struct {
	Providers map[string]oauth2.TokenSource
}

func (s *Authorizer) GetClient(ctx context.Context, key string) *http.Client {
	return oauth2.NewClient(ctx, s.Providers[key])
}

func (s *Authorizer) Do(ctx context.Context, key string, req *http.Request) ([]byte, error) {
	resp, err := oauth2.NewClient(ctx, s.Providers[key]).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (s *Authorizer) Put(key string, tok oauth2.TokenSource) {
	s.Providers[key] = tok
}

func (s *Authorizer) Exists(key string) bool {
	t := s.Providers[key]
	if t == nil {
		return false
	}
	return true
}

func (s *Authorizer) Length() int {
	return len(s.Providers)
}

func (s *Authorizer) DeepEqual(y interface{}) bool {
	return reflect.DeepEqual(s, y)
}

func (s *Authorizer) Validate(fn func(a *Authorizer) error) error {
	return fn(s)
}
