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

import "math"

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

const DefaultSimilarityMetric = EuclideanDistance

// Supported similarity metrics for vector fields
var SupportedSimilarityMetrics = map[string]struct{}{
	EuclideanDistance: {},
	InnerProduct:      {},
	CosineSimilarity:  {},
}

func NormalizeVector(vector []float32) []float32 {
	// first calculate the magnitude of the vector
	var mag float32
	for _, v := range vector {
		mag += v * v
	}
	mag = float32(math.Sqrt(float64(mag)))
	// cannot normalize a zero vector
	// if the magnitude is 1, then the vector is already normalized
	if mag != 0 && mag != 1 {
		// normalize the vector
		for i, v := range vector {
			vector[i] = v / mag
		}
	}
	return vector
}

// -----------------------------------------------------------------------------

const (
	IndexOptimizedForRecall          = "recall"
	IndexOptimizedForLatency         = "latency"
	IndexOptimizedForMemoryEfficient = "memory-efficient"
)

const DefaultIndexOptimization = IndexOptimizedForRecall

var SupportedVectorIndexOptimizations = map[string]int{
	IndexOptimizedForRecall:          0,
	IndexOptimizedForLatency:         1,
	IndexOptimizedForMemoryEfficient: 2,
}

// Reverse maps vector index optimizations': int -> string
var VectorIndexOptimizationsReverseLookup = map[int]string{
	0: IndexOptimizedForRecall,
	1: IndexOptimizedForLatency,
	2: IndexOptimizedForMemoryEfficient,
}
