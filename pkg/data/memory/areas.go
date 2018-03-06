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

package memory

import (
	"log"
	"strconv"
	"strings"

	"github.com/mbertschler/bunny/pkg/data/stored"
	"github.com/tidwall/buntdb"
)

type areasTx struct {
	parent *Tx
	tx     *buntdb.Tx
}

func (t *areasTx) Key(id int) string {
	return areaPrefix + strconv.Itoa(id)
}

func (t *areasTx) ID(key string) int {
	key = strings.TrimPrefix(key, areaPrefix)
	i, err := strconv.Atoi(key)
	if err != nil {
		log.Println("KEY ERROR:", err)
	}
	return i
}

func (t *areasTx) Get(id int) (stored.Area, error) {
	var a stored.Area
	val, err := t.tx.Get(t.Key(id))
	if err != nil {
		return a, err
	}
	err = decode(val, &a)
	return a, err
}

func (t *areasTx) Set(a stored.Area) error {
	val, err := encode(a)
	if err != nil {
		return err
	}
	_, _, err = t.tx.Set(t.Key(a.ID), val, nil)
	return err
}
