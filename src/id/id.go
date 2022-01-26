/*
    This file is part of manifold.
    Manifold is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.
    Manifold is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.
    You should have received a copy of the GNU General Public License
    along with manifold.  If not, see <https://www.gnu.org/licenses/>.
*/
package id

import (
	"math"
	"math/big"
	"crypto/rand"
)

type Id = uint64

type IdSlice []Id

func (list IdSlice) Has(a Id) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func NewID() (Id, error) {
    val, err := rand.Int(rand.Reader, big.NewInt(int64(math.MaxInt64)))
    if err != nil {
        return 0, err
    }
    return val.Uint64(), err
}
