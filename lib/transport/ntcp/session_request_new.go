package ntcp

import (
	"crypto/rand"
	"math/big"
	"net"

	"github.com/go-i2p/go-i2p/lib/common/data"
	"github.com/go-i2p/go-i2p/lib/transport/ntcp/handshake"
	"github.com/go-i2p/go-i2p/lib/transport/ntcp/messages"
	"github.com/samber/oops"
)

/*
SessionRequestProcessor implements NTCP2 Message 1 (SessionRequest):
1. Create session request message with options block (version, padding length, etc.)
2. Set timeout deadline for the connection
3. Obfuscate ephemeral key (X) using AES with Bob's router hash as key
4. Encrypt options block using ChaCha20-Poly1305
5. Assemble final message: obfuscated X + encrypted options + padding
6. Write complete message to connection
*/
type SessionRequestProcessor struct {
	*NTCP2Session
}

// EncryptPayload encrypts the payload portion of the message
func (p *SessionRequestProcessor) EncryptPayload(
	message messages.Message,
	obfuscatedKey []byte,
	hs *handshake.HandshakeState,
) ([]byte, error) {
	req, ok := message.(*messages.SessionRequest)
	if !ok {
		return nil, oops.Errorf("expected SessionRequest message")
	}

	return p.NTCP2Session.encryptSessionRequestOptions(req, obfuscatedKey)
}

// MessageType implements handshake.HandshakeMessageProcessor.
func (s *SessionRequestProcessor) MessageType() messages.MessageType {
	return messages.MessageTypeSessionRequest
}

// ProcessMessage implements handshake.HandshakeMessageProcessor.
func (s *SessionRequestProcessor) ProcessMessage(message messages.Message, hs *handshake.HandshakeState) error {
	panic("unimplemented")
}

// ReadMessage reads a SessionRequest message from the connection
func (p *SessionRequestProcessor) ReadMessage(conn net.Conn, hs *handshake.HandshakeState) (messages.Message, error) {
	// 1. Read ephemeral key
	obfuscatedX, err := p.NTCP2Session.readEphemeralKey(conn)
	if err != nil {
		return nil, err
	}

	// 2. Process ephemeral key
	deobfuscatedX, err := p.NTCP2Session.processEphemeralKey(obfuscatedX, hs)
	if err != nil {
		return nil, err
	}

	// 3. Read options block
	encryptedOptions, err := p.NTCP2Session.readOptionsBlock(conn)
	if err != nil {
		return nil, err
	}

	// 4. Process options block
	options, err := p.NTCP2Session.processOptionsBlock(encryptedOptions, obfuscatedX, deobfuscatedX, hs)
	if err != nil {
		return nil, err
	}

	// 5. Read padding if present
	paddingLen := options.PaddingLength.Int()
	if paddingLen > 0 {
		if err := p.NTCP2Session.readAndValidatePadding(conn, paddingLen); err != nil {
			return nil, err
		}
	}

	// Construct the full message
	return &messages.SessionRequest{
		XContent: [32]byte{}, // We've already processed this
		Options:  *options,
		Padding:  make([]byte, paddingLen), // Padding content doesn't matter after validation
	}, nil
}

// CreateMessage implements HandshakeMessageProcessor.
func (s *SessionRequestProcessor) CreateMessage(hs *handshake.HandshakeState) (messages.Message, error) {
	// Get our ephemeral key pair
	ephemeralKey := make([]byte, 32)
	if _, err := rand.Read(ephemeralKey); err != nil {
		return nil, err
	}

	// Add random padding (implementation specific)
	randomInt, err := rand.Int(rand.Reader, big.NewInt(16))
	if err != nil {
		return nil, err
	}

	padding := make([]byte, randomInt.Int64()) // Up to 16 bytes of padding
	if err != nil {
		return nil, err
	}

	netId, err := data.NewIntegerFromInt(2, 1)
	if err != nil {
		return nil, err
	}
	version, err := data.NewIntegerFromInt(2, 1)
	if err != nil {
		return nil, err
	}
	paddingLen, _, err := data.NewInteger([]byte{byte(len(padding))}, 1)
	if err != nil {
		return nil, err
	}
	//message3Part2Len, err := data.NewInteger()
	//if err != nil {
	//	return nil, err
	//}
	timestamp, err := data.DateFromTime(s.GetCurrentTime())
	if err != nil {
		return nil, err
	}
	requestOptions := &messages.RequestOptions{
		NetworkID:       netId,
		ProtocolVersion: version,
		PaddingLength:   paddingLen,
		// Message3Part2Length: ,
		Timestamp: timestamp,
	}

	return &messages.SessionRequest{
		XContent: [32]byte(ephemeralKey),
		Options:  *requestOptions,
		Padding:  padding,
	}, nil
}

// GetPadding retrieves padding from a message
func (p *SessionRequestProcessor) GetPadding(message messages.Message) []byte {
	req, ok := message.(*messages.SessionRequest)
	if !ok {
		return nil
	}

	return req.Padding
}

// ObfuscateKey obfuscates the ephemeral key for transmission
func (p *SessionRequestProcessor) ObfuscateKey(message messages.Message, handshake *handshake.HandshakeState) ([]byte, error) {
	req, ok := message.(*messages.SessionRequest)
	if !ok {
		return nil, oops.Errorf("expected SessionRequest message")
	}

	return p.NTCP2Session.ObfuscateEphemeral(req.XContent[:])
}

var _ handshake.HandshakeMessageProcessor = (*SessionRequestProcessor)(nil)
