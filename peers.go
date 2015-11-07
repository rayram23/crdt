package crdt

import (
	"errors"
	"net"

	"google.golang.org/grpc"
)

// PeerPool consists of list of peers to talk to. Each peer
// will have a connection established to every other peer.
// PeerPool stores in the list of host:port addresses to talk
// to.
type PeerPool struct {
	peerAddrs []net.Addr
	conns     []*grpc.ClientConn
	clients   []*ReplicationClient
}

var (
	// ErrInvalidPeers is thrown if the list of peers is passed in
	// doesn't contain valid peers.
	ErrInvalidPeers = errors.New("invalid no. of peers found")
)

// NewPeerPool takes in a list of peers and options to initialize connections
// to those. It also tries to establish connections between the peers and any
// errors incurred during this process are returned.
func NewPeerPool(peers []string, opts ...grpc.DialOption) (*PeerPool, error) {
	if len(peers) <= 0 {
		return nil, ErrInvalidPeers
	}

	peerAddrs := make([]net.Addr, 0, len(peers))

	for _, p := range peers {
		tcpAddr, err := net.ResolveTCPAddr("tcp", p)
		if err != nil {
			return nil, err
		}

		peerAddrs = append(peerAddrs, tcpAddr)
	}

	clients := make([]*ReplicationClient, 0, len(peerAddrs))
	conns := make([]*grpc.ClientConn, 0, len(peerAddrs))

	for _, addr := range peerAddrs {
		conn, err := grpc.Dial(addr.String(), opts...)
		if err != nil {
			return nil, err
		}

		conns = append(conns, conn)
		clients = append(clients, NewReplicationClient(conn, addr.String()))
	}

	return &PeerPool{
		peerAddrs: peerAddrs,
		conns:     conns,
		clients:   clients,
	}, nil
}

// Addrs returns the list of unique addresses of each peer.
func (p *PeerPool) Addrs() []net.Addr {
	return p.peerAddrs
}

// Clients returns list of clients.
func (p *PeerPool) Clients() []*ReplicationClient {
	return p.clients
}

// CloseAll attempts to close established connections to all currently
// connected peers. If there are any errors in closing connection to any
// of the peers, then those failed connections are kept for future retries
// and the last seen non-nil error is returned.
func (p *PeerPool) CloseAll() error {
	var err error
	existingConns := make([]*grpc.ClientConn, 0, len(p.conns))

	for _, conn := range p.conns {
		if er := conn.Close(); er != nil {
			err = er
			existingConns = append(existingConns, conn)
		}
	}

	p.conns = existingConns
	return err
}
