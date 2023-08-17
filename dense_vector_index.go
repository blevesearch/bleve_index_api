//  Copyright (c) 2023 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build densevector
// +build densevector

package index

import (
	"context"
)

type VectorReader interface {
	// Next returns the next document similar to the vector, in this field, or nil
	// when it reaches the end of the enumeration.  The preAlloced VectorDoc
	// is optional, and when non-nil, will be used instead of allocating memory.
	Next(preAlloced *VectorDoc) (*VectorDoc, error)

	// Advance resets the enumeration at specified document or its immediate
	// follower.
	Advance(ID IndexInternalID, preAlloced *VectorDoc) (*VectorDoc, error)

	// Count returns the number of documents similar to the vector, in this field.
	Count() uint64
	Close() error

	Size() int
}

type VectorIndexReader interface {
	VectorReader(ctx context.Context, vector []float32, field string, k int64) (
		VectorReader, error)
}

type VectorDoc struct {
	Vector []float32
	ID     IndexInternalID
	Score  float64
}

func (vd *VectorDoc) Size() int {
	return 1
}

// Reset allows an already allocated VectorDoc to be reused
func (vd *VectorDoc) Reset() *VectorDoc {
	// remember the []byte used for the ID
	id := vd.ID
	// idiom to copy over from empty VectorDoc (0 allocations)
	*vd = VectorDoc{}
	// reuse the []byte already allocated (and reset len to 0)
	vd.ID = id[:0]
	return vd
}
