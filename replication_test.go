package crdt

import (
	"errors"
	"net"
	"strconv"
	"testing"
	"time"

	"google.golang.org/grpc"
)

const (
	NoOfPeers = 3
)

type ReplicationTestServer struct {
	listener net.Listener

	srv *grpc.Server

	rsrv *ReplicationServer
}

func NewTestReplicationServer() (*ReplicationTestServer, error) {
	addr, err := net.ResolveTCPAddr("tcp", ":0")
	if err != nil {
		return nil, err
	}

	lner, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}

	srv := grpc.NewServer()

	return &ReplicationTestServer{
		srv:      srv,
		listener: lner,
	}, nil
}

func TestReplicationPing(t *testing.T) {
	t.Parallel()

	testServers := make([]*ReplicationTestServer, 0, NoOfPeers)

	for i := 0; i < NoOfPeers; i++ {
		srv, err := NewTestReplicationServer()
		if err != nil {
			t.Fatalf("cannot create test server: %s", err)
		}

		rsrv := NewReplicationServer(strconv.Itoa(i+1), 3*time.Second)
		srv.srv.RegisterService(&_ReplicationTransport_serviceDesc, rsrv)

		go func() {
			if err := srv.srv.Serve(srv.listener); err != nil {
				t.Errorf("cannot start server:%d %q", i+1, err)
			}
		}()

		testServers = append(testServers, srv)
	}

	defer func() {
		for _, srv := range testServers {
			srv.srv.Stop()
		}
	}()

	clients := make([]*ReplicationClient, 0, len(testServers))

	for i, srv := range testServers {
		conn, err := grpc.Dial(srv.listener.Addr().String(), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("cannot create client[%d]: %s", i, err)
		}
		// These defers will not be called until the entire test
		// completes running, which is exactly what we want.
		defer conn.Close()

		clients = append(clients, NewReplicationClient(conn, strconv.Itoa(i+1)))
	}

	var attempts int32 = 5

	for _, client := range clients {
		var count = client.count

		for i := 0; i < int(attempts); i++ {
			if err := client.Ping(); err != nil {
				t.Errorf("cannot ping client:%d %q", client.Id(), err)
			}
		}

		if count += attempts; client.count != count {
			t.Errorf("expected ping count: %d, actual: %d", count, client.count)
		}
	}
}

type crdtRequest struct{}

func (c crdtRequest) GetQuery() ([]byte, error) {
	return []byte("query"), nil
}

var _ CRDTRequest = crdtRequest{}

type crdtBackend struct{}

func (c crdtBackend) Query(query []byte) (<-chan []byte, <-chan error) {
	var noOfResults = 5
	dChan, eChan := make(chan []byte, noOfResults), make(chan error, 1)

	defer func() {
		close(dChan)
		close(eChan)
	}()

	switch string(query) {
	case "query":
		for i := 0; i < noOfResults; i++ {
			dChan <- []byte("result-" + strconv.Itoa(i+1))
		}
		eChan <- nil
	default:
		eChan <- errors.New("something's borked")
	}

	return dChan, eChan
}

var _ CRDTBackend = crdtBackend{}

func TestReplicationQuery(t *testing.T) {
	t.Parallel()

	backend := crdtBackend{}
	request := crdtRequest{}

	testServers := make([]*ReplicationTestServer, 0, NoOfPeers)

	for i := 0; i < NoOfPeers; i++ {
		srv, err := NewTestReplicationServer()
		if err != nil {
			t.Fatalf("cannot create test server: %s", err)
		}

		rsrv := NewReplicationServer(strconv.Itoa(i+1), 3*time.Second)
		rsrv.Backend = backend

		srv.srv.RegisterService(&_ReplicationTransport_serviceDesc, rsrv)

		go func() {
			if err := srv.srv.Serve(srv.listener); err != nil {
				t.Errorf("cannot start server:%d %q", i+1, err)
			}
		}()

		testServers = append(testServers, srv)
	}

	defer func() {
		for _, srv := range testServers {
			srv.srv.Stop()
		}
	}()

	clients := make([]*ReplicationClient, 0, len(testServers))

	for i, srv := range testServers {
		conn, err := grpc.Dial(srv.listener.Addr().String(), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("cannot create client[%d]: %s", i, err)
		}
		// These defers will not be called until the entire test
		// completes running, which is exactly what we want.
		defer conn.Close()

		clients = append(clients, NewReplicationClient(conn, strconv.Itoa(i+1)))
	}

	for _, client := range clients {
		dChan, eChan := client.Query(request)
		for d := range dChan {
			t.Logf("retrieved: %q", d)
		}
		for e := range eChan {
			if e != nil {
				t.Errorf("error occurred while querying: %s", e)
			}
		}
	}
}
