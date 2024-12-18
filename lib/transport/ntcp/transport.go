package ntcp

/**
 * https://geti2p.net/spec/ntcp2
**/

import (
	"fmt"
	"net"

	"github.com/go-i2p/go-i2p/lib/common/router_info"
	"github.com/go-i2p/go-i2p/lib/transport"
	"github.com/go-i2p/go-i2p/lib/transport/noise"
)

const (
	NTCP_PROTOCOL_VERSION = 2
	NTCP_PROTOCOL_NAME    = "NTCP2"
	NTCP_MESSAGE_MAX_SIZE = 65537
)

var exampleNTCPTransport transport.Transport = &Transport{}

// Transport is an ntcp2 transport implementing transport.Transport interface
type Transport struct {
	*noise.NoiseTransport
}

func (t *Transport) Name() string {
	return NTCP_PROTOCOL_NAME
}

func (t *Transport) Compatible(routerInfo router_info.RouterInfo) bool {
	// Check if the router info contains NTCP2 address and capabilities
	addresses := routerInfo.RouterAddresses()
	for _, addr := range addresses {
		transportStyle, err := addr.TransportStyle().Data()
		if err != nil {
			continue
		}
		if transportStyle == NTCP_PROTOCOL_NAME {
			return true
		}
	}
	return false
}

func (t *Transport) GetSession(routerInfo router_info.RouterInfo) (transport.TransportSession, error) {
	// Create new NTCP2 session
	session, err := NewNTCP2Session(routerInfo)
	if err != nil {
		return nil, err
	}

	// Perform handshake
	if err := session.Handshake(routerInfo); err != nil {
		return nil, err
	}

	return session, nil
}

func (t *Transport) Accept() (net.Conn, error) {
	conn, err := t.NoiseTransport.Accept()
	if err != nil {
		return nil, err
	}
	// check if remote router address contains a compatible transport
	// first get the RemoteAddr
	remoteAddr := conn.LocalAddr()
	// then check if it's a router address
	routerAddr, ok := remoteAddr.(*router_info.RouterInfo)
	if !ok {
		return nil, fmt.Errorf("remote address is not a router address")
	}
	// then check if it's compatible
	if !t.Compatible(*routerAddr) {
		return nil, fmt.Errorf("remote router address is not compatible with NTCP2")
	}
	// Wrap connection with NTCP2 session
	session, err := NewNTCP2Session(remoteAddr.(router_info.RouterInfo)) // nil for incoming connections
	if err != nil {
		conn.Close()
		return nil, err
	}

	return session, nil
}
