// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package grpc_test

import (
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/pkg/uuid"
	"github.com/mainflux/mainflux/things"
	grpcapi "github.com/mainflux/mainflux/things/api/auth/grpc"
	"github.com/mainflux/mainflux/things/mocks"
	"github.com/opentracing/opentracing-go/mocktracer"
	"google.golang.org/grpc"
)

const (
	port  = 7000
	token = "token"
	wrong = "wrong"
	email = "john.doe@email.com"
)

var svc things.Service

func TestMain(m *testing.M) {
	serverErr := make(chan error)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("got unexpected error while creating new listerner: %s", err)
	}

	svc = newService(map[string]string{token: email})
	server := grpc.NewServer()
	mainflux.RegisterThingsServiceServer(server, grpcapi.NewServer(mocktracer.New(), svc))

	// Start gRPC server in detached mode.
	go func() {
		serverErr <- server.Serve(listener)
	}()

	code := m.Run()

	server.GracefulStop()
	err = <-serverErr
	if err != nil {
		log.Fatalln("gPRC Server Terminated : ", err)
	}
	close(serverErr)
	os.Exit(code)
}

func newService(tokens map[string]string) things.Service {
	policies := []mocks.MockSubjectSet{{Object: "users", Relation: "member"}}
	auth := mocks.NewAuthService(tokens, map[string][]mocks.MockSubjectSet{email: policies})
	conns := make(chan mocks.Connection)
	thingsRepo := mocks.NewThingRepository(conns)
	channelsRepo := mocks.NewChannelRepository(thingsRepo, conns)
	chanCache := mocks.NewChannelCache()
	thingCache := mocks.NewThingCache()
	idProvider := uuid.NewMock()

	return things.New(auth, thingsRepo, channelsRepo, chanCache, thingCache, idProvider)
}
