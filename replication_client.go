package crdt

import (
	"io"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ReplicationClient is designed to be used by the CRDT server
// to issue pings and queries to its peers. Querying for information
// needs to pass its current state in and acquire the new server
// state in return.
type ReplicationClient struct {
	serverId string
	count    int32

	client ReplicationTransportClient
}

// CRDTRequest is a generic interface that translates a CRDT query
// into something that the other peer endpoint can understand.
type CRDTRequest interface {
	GetQuery() ([]byte, error)
}

// NewReplicationClient supplies the given connection to the newly created client.
// Each client is supposed to start on one of the peer servers and hence needs to
// have a unique id.
func NewReplicationClient(conn *grpc.ClientConn, serverId string) *ReplicationClient {
	return &ReplicationClient{
		serverId: serverId,
		client:   NewReplicationTransportClient(conn),
	}
}

// Id returns the unique id of the given server this client is a part of.
func (r *ReplicationClient) Id() string {
	return r.serverId
}

// Ping pings the server that is on the other end of this client.
func (r *ReplicationClient) Ping() error {
	// Currently we just increment our counter to keep a track
	// of issued ping commands. This behavior is not atomic
	// in nature and needs to be made thread-safe.
	r.count++

	_, err := r.client.Ping(context.Background(), &PingRequest{
		Id:        r.count,
		Timestamp: time.Now().Unix(),
	})
	return err
}

// Query takes the query that our crdt client gives us and uses that to
// query the server. It returns two channels, one for data and another
// one for error, respectively. It serves the data and errors (if any)
// asynchronously to the caller.
func (r *ReplicationClient) Query(req CRDTRequest) (<-chan []byte, <-chan error) {
	dataChan := make(chan []byte)
	errChan := make(chan error, 1)

	in, err := req.GetQuery()
	if err != nil {
		errChan <- err
		close(dataChan)
		close(errChan)

		return dataChan, errChan
	}

	go func() {
		defer func() {
			close(dataChan)
			close(errChan)
		}()

		conn, err := r.client.Query(context.Background(), &QueryRequest{
			Server: r.serverId,
			Query:  in,
		})
		if err != nil {
			errChan <- err
			return
		}

		for {
			res, err := conn.Recv()
			if err != nil {
				if err == io.EOF {
					err = nil
				}
				errChan <- err
				return
			}

			dataChan <- res.Result
		}
	}()

	return dataChan, errChan
}
