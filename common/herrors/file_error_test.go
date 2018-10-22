// Copyright 2018 The Hugo Authors. All rights reserved.
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

package herrors

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToLineNumberError(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	for i, test := range []struct {
		in           error
		offset       int
		lineNumber   int
		columnNumber int
	}{
		{errors.New("no line number for you"), 0, -1, 1},
		{errors.New(`template: _default/single.html:4:15: executing "_default/single.html" at <.Titles>: can't evaluate field Titles in type *hugolib.PageOutput`), 0, 4, 15},
		{errors.New("parse failed: template: _default/bundle-resource-meta.html:11: unexpected in operand"), 0, 11, 1},
		{errors.New(`failed:: template: _default/bundle-resource-meta.html:2:7: executing "main" at <.Titles>`), 0, 2, 7},
		{errors.New("error in front matter: Near line 32 (last key parsed 'title')"), 0, 32, 1},
	} {

		got := ToFileErrorWithOffset("template", test.in, test.offset)

		errMsg := fmt.Sprintf("[%d][%T]", i, got)
		le, ok := got.(FileError)

		if test.lineNumber > 0 {
			assert.True(ok)
			assert.Equal(test.lineNumber, le.LineNumber(), errMsg)
			assert.Equal(test.columnNumber, le.ColumnNumber(), errMsg)
			assert.Contains(got.Error(), strconv.Itoa(le.LineNumber()))
		} else {
			assert.False(ok)
		}
	}
}