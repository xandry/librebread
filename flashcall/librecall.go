package flashcall

import (
	"errors"
)

var ErrBadCodeLength = errors.New("code length should be 1..4")
var ErrBadCodeFormat = errors.New("code char should be 0..9")

type code struct {
	code string
}

type Librecall struct {
	stor *MemoryStorage
}

func newCode(confirmCode string) (code, error) {
	if len(confirmCode) < 1 || len(confirmCode) > 4 {
		return code{}, ErrBadCodeLength
	}

	for _, r := range confirmCode {
		if r < '0' || r > '9' {
			return code{}, ErrBadCodeFormat
		}
	}

	return code{code: confirmCode}, nil
}

func (call *Librecall) Call(clientPhone, confirmCode code) error {

}

func generatePhoneWithCode(confirmCode code) string {

}
