// Package router_address implements the I2P RouterAddress common data structure
package router_address

import (
	"errors"

	. "github.com/go-i2p/go-i2p/lib/common/data"
	log "github.com/sirupsen/logrus"
)

// Minimum number of bytes in a valid RouterAddress
const (
	ROUTER_ADDRESS_MIN_SIZE = 9
)

/*
[RouterAddress]
Accurate for version 0.9.49

Description
This structure defines the means to contact a router through a transport protocol.

Contents
1 byte Integer defining the relative cost of using the address, where 0 is free and 255 is
expensive, followed by the expiration Date after which the address should not be used,
of if null, the address never expires. After that comes a String defining the transport
protocol this router address uses. Finally there is a Mapping containing all of the
transport specific options necessary to establish the connection, such as IP address,
port number, email address, URL, etc.


+----+----+----+----+----+----+----+----+
|cost|           expiration
+----+----+----+----+----+----+----+----+
     |        transport_style           |
+----+----+----+----+-//-+----+----+----+
|                                       |
+                                       +
|               options                 |
~                                       ~
~                                       ~
|                                       |
+----+----+----+----+----+----+----+----+

cost :: Integer
        length -> 1 byte

        case 0 -> free
        case 255 -> expensive

expiration :: Date (must be all zeros, see notes below)
              length -> 8 bytes

              case null -> never expires

transport_style :: String
                   length -> 1-256 bytes

options :: Mapping
*/

// RouterAddress is the represenation of an I2P RouterAddress.
//
// https://geti2p.net/spec/common-structures#routeraddress
type RouterAddress struct {
	cost            *Integer
	expiration      *Date
	transport_style *I2PString
	options         *Mapping
}

// Bytes returns the router address as a []byte.
func (router_address RouterAddress) Bytes() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, router_address.cost.Bytes()...)
	bytes = append(bytes, router_address.expiration.Bytes()...)
	strData, err := router_address.transport_style.Data()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("RouterAddress.Bytes: error getting transport_style bytes")
	} else {
		bytes = append(bytes, strData...)
	}
	bytes = append(bytes, router_address.options.Data()...)
	return bytes
}

// Cost returns the cost for this RouterAddress as a Go integer.
func (router_address RouterAddress) Cost() int {
	return router_address.cost.Int()
}

// Expiration returns the expiration for this RouterAddress as an I2P Date.
func (router_address RouterAddress) Expiration() Date {
	return *router_address.expiration
}

// TransportStyle returns the transport style for this RouterAddress as an I2PString.
func (router_address RouterAddress) TransportStyle() I2PString {
	return *router_address.transport_style
}

// GetOption returns the value of the option specified by the key
func (router_address RouterAddress) GetOption(key I2PString) I2PString {
	return router_address.Options().Values().Get(key)
}

func (router_address RouterAddress) Host() I2PString {
	host, _ := ToI2PString("host")
	return router_address.GetOption(host)
}

// Options returns the options for this RouterAddress as an I2P Mapping.
func (router_address RouterAddress) Options() Mapping {
	return *router_address.options
}

// Check if the RouterAddress is empty or if it is too small to contain valid data.
func (router_address RouterAddress) checkValid() (err error, exit bool) {
	return
}

// ReadRouterAddress returns RouterAddress from a []byte.
// The remaining bytes after the specified length are also returned.
// Returns a list of errors that occurred during parsing.
func ReadRouterAddress(data []byte) (router_address RouterAddress, remainder []byte, err error) {
	if len(data) == 0 {
		log.WithField("at", "(RouterAddress) ReadRouterAddress").Error("error parsing RouterAddress: no data")
		err = errors.New("error parsing RouterAddress: no data")
		return
	}
	router_address.cost, remainder, err = NewInteger(data, 1)
	if err != nil {
		log.WithFields(log.Fields{
			"at":     "(RouterAddress) ReadNewRouterAddress",
			"reason": "error parsing cost",
		}).Warn("error parsing RouterAddress")
	}
	router_address.expiration, remainder, err = NewDate(remainder)
	if err != nil {
		log.WithFields(log.Fields{
			"at":     "(RouterAddress) ReadNewRouterAddress",
			"reason": "error parsing expiration",
		}).Error("error parsing RouterAddress")
	}
	router_address.transport_style, remainder, err = NewI2PString(remainder)
	if err != nil {
		log.WithFields(log.Fields{
			"at":     "(RouterAddress) ReadNewRouterAddress",
			"reason": "error parsing transport_style",
		}).Error("error parsing RouterAddress")
	}
	var errs []error
	router_address.options, remainder, errs = NewMapping(remainder)
	for _, err := range errs {
		log.WithFields(log.Fields{
			"at":     "(RouterAddress) ReadNewRouterAddress",
			"reason": "error parsing options",
			"error":  err,
		}).Error("error parsing RozuterAddress")
	}
	return
}

// NewRouterAddress creates a new *RouterAddress from []byte using ReadRouterAddress.
// Returns a pointer to RouterAddress unlike ReadRouterAddress.
func NewRouterAddress(data []byte) (router_address *RouterAddress, remainder []byte, err error) {
	objrouteraddress, remainder, err := ReadRouterAddress(data)
	router_address = &objrouteraddress
	return
}
