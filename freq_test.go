//  Copyright (c) 2014 Couchbase, Inc.
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
	"reflect"
	"testing"
)

func TestTokenFrequenciesMergeAll(t *testing.T) {
	tf1 := TokenFrequencies{
		"water": &TokenFreq{
			Term: []byte("water"),
			Locations: []*TokenLocation{
				{
					Position: 1,
					Start:    0,
					End:      5,
				},
				{
					Position: 2,
					Start:    6,
					End:      11,
				},
			},
		},
	}
	tf2 := TokenFrequencies{
		"water": &TokenFreq{
			Term: []byte("water"),
			Locations: []*TokenLocation{
				{
					Position: 1,
					Start:    0,
					End:      5,
				},
				{
					Position: 2,
					Start:    6,
					End:      11,
				},
			},
		},
	}
	expectedResult := TokenFrequencies{
		"water": &TokenFreq{
			Term: []byte("water"),
			Locations: []*TokenLocation{
				{
					Position: 1,
					Start:    0,
					End:      5,
				},
				{
					Position: 2,
					Start:    6,
					End:      11,
				},
				{
					Field:    "tf2",
					Position: 1,
					Start:    0,
					End:      5,
				},
				{
					Field:    "tf2",
					Position: 2,
					Start:    6,
					End:      11,
				},
			},
		},
	}
	tf1.MergeAll("tf2", tf2)
	if !reflect.DeepEqual(tf1, expectedResult) {
		t.Errorf("expected %#v, got %#v", expectedResult, tf1)
	}
}

func TestTokenFrequenciesMergeAllLeftEmpty(t *testing.T) {
	tf1 := TokenFrequencies{}
	tf2 := TokenFrequencies{
		"water": &TokenFreq{
			Term: []byte("water"),
			Locations: []*TokenLocation{
				{
					Position: 1,
					Start:    0,
					End:      5,
				},
				{
					Position: 2,
					Start:    6,
					End:      11,
				},
			},
		},
	}
	expectedResult := TokenFrequencies{
		"water": &TokenFreq{
			Term: []byte("water"),
			Locations: []*TokenLocation{
				{
					Field:    "tf2",
					Position: 1,
					Start:    0,
					End:      5,
				},
				{
					Field:    "tf2",
					Position: 2,
					Start:    6,
					End:      11,
				},
			},
		},
	}
	tf1.MergeAll("tf2", tf2)
	if !reflect.DeepEqual(tf1, expectedResult) {
		t.Errorf("expected %#v, got %#v", expectedResult, tf1)
	}
}

func TestTermFieldDocResetClearsNormByte(t *testing.T) {
	tfd := &TermFieldDoc{
		Term:     "beer",
		Freq:     7,
		Norm:     0.5,
		NormByte: 42,
		ID:       IndexInternalID("abc"),
	}
	tfd.Reset()
	if tfd.NormByte != 0 {
		t.Errorf("expected NormByte=0 after Reset, got %d", tfd.NormByte)
	}
	if tfd.Freq != 0 {
		t.Errorf("expected Freq=0 after Reset, got %d", tfd.Freq)
	}
	if tfd.Term != "" {
		t.Errorf("expected Term=\"\" after Reset, got %q", tfd.Term)
	}
	// ID and Vectors backing arrays should be preserved (zero-length, not nil)
	if tfd.ID == nil {
		t.Error("expected ID backing array preserved after Reset")
	}
}
