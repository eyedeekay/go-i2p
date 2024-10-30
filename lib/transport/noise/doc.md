# noise
--
    import "github.com/go-i2p/go-i2p/lib/transport/noise"


## Usage

```go
const (
	NOISE_DH_CURVE25519 = 1

	NOISE_CIPHER_CHACHAPOLY = 1
	NOISE_CIPHER_AESGCM     = 2

	NOISE_HASH_SHA256 = 3

	NOISE_PATTERN_XK = 11

	MaxPayloadSize = 65537
)
```

```go
var ExampleNoiseListener net.Listener = exampleNoiseTransport
```
ExampleNoiseListener is not a real Noise Listener, do not use it. It is exported
so that it can be confirmed that the transport implements net.Listener

```go
var (
	ExampleNoiseSession net.Conn = exampleNoiseSession.(*NoiseSession)
)
```

#### func  ComposeInitiatorHandshakeMessage

```go
func ComposeInitiatorHandshakeMessage(s noise.DHKey, rs []byte, payload []byte, ePrivate []byte) (negData, msg []byte, state *noise.HandshakeState, err error)
```

#### func  ComposeReceiverHandshakeMessage

```go
func ComposeReceiverHandshakeMessage(s noise.DHKey, rs []byte, payload []byte, ePrivate []byte) (negData, msg []byte, state *noise.HandshakeState, err error)
```

#### func  NewNoiseTransportSession

```go
func NewNoiseTransportSession(ri router_info.RouterInfo) (transport.TransportSession, error)
```

#### type NoiseSession

```go
type NoiseSession struct {
	router_info.RouterInfo
	*noise.CipherState
	sync.Mutex
	*sync.Cond
	*NoiseTransport // The parent transport, which "Dialed" the connection to the peer whith whom we established the session
	RecvQueue       *cb.Queue
	SendQueue       *cb.Queue
	SendKey         noise.DHKey
	RecvKey         noise.DHKey
	HandKey         noise.DHKey
	VerifyCallback  VerifyCallbackFunc

	Conn net.Conn
}
```


#### func (*NoiseSession) Close

```go
func (s *NoiseSession) Close() error
```

#### func (*NoiseSession) LocalAddr

```go
func (s *NoiseSession) LocalAddr() net.Addr
```

#### func (*NoiseSession) QueueSendI2NP

```go
func (s *NoiseSession) QueueSendI2NP(msg i2np.I2NPMessage)
```

#### func (*NoiseSession) Read

```go
func (c *NoiseSession) Read(b []byte) (int, error)
```

#### func (*NoiseSession) ReadNextI2NP

```go
func (s *NoiseSession) ReadNextI2NP() (i2np.I2NPMessage, error)
```

#### func (*NoiseSession) RemoteAddr

```go
func (noise_session *NoiseSession) RemoteAddr() net.Addr
```
RemoteAddr implements net.Conn

#### func (*NoiseSession) RunIncomingHandshake

```go
func (c *NoiseSession) RunIncomingHandshake() error
```

#### func (*NoiseSession) RunOutgoingHandshake

```go
func (c *NoiseSession) RunOutgoingHandshake() error
```

#### func (*NoiseSession) SendQueueSize

```go
func (s *NoiseSession) SendQueueSize() int
```

#### func (*NoiseSession) SetDeadline

```go
func (noise_session *NoiseSession) SetDeadline(t time.Time) error
```
SetDeadline implements net.Conn

#### func (*NoiseSession) SetReadDeadline

```go
func (noise_session *NoiseSession) SetReadDeadline(t time.Time) error
```
SetReadDeadline implements net.Conn

#### func (*NoiseSession) SetWriteDeadline

```go
func (noise_session *NoiseSession) SetWriteDeadline(t time.Time) error
```
SetWriteDeadline implements net.Conn

#### func (*NoiseSession) Write

```go
func (c *NoiseSession) Write(b []byte) (int, error)
```

#### type NoiseTransport

```go
type NoiseTransport struct {
	sync.Mutex
	router_identity.RouterIdentity
	*noise.CipherState
	Listener net.Listener
}
```


#### func  NewNoiseTransport

```go
func NewNoiseTransport(netSocket net.Listener) *NoiseTransport
```
NewNoiseTransport create a NoiseTransport using a supplied net.Listener

#### func  NewNoiseTransportSocket

```go
func NewNoiseTransportSocket() (*NoiseTransport, error)
```
NewNoiseTransportSocket creates a Noise transport socket with a random host and
port.

#### func (*NoiseTransport) Accept

```go
func (noopt *NoiseTransport) Accept() (net.Conn, error)
```
Accept a connection on a listening socket.

#### func (*NoiseTransport) Addr

```go
func (noopt *NoiseTransport) Addr() net.Addr
```
Addr of the transport, for now this is returning the IP:Port the transport is
listening on, but this might actually be the router identity

#### func (*NoiseTransport) Close

```go
func (noopt *NoiseTransport) Close() error
```
close the transport cleanly blocks until done returns an error if one happens

#### func (*NoiseTransport) Compatable

```go
func (noopt *NoiseTransport) Compatable(routerInfo router_info.RouterInfo) bool
```
Compatable return true if a routerInfo is compatable with this transport

#### func (*NoiseTransport) Compatible

```go
func (noopt *NoiseTransport) Compatible(routerInfo router_info.RouterInfo) bool
```

#### func (*NoiseTransport) GetSession

```go
func (noopt *NoiseTransport) GetSession(routerInfo router_info.RouterInfo) (transport.TransportSession, error)
```
Obtain a transport session with a router given its RouterInfo. If a session with
this router is NOT already made attempt to create one and block until made or
until an error happens returns an established TransportSession and nil on
success returns nil and an error on error

#### func (*NoiseTransport) Handshake

```go
func (c *NoiseTransport) Handshake(routerInfo router_info.RouterInfo) error
```

#### func (*NoiseTransport) Name

```go
func (noopt *NoiseTransport) Name() string
```
Name of the transport TYPE, in this case `noise`

#### func (*NoiseTransport) SetIdentity

```go
func (noopt *NoiseTransport) SetIdentity(ident router_identity.RouterIdentity) (err error)
```
SetIdentity will set the router identity for this transport. will bind if the
underlying socket is not already if the underlying socket is already bound
update the RouterIdentity returns any errors that happen if they do

#### type VerifyCallbackFunc

```go
type VerifyCallbackFunc func(publicKey []byte, data []byte) error
```
