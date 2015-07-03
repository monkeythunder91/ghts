/* This file is part of GHTS.

GHTS is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

GHTS is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with GHTS.  If not, see <http://www.gnu.org/licenses/>.

@author: UnHa Kim <unha.kim.ghts@gmail.com> */

package shared

import (
	"strings"
	"testing"
)

func TestT주소String(테스트 *testing.T) {
	F테스트_참임(테스트, strings.HasPrefix(P주소_주소정보.String(), "tcp://127.0.0.1:"))
}

func TestT통화단위String(테스트 *testing.T) {
	F테스트_같음(테스트, KRW.String(), "KRW")
}