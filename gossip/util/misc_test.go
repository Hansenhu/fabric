/*
Copyright IBM Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"crypto/rand"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func testHappyPath(t *testing.T) {
	n1 := RandomInt(10000)
	n2 := RandomInt(10000)
	assert.NotEqual(t, n1, n2)
	n3 := RandomUInt64()
	n4 := RandomUInt64()
	assert.NotEqual(t, n3, n4)
}

func TestGetRandomInt(t *testing.T) {
	testHappyPath(t)
}

func TestNonNegativeValues(t *testing.T) {
	assert.True(t, RandomInt(1000000) >= 0)
}

func TestGetRandomIntBadInput(t *testing.T) {
	f1 := func() {
		RandomInt(0)
	}
	f2 := func() {
		RandomInt(-500)
	}
	assert.Panics(t, f1)
	assert.Panics(t, f2)
}

type reader struct {
	mock.Mock
}

func (r *reader) Read(p []byte) (int, error) {
	args := r.Mock.Called(p)
	n := args.Get(0).(int)
	err := args.Get(1)
	if err == nil {
		return n, nil
	}
	return n, err.(error)
}

func TestGetRandomIntNoEntropy(t *testing.T) {
	rr := rand.Reader
	defer func() {
		rand.Reader = rr
	}()
	r := &reader{}
	r.On("Read", mock.Anything).Return(0, errors.New("Not enough entropy"))
	rand.Reader = r
	// Make sure randomness still works even when we have no entropy
	testHappyPath(t)
}
