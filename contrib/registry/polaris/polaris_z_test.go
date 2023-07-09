// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package polaris

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/config"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
)

// TestRegistry TestRegistryManyService
func TestRegistry(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	conf.GetGlobal().GetStatReporter().SetEnable(false)
	conf.Consumer.LocalCache.SetPersistDir(os.TempDir() + "/polaris-registry/backup")
	if err := api.SetLoggersDir(os.TempDir() + "/polaris-registry/log"); err != nil {
		t.Fatal(err)
	}

	r := NewWithConfig(
		conf,
		WithTimeout(time.Second*10),
		WithTTL(100),
	)

	svc := &gsvc.LocalService{
		Name:      "goframe-provider-0-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:9000"),
	}

	s, err := r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}

	if err = r.Deregister(context.Background(), s); err != nil {
		t.Fatal(err)
	}
}

// TestRegistryMany TestRegistryManyService
func TestRegistryMany(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	conf.GetGlobal().GetStatReporter().SetEnable(false)
	conf.Consumer.LocalCache.SetPersistDir(os.TempDir() + "/polaris-registry-many/backup")
	if err := api.SetLoggersDir(os.TempDir() + "/polaris-registry-many/log"); err != nil {
		t.Fatal(err)
	}

	r := NewWithConfig(
		conf,
		WithTimeout(time.Second*10),
		WithTTL(100),
	)

	svc := &gsvc.LocalService{
		Name:      "goframe-provider-1-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:9000"),
	}
	svc1 := &gsvc.LocalService{
		Name:      "goframe-provider-2-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:9001"),
	}
	svc2 := &gsvc.LocalService{
		Name:      "goframe-provider-3-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:9002"),
	}

	s0, err := r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}

	s1, err := r.Register(context.Background(), svc1)
	if err != nil {
		t.Fatal(err)
	}

	s2, err := r.Register(context.Background(), svc2)
	if err != nil {
		t.Fatal(err)
	}

	if err = r.Deregister(context.Background(), s0); err != nil {
		t.Fatal(err)
	}

	if err = r.Deregister(context.Background(), s1); err != nil {
		t.Fatal(err)
	}

	if err = r.Deregister(context.Background(), s2); err != nil {
		t.Fatal(err)
	}
}

// TestGetService Test GetService
func TestGetService(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	conf.GetGlobal().GetStatReporter().SetEnable(false)
	conf.Consumer.LocalCache.SetPersistDir(os.TempDir() + "/polaris-get-service/backup")
	if err := api.SetLoggersDir(os.TempDir() + "/polaris-get-service/log"); err != nil {
		t.Fatal(err)
	}

	r := NewWithConfig(
		conf,
		WithTimeout(time.Second*10),
		WithTTL(100),
	)

	svc := &gsvc.LocalService{
		Name:      "goframe-provider-4-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:9000"),
	}

	s, err := r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 1)
	serviceInstances, err := r.Search(context.Background(), gsvc.SearchInput{
		Prefix:   s.GetPrefix(),
		Name:     svc.Name,
		Version:  svc.Version,
		Metadata: svc.Metadata,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range serviceInstances {
		t.Log(instance)
	}

	if err = r.Deregister(context.Background(), s); err != nil {
		t.Fatal(err)
	}
}

// TestWatch Test Watch
func TestWatch(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	conf.GetGlobal().GetStatReporter().SetEnable(false)
	conf.Consumer.LocalCache.SetPersistDir(os.TempDir() + "/polaris-watch/backup")
	if err := api.SetLoggersDir(os.TempDir() + "/polaris-watch/log"); err != nil {
		t.Fatal(err)
	}
	r := NewWithConfig(
		conf,
		WithTimeout(time.Second*10),
		WithTTL(100),
	)

	svc := &gsvc.LocalService{
		Name:      "goframe-provider-5-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:9000"),
	}

	s := &Service{
		Service: svc,
	}

	watch, err := r.Watch(context.Background(), s.GetPrefix())
	if err != nil {
		t.Fatal(err)
	}

	s1, err := r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}
	// watch svc
	time.Sleep(time.Second * 1)

	// svc register, AddEvent
	next, err := watch.Proceed()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output one instance
		t.Log("Register Proceed service: ", instance)
	}

	if err = r.Deregister(context.Background(), s1); err != nil {
		t.Fatal(err)
	}

	// svc deregister, DeleteEvent
	next, err = watch.Proceed()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output nothing
		t.Log("Deregister Proceed service: ", instance)
	}

	if err = watch.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err = watch.Proceed(); err == nil {
		// if nil, stop failed
		t.Fatal()
	}
}

