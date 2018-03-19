// Copyright 2018 Martin Bertschler.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"os"
	"testing"
)

func TestEnvOrFallback(t *testing.T) {
	err := os.Setenv("BUNNY_TEST", "value")
	if err != nil {
		t.Error(err)
	}
	str := envOrFallback("BUNNY_TEST", "fall")
	if str != "value" {
		t.Error("expected", str, "to be \"value\"")
	}
	str = envOrFallback("BUNNY_TEST2", "fall")
	if str != "fall" {
		t.Error("expected", str, "to be \"fall\"")
	}
}

func TestFindProjectFolder(t *testing.T) {
	_, err := findProjectFolder()
	if err != nil {
		t.Error(err)
	}
	gopath := os.Getenv("GOPATH")
	err = os.Unsetenv("GOPATH")
	if err != nil {
		t.Error(err)
	}

	// should fail without GOPATH
	_, err = findProjectFolder()
	if err == nil {
		t.Error("expected an error")
	}
	err = os.Setenv("GOPATH", gopath)
	if err != nil {
		t.Error(err)
	}
}

func TestSetup(t *testing.T) {
	// test env config
	err := os.Setenv("BUNNY_PORT", "1234")
	if err != nil {
		t.Error(err)
	}
	err = os.Setenv("BUNNY_ROOT", "/a/b/c")
	if err != nil {
		t.Error(err)
	}
	err = Setup()
	if err != nil {
		t.Error(err)
	}
	if Port != "1234" {
		t.Error("expected", Port, "to be \"1234\"")
	}
	if Root != "/a/b/c" {
		t.Error("expected", Root, "to be \"/a/b/c\"")
	}

	// test defaults
	err = os.Unsetenv("BUNNY_PORT")
	if err != nil {
		t.Error(err)
	}
	err = os.Unsetenv("BUNNY_ROOT")
	if err != nil {
		t.Error(err)
	}
	err = Setup()
	if err != nil {
		t.Error(err)
	}
	if Port != "3080" {
		t.Error("expected", Port, "to be \"1234\"")
	}
	folder, err := findProjectFolder()
	if err != nil {
		t.Error(err)
	}
	if Root != folder {
		t.Error("expected", Root, "to be", folder)
	}
}
