# ntcp
--
    import "github.com/go-i2p/go-i2p/lib/transport/ntcp"


## Usage

```go
const (
	NOISE_DH_CURVE25519 = 1

	NOISE_CIPHER_CHACHAPOLY = 1
	NOISE_CIPHER_AESGCM     = 2

	NOISE_HASH_SHA256 = 3

	NOISE_PATTERN_XK = 11

	MaxPayloadSize = math.MaxUint16 - 16 - uint16Size /*data len*/
)
```

```go
const (
	NTCP_PROTOCOL_VERSION = 2
	NTCP_PROTOCOL_NAME    = "NTCP2"
	NTCP_MESSAGE_MAX_SIZE = 65537
)
```

#### func  DeobfuscateEphemeralKey

```go
func DeobfuscateEphemeralKey(message []byte, aesKey *crypto.AESSymmetricKey) ([]byte, error)
```
DeobfuscateEphemeralKey decrypts the ephemeral public key in the message using
AES-256-CBC without padding

#### func  ObfuscateEphemeralKey

```go
func ObfuscateEphemeralKey(message []byte, aesKey *crypto.AESSymmetricKey) ([]byte, error)
```
ObfuscateEphemeralKey encrypts the ephemeral public key in the message using
AES-256-CBC without padding

#### type Session

```go
type Session struct{}
```

Session implements TransportSession An established transport session

#### type Transport

```go
type Transport struct{}
```

Transport is an ntcp transport implementing transport.Transport interface
