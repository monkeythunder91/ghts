/*
This file is part of GHTS.

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

Created on 2015. 4. 5.

@author: UnHa Kim <unha.kim@gh-system.com>
*/

package shared

var p파이썬_경로 string = ""

func init() {
	p파이썬_경로 = F실행화일_검색("python.exe")
}
