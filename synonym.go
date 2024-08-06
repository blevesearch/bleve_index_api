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

package index

import (
	"fmt"
	"reflect"
	"strconv"
)

var EquivalentSynonymType = "equivalent"
var ExplicitSynonymType = "explicit"
var SynonymKeySeparator = []byte{'\xff'}

var reflectStaticSizeSynonymDefinition int

func init() {
	var sd SynonymDefinition
	reflectStaticSizeSynonymDefinition = int(reflect.TypeOf(sd).Size())
}

type SynonymDefinition struct {
	MappingType string   `json:"mappingType"`
	Input       []string `json:"input,omitempty"`
	Synonyms    []string `json:"synonyms"`
}

func (s *SynonymDefinition) Size() int {
	sizeInBytes := reflectStaticSizeSynonymDefinition +
		len(s.MappingType)
	for _, entry := range s.Input {
		sizeInBytes += len(entry)
	}
	for _, entry := range s.Synonyms {
		sizeInBytes += len(entry)
	}
	return sizeInBytes
}

func (s *SynonymDefinition) Validate() error {
	if s.MappingType != EquivalentSynonymType &&
		s.MappingType != ExplicitSynonymType {
		return fmt.Errorf("invalid mappingType; must be either `%s` or `%s`",
			EquivalentSynonymType, ExplicitSynonymType)
	}
	if len(s.Synonyms) == 0 {
		return fmt.Errorf("`synonyms` field missing or empty")
	}
	return nil
}

type SynonymMetadata struct {
	HashToSynonyms map[uint64][]*SynonymDocumentMap `json:"hashToSynonyms"`
	HashToPhrase   map[uint64]*PhraseDocumentMap    `json:"hashToPhrase"`
	SynonymFST     []byte                           `json:"synonymFST"`
}

type PhraseDocumentMap struct {
	// the actual phrase
	Phrase string

	// set of the doc numbers that contain this phrase on the LHS
	DocNums map[uint32]struct{}

	// value that indicates whether this LHS value is invalid because
	// all docs in docnums have been deleted
	IsInvalid bool
}

type SynonymDocumentMap struct {
	// the hash of the phrase
	Hash uint64

	// set of the doc numbers that contain this phrase on the RHS
	DocNums map[uint32]struct{}

	// value that indicates whether this RHS value is invalid because
	// all docs in docnums have been deleted
	IsInvalid bool
}

func CreateSynonymMetadataKey(collection string, analyzerName string) string {
	return collection + string(SynonymKeySeparator) + analyzerName
}

func (s *SynonymMetadata) String() string {
	if s.HashToSynonyms == nil {
		return "Empty Metadata"
	}
	rv := ""
	rv += "LHS, [DocNums], Invalid / Valid\t-\t[[RHS, DocNums, Invalid / Valid]]\n\n"
	for hash, synonyms := range s.HashToSynonyms {
		phraseDocMap := s.HashToPhrase[hash]
		if phraseDocMap == nil {
			return "Error: LHS Struct is nil"
		}
		LHS := phraseDocMap.Phrase
		LHSDocNums := "["
		idx := 0
		for docNum, _ := range phraseDocMap.DocNums {
			if idx == len(phraseDocMap.DocNums)-1 {
				LHSDocNums += strconv.FormatUint(uint64(docNum), 10) + "]"
			} else {
				LHSDocNums += strconv.FormatUint(uint64(docNum), 10) + ", "
			}
			idx++
		}
		lhsInvalidString := "Valid"
		if phraseDocMap.IsInvalid {
			lhsInvalidString = "Invalid"
		}

		rv += string(LHS) + ", " + LHSDocNums + ", " + lhsInvalidString + "\t-\t["
		for oidx, synonym := range synonyms {
			rhsStruct := s.HashToPhrase[synonym.Hash]
			if rhsStruct == nil {
				return "Error: RHS Struct is nil"
			}
			RHS := rhsStruct.Phrase
			rv += string(RHS) + ", "
			RHSDocNums := "["
			idx := 0
			for docNum, _ := range synonym.DocNums {
				if idx == len(synonym.DocNums)-1 {
					RHSDocNums += strconv.FormatUint(uint64(docNum), 10) + "]"
				} else {
					RHSDocNums += strconv.FormatUint(uint64(docNum), 10) + ", "
				}
				idx++
			}

			rhsInvalidString := "Valid"
			if synonym.IsInvalid {
				rhsInvalidString = "Invalid"
			}

			rv += RHSDocNums + ", " + rhsInvalidString
			if oidx == len(synonyms)-1 {
				rv += "]\n"
			} else {
				rv += "], ["
			}
		}
	}
	return rv
}
