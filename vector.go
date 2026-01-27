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

//go:build vectors
// +build vectors

package index

type VectorField interface {
	Vector() []float32
	// Dimensionality of the vector
	Dims() int
	// Similarity metric to be used for scoring the vectors
	Similarity() string
	// nlist/nprobe config (recall/latency) the index is optimized for
	IndexOptimizedFor() string
}

// -----------------------------------------------------------------------------

const (
	EuclideanDistance = "l2_norm"

	InnerProduct = "dot_product"

	CosineSimilarity = "cosine"
)

const DefaultVectorSimilarityMetric = EuclideanDistance

// Supported similarity metrics for vector fields
var SupportedVectorSimilarityMetrics = map[string]struct{}{
	EuclideanDistance: {},
	InnerProduct:      {},
	CosineSimilarity:  {},
}

// -----------------------------------------------------------------------------

// Types of vector indexes
const (
	FloatVectorIndex = "float32"
	BinaryVectorIndex = "binary"
)

var SupportedVectorIndexTypes = map[string]int{
	FloatVectorIndex:    0,
	BinaryVectorIndex:   1,
}

var SupportedVectorIndexTypesReverse = map[int]string{
	0: FloatVectorIndex,
	1: BinaryVectorIndex,
}

// -----------------------------------------------------------------------------

const (
	IndexOptimizedForRecall          = "recall"
	IndexOptimizedForLatency         = "latency"
	IndexOptimizedForMemoryEfficient = "memory-efficient"
	IndexOptimizedForRecallBinary   = "recall,binary"
	IndexOptimizedForLatencyBinary  = "latency,binary"
)

const DefaultIndexOptimization = IndexOptimizedForRecall

var SupportedVectorIndexOptimizations = map[string]int{
	IndexOptimizedForRecall:          0,
	IndexOptimizedForLatency:         1,
	IndexOptimizedForMemoryEfficient: 2,
	IndexOptimizedForRecallBinary:    3,
	IndexOptimizedForLatencyBinary:   4,
}

// Reverse maps vector index optimizations': int -> string
var VectorIndexOptimizationsReverseLookup = map[int]string{
	0: IndexOptimizedForRecall,
	1: IndexOptimizedForLatency,
	2: IndexOptimizedForMemoryEfficient,
	3: IndexOptimizedForRecallBinary,
	4: IndexOptimizedForLatencyBinary,
}
