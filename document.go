//  Copyright (c) 2015 Couchbase, Inc.
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

package index

import (
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"github.com/blevesearch/bleve/v2/size"
)

// A synonym document is a json object with the following fields:
//  1. mapping type: either "equivalent" or "explicit"
//     a. equivalent: all the phrases in the synonym list are equivalent to each other
//     b. explicit: each phrase in the input list is equivalent to the phrases in the synonym list,
//     but not to each other
//  2. input: a list of phrases
//  3. synonyms: a list of phrases
//
// A phrase is a sequence of words separated by spaces, and a word is a sequence of characters.
// A phrase can be a single word.
type SynonymDefinition struct {
	MappingType json.RawMessage   `json:"mappingType"`
	Input       []json.RawMessage `json:"input"`
	Synonyms    []json.RawMessage `json:"synonyms"`
}

func (s *SynonymDefinition) Size() int {
	var sd SynonymDefinition
	sizeInBytes := len(s.MappingType) + int(reflect.TypeOf(sd).Size()) + size.SizeOfPtr
	for _, entry := range s.Input {
		sizeInBytes += len(entry)
	}
	for _, entry := range s.Synonyms {
		sizeInBytes += len(entry)
	}
	return sizeInBytes
}

type Document interface {
	ID() string
	Size() int

	VisitFields(visitor FieldVisitor)
	VisitComposite(visitor CompositeFieldVisitor)
	HasComposite() bool

	NumPlainTextBytes() uint64

	AddIDField()

	StoredFieldsBytes() uint64
	SynonymInfo() *SynonymDefinition
}

type FieldVisitor func(Field)

type Field interface {
	Name() string
	Value() []byte
	ArrayPositions() []uint64

	EncodedFieldType() byte

	Analyze()

	AnalyzeSynonyms([]*SynonymDefinition, *sync.Map)

	Options() FieldIndexingOptions

	AnalyzedLength() int
	AnalyzedTokenFrequencies() TokenFrequencies

	NumPlainTextBytes() uint64
}

type CompositeFieldVisitor func(field CompositeField)

type CompositeField interface {
	Field

	Compose(field string, length int, freq TokenFrequencies)
}

type TextField interface {
	Text() string
}

type NumericField interface {
	Number() (float64, error)
}

type DateTimeField interface {
	DateTime() (time.Time, error)
}

type BooleanField interface {
	Boolean() (bool, error)
}

type GeoPointField interface {
	Lon() (float64, error)
	Lat() (float64, error)
}

type GeoShapeField interface {
	GeoShape() (GeoJSON, error)
}

// TokenizableSpatialField is an optional interface for fields that
// supports pluggable custom hierarchial spatial token generation.
type TokenizableSpatialField interface {
	// SetSpatialAnalyzerPlugin lets the index implementations to
	// initialise relevant spatial analyzer plugins for the field
	// to override the spatial token generations during the analysis phase.
	SetSpatialAnalyzerPlugin(SpatialAnalyzerPlugin)
}
