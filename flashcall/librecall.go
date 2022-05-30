package flashcall

import (
	"errors"
	"time"
)

var errBadCodeLength = errors.New("code length should be 1..4")
var errBadCodeFormat = errors.New("code char should be 0..9")
var errBadPhoneFormat = errors.New("bad phone format")

type code struct {
	code string
}

type Phone struct {
	phone string
}

type CallRecord struct {
	CallAt time.Time
	To     string
	From   string
	Code   string
}

type LibreCall struct {
	stor *MemoryStorage
}

func newCode(confirmCode string) (code, error) {
	if len(confirmCode) < 1 || len(confirmCode) > 4 {
		return code{}, errBadCodeLength
	}

	for _, r := range confirmCode {
		if r < '0' || r > '9' {
			return code{}, errBadCodeFormat
		}
	}

	return code{code: confirmCode}, nil
}

func newPhone(strPhone string) (Phone, error) {
	if len(strPhone) < 10 || len(strPhone) > 11 {
		return Phone{}, errBadPhoneFormat
	}

	for i, r := range strPhone {
		if (r < '0' || r > '9') && r != '+' && i != 0 {
			return Phone{}, errBadCodeFormat
		}
	}

	return Phone{phone: strPhone}, nil
}

func (p *Phone) replaceSuffixWithCode(suffixCode code) {
	p.phone = p.phone[:len(p.phone)-len(suffixCode.code)] + suffixCode.code
}

func NewLibrecall(storage *MemoryStorage) *LibreCall {
	return &LibreCall{
		stor: storage,
	}
}

func (call *LibreCall) Call(callTo Phone, confirmCode code) {
	phoneFrom := generateRandomPhoneWithSuffixCode(confirmCode)
	call.stor.addRecord(callTo.phone, phoneFrom.phone, confirmCode.code)
}

func (call *LibreCall) AllCallsSortedDesc() []CallRecord {
	all := call.stor.allRecordsReversed()

	result := make([]CallRecord, 0, len(all))
	for _, v := range all {
		result = append(result, CallRecord{
			CallAt: v.at,
			To:     v.to,
			From:   v.from,
			Code:   v.code,
		})
	}

	return result
}

func generateRandomPhoneWithSuffixCode(code code) Phone {
	p, err := newPhone("89441231230")
	if err != nil {
		panic(err)
	}

	p.replaceSuffixWithCode(code)

	return p
}
