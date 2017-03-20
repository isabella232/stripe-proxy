// Copyright © 2017 stripe-proxy authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxy

import (
	"encoding/binary"
)

type StripeResource int

// Note: these do not use iota so that they are stable through modifications
// of the list.
const (
	ResourceAll               StripeResource = 0
	ResourceBalance                          = 1
	ResourceCharges                          = 2
	ResourceCustomers                        = 3
	ResourceDisputes                         = 4
	ResourceEvents                           = 5
	ResourceFileUploads                      = 6
	ResourceRefunds                          = 7
	ResourceTokens                           = 8
	ResourceTransfers                        = 9
	ResourceTransferReversals                = 10
)

type Access int

// Note: these do not use iota so that they are stable through modifications
// of the list.
const (
	None      = 0
	Read      = 1
	Write     = 2
	ReadWrite = 3
)

type Permission struct {
	encoded uint32
}

func NewPermission(initialValue uint32) *Permission {
	return &Permission{initialValue}
}

func (p *Permission) MarshalBinary() ([]byte, error) {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, p.encoded)
	return bs, nil
}

func (p *Permission) BinaryUnmarshaler(data []byte) error {
	p.encoded = binary.BigEndian.Uint32(data)
	return nil
}

func resourceMask(resource StripeResource, access Access) uint32 {
	return uint32(access << (uint32(resource) * 2))
}

func (p *Permission) Can(access Access, resource StripeResource) bool {
	mask := resourceMask(resource, access)
	allMask := resourceMask(ResourceAll, access)
	return mask&p.encoded == mask || allMask&p.encoded == allMask
}

func (p *Permission) SetAccess(access Access, resources ...StripeResource) {
	for _, resource := range resources {
		p.encoded |= resourceMask(resource, access)
	}
}
