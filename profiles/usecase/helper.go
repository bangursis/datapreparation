package usecase

import (
	"datapreparation/pkg/cryptohelper"
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	BadFileErr = errors.New("File isn't valid")
)

type header struct {
	M1    int    `asn1:"optional,tag:0"`
	M2    int    `asn1:"optional,tag:1"`
	Pt    string `asn1:"optional,tag:2"`
	ICCID []byte `asn1:"tag:3"`
}

func (uc *profilesUseCase) readFromFile(filepath string) ([]byte, error) {

	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	pe, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return pe, err
}

func (uc *profilesUseCase) getICCID(pe []byte) ([]byte, error) {
	cppe := make([]byte, len(pe))
	copy(cppe, pe)
	cppe[0] = '3'

	peHeader := make([]byte, hex.DecodedLen(len(cppe)))
	_, err := hex.Decode(peHeader, cppe)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", BadFileErr, err)
	}

	if len(peHeader) < 2 {
		return nil, BadFileErr
	}
	headerLength := peHeader[1]

	if len(peHeader) < int(2+headerLength) {
		return nil, BadFileErr
	}

	peHeader = peHeader[:2+headerLength]
	h := header{}
	_, err = asn1.Unmarshal(peHeader, &h)
	if err != nil {
		return nil, err
	}

	iccid := make([]byte, hex.EncodedLen(len(h.ICCID)))
	hex.Encode(iccid, h.ICCID)

	return iccid, nil
}

func (uc *profilesUseCase) generateKeys() ([]byte, []byte, []byte) {
	return cryptohelper.HexGenerator(16), cryptohelper.HexGenerator(16), cryptohelper.HexGenerator(16)
}

func (uc *profilesUseCase) splitToChunks(iccID []byte, size int) [][]byte {
	out := make([][]byte, 0, (len(iccID)/size)+1)

	for len(iccID) > 0 {
		out = append(out, iccID[:size])
		iccID = iccID[size:]

		if l := len(iccID); l < size && l%size != 0 {
			size = len(iccID)
		}
	}

	return out
}
