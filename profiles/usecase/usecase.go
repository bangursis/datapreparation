package usecase

import (
	"context"
	"datapreparation/pkg/cryptohelper"
	"datapreparation/profiles"
	"errors"
	"fmt"
	"regexp"

	"github.com/bangursis/crypto"
)

var (
	DecryptionErr             = errors.New("Failed to decrypt request")
	GetICCIDErr               = errors.New("Failed to get ICCID")
	ValidateIDOverFileNameErr = errors.New("Filename does not contain the ICCID stored")
)

type profilesUseCase struct {
	Decrypt cryptohelper.Decrypt
	Repo    profiles.Repository
}

func Init(repo profiles.Repository, dec cryptohelper.Decrypt) profiles.UseCase {
	return &profilesUseCase{
		Repo:    repo,
		Decrypt: dec,
	}
}

func (uc *profilesUseCase) Import(ctx context.Context, filepath string) ([]byte, error) {
	filepath, err := uc.Decrypt(filepath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", DecryptionErr, err)
	}

	pe, err := uc.readFromFile(filepath)
	if err != nil {
		return nil, err
	}

	iccID, err := uc.getICCID(pe)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", GetICCIDErr, err)
	}
	readableICCID := string(iccID)

	pattern := fmt.Sprintf(`\S*\/?%s\.?[a-zA-z0-9]*$`, readableICCID)
	if matched, _ := regexp.MatchString(pattern, filepath); !matched {
		return nil, ValidateIDOverFileNameErr
	}

	encKey, macKey, macChaining := uc.generateKeys()
	chunks := uc.splitToChunks(pe, 1007)
	encrypteds := make([][]byte, 0, len(chunks))
	load := make([][]byte, 0, len(chunks))

	for i, c := range chunks {
		enc := crypto.EncryptScp11(byte(i+1), c, encKey)
		singed, macChaining := crypto.SignScp11(enc, 0x86, macKey, macChaining)

		encrypteds = append(encrypteds, enc)
		load = append(load, [][]byte{enc, macChaining, singed}...)
	}

	if err := uc.Repo.Save(ctx, readableICCID, encrypteds); err != nil {
		return nil, err
	}

	return flattenSlice(load), err
}