// BenchmarkRegister
func BenchmarkRegister(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
		conf.GetGlobal().GetStatReporter().SetEnable(false)
		conf.Consumer.LocalCache.SetPersistDir(os.TempDir() + "/polaris-registry/backup")
		if err := api.SetLoggersDir(os.TempDir() + "/polaris-registry/log"); err != nil {
			b.Fatal(err)
		}

		r := NewWithConfig(
			conf,
			WithTimeout(time.Second*10),
			WithTTL(100),
		)

		svc := &gsvc.LocalService{
			Name:      "goframe-provider-0-tcp",
			Version:   "test",
			Metadata:  map[string]interface{}{"app": "goframe", gsvc.MDProtocol: "tcp"},
			Endpoints: gsvc.NewEndpoints("127.0.0.1:9000"),
		}

		s, err := r.Register(context.Background(), svc)
		if err != nil {
			b.Fatal(err)
		}

		if err = r.Deregister(context.Background(), s); err != nil {
			b.Fatal(err)
		}
	}
}

// TestRegistryManyForEndpoints TestRegistryManyForEndpointsService
func TestRegistryManyForEndpoints(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	conf.GetGlobal().GetStatReporter().SetEnable(false)
	conf.Consumer.LocalCache.SetPersistDir(os.TempDir() + "/polaris-registry-many/backup")
	if err := api.SetLoggersDir(os.TempDir() + "/polaris-registry-many/log"); err != nil {
		t.Fatal(err)
	}

	r := NewWithConfig(
		conf,
		WithTimeout(time.Second*10),
		WithTTL(100),
	)

	var (
		serviceName   = "goframe-provider-tcp"
		version       = "latest"
		endpointOne   = "127.0.0.1:9000"
		endpointTwo   = "127.0.0.1:9001"
		endpointThree = "127.0.0.1:9002"
	)

	svc := &gsvc.LocalService{
		Name:      serviceName,
		Version:   version,
		Metadata:  map[string]interface{}{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints(endpointOne),
	}

	svc1 := &gsvc.LocalService{
		Name:      serviceName,
		Version:   version,
		Metadata:  map[string]interface{}{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints(endpointTwo),
	}

	svc2 := &gsvc.LocalService{
		Name:      serviceName,
		Version:   version,
		Metadata:  map[string]interface{}{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints(endpointThree),
	}

	s0, err := r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}

	s1, err := r.Register(context.Background(), svc1)
	if err != nil {
		t.Fatal(err)
	}

	s2, err := r.Register(context.Background(), svc2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Register service success sleep 1s")
	time.Sleep(time.Second * 1)
	// serviceName = "service-default-default-goframe-provider-tcp-latest"
	result, err := r.Search(context.Background(), gsvc.SearchInput{
		Name: serviceName,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Search service success size:", len(result))
	for i := 0; i < len(result); i++ {
		t.Log("Endpoints:", result[i].GetEndpoints().String())
		if !gstr.Contains(result[i].GetEndpoints().String(), endpointOne) {
			t.Fatal("endpointOne not found")
		}
		if !gstr.Contains(result[i].GetEndpoints().String(), endpointTwo) {
			t.Fatal("endpointTwo not found")
		}
		if !gstr.Contains(result[i].GetEndpoints().String(), endpointThree) {
			t.Fatal("endpointThree not found")
		}
	}
	t.Log("Search service success sleep 1s")
	time.Sleep(time.Second * 1)
	if err = r.Deregister(context.Background(), s0); err != nil {
		t.Fatal(err)
	}

	if err = r.Deregister(context.Background(), s1); err != nil {
		t.Fatal(err)
	}

	if err = r.Deregister(context.Background(), s2); err != nil {
		t.Fatal(err)
	}
	t.Log("Deregister success")
}
