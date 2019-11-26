/*
 * Copyright 2019, Offchain Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package protocol

import (
	"fmt"

	"github.com/offchainlabs/arbitrum/packages/arb-util/value"

	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

type AssertionStub struct {
	AfterHash        [32]byte
	NumSteps         uint32
	FirstMessageHash [32]byte
	LastMessageHash  [32]byte
	FirstLogHash     [32]byte
	LastLogHash      [32]byte
}

func (a *AssertionStub) String() string {
	return fmt.Sprintf("AssertionStub(%x, %v, %x, %x)", a.AfterHash, a.NumSteps, a.FirstMessageHash, a.LastMessageHash)
}

func (a *AssertionStub) Equals(b *AssertionStub) bool {
	if a.AfterHash != b.AfterHash ||
		a.NumSteps != b.NumSteps ||
		a.FirstMessageHash != b.FirstMessageHash ||
		a.LastMessageHash != b.LastMessageHash {
		return false
	}
	return true
}

func (a *AssertionStub) Hash() [32]byte {
	var ret [32]byte
	hashVal := solsha3.SoliditySHA3(
		solsha3.Bytes32(a.AfterHash),
		solsha3.Uint32(a.NumSteps),
		solsha3.Bytes32(a.FirstMessageHash),
		solsha3.Bytes32(a.LastMessageHash),
		solsha3.Bytes32(a.FirstLogHash),
		solsha3.Bytes32(a.LastLogHash),
	)
	copy(ret[:], hashVal)
	return ret
}

func (a *AssertionStub) GeneratePostcondition(pre *Precondition) *Precondition {
	return &Precondition{
		BeforeHash:  value.NewHashBuf(a.AfterHash),
		TimeBounds:  pre.TimeBounds,
		BeforeInbox: pre.BeforeInbox,
	}
}

func GeneratePreconditions(pre *Precondition, assertions []*AssertionStub) []*Precondition {
	preconditions := make([]*Precondition, 0, len(assertions))
	for _, assertion := range assertions {
		preconditions = append(preconditions, pre)
		pre = assertion.GeneratePostcondition(pre)
	}
	return preconditions
}
