package crdt

import (
	"errors"
	"time"

	"golang.org/x/net/context"
)

var (
	// An ErrQueryTimeout is returned if the query times out
	// on the server.
	ErrQueryTimeout = errors.New("query timed out")
)

// ReplicationServer describes the server-side functionality a server
// should have. It needs a CRDTBackend that is able to suffice a Query
// issued by the client. If the query to CRDT Backend takes more than
// `queryTimeout` time to respond than we return an error.
type ReplicationServer struct {
	Id           string
	queryTimeout time.Duration

	Backend CRDTBackend
}

// NewReplicationServer creates an instance of server and attaches a given
// id to it along with the query timeout parameter.
func NewReplicationServer(Id string, timeout time.Duration) *ReplicationServer {
	return &ReplicationServer{
		Id:           Id,
		queryTimeout: timeout,
	}
}

// CRDTBackend is an interface that our CRDT backend should implement that takes
// in a query and returns the resulting response as data passed into the channel.
// If there is an error transmitting the data then error should be thrown in
// the error channel.
type CRDTBackend interface {
	Query([]byte) (<-chan []byte, <-chan error)
}

// ReplicationServer should implement the ReplicationTransportServer interface.
var _ ReplicationTransportServer = &ReplicationServer{}

// Ping returns back the result of a ping received from a client. It returns back
// the same id that the client has sent along with the time when it was received.
func (r *ReplicationServer) Ping(ctx context.Context, req *PingRequest) (*PingResult, error) {
	return &PingResult{
		Id:        req.Id,
		Timestamp: time.Now().Unix(),
	}, nil
}

// Query takes in a query from the client, issues it to the backend and streams
// the received responses back to the client.
// If the backend doesn't send its responses within a given amount of time then
// it automatically hangs up and returns an error to the client.
func (r *ReplicationServer) Query(req *QueryRequest, srv ReplicationTransport_QueryServer) error {
	dataChan, errChan := r.Backend.Query(req.Query)

	for {
		select {
		case data, open := <-dataChan:
			if !open {
				return nil
			}

			if err := srv.Send(&QueryResult{
				Result: data,
			}); err != nil {
				return err
			}
		case err := <-errChan:
			if err != nil {
				return err
			}
		case <-time.After(r.queryTimeout):
			return ErrQueryTimeout
		}
	}

	return nil
}
