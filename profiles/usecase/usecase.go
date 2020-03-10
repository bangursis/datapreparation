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

func (uc *profilesUseCase) Import(ctx context.Context, filepath string) error {
	filepath, err := uc.Decrypt(filepath)
	if err != nil {
		return fmt.Errorf("%w: %v", DecryptionErr, err)
	}

	pe, err := uc.readFromFile(filepath)
	if err != nil {
		return err
	}

	iccID, err := uc.getICCID(pe)
	if err != nil {
		return fmt.Errorf("%w: %v", GetICCIDErr, err)
	}
	readableICCID := string(iccID)

	pattern := fmt.Sprintf(`\S*\/?%s\.?[a-zA-z0-9]*$`, readableICCID)
	if matched, _ := regexp.MatchString(pattern, filepath); !matched {
		return ValidateIDOverFileNameErr
	}

	encKey, macKey, macChaining := uc.generateKeys()
	chunks := uc.splitToChunks(iccID, 1007)
	encrypted := [][]byte{}

	for i, c := range chunks {
		enc := crypto.EncryptScp11(byte(i+1), c, encKey)
		_, macChaining = crypto.SignScp11(enc, 0x86, macKey, macChaining)

		encrypted = append(encrypted, enc)
	}

	if err := uc.Repo.Save(ctx, readableICCID, encrypted); err != nil {
		return err
	}

	return nil
}
