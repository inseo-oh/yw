package csssyntax

import (
	"fmt"

	"github.com/inseo-oh/yw/css/float"
)

// https://www.w3.org/TR/CSS2/visuren.html#propdef-float
func (ts *tokenStream) parseFloat() (res float.Float, err error) {
	if err := ts.consumeIdentTokenWith("none"); err == nil {
		return float.None, nil
	} else if err := ts.consumeIdentTokenWith("left"); err == nil {
		return float.Left, nil
	} else if err := ts.consumeIdentTokenWith("right"); err == nil {
		return float.Right, nil
	}
	return res, fmt.Errorf("%s: invalid float value", ts.errorHeader())
}
