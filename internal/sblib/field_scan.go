package sblib

import (
	"encoding/hex"
	"math/rand"
	"os"
	"path"
	"strconv"

	"github.com/pkg/errors"
)

type SnowFieldInfo struct {
	Name           string
	Length         int64
	MerkleRootHash []byte
	// ActivationTarget
}

func mustDecodeHex(byteString string) []byte {
	bytes, err := hex.DecodeString(byteString)
	if err != nil {
		panic(err)
	}
	return bytes
}

var SnowFields = map[int]SnowFieldInfo{
	0:  {"cricket", 1 * 1024 * 1024 * 1024, mustDecodeHex("2c363d33550f5c4da16279d1fa89f7a9")},
	1:  {"shrew", 2 * 1024 * 1024 * 1024, mustDecodeHex("4626ba6c78e0d777d35ae2044f8882fa")},
	2:  {"stoat", 4 * 1024 * 1024 * 1024, mustDecodeHex("2948b321266dfdec11fbdc45f46cf959")},
	3:  {"ocelot", 8 * 1024 * 1024 * 1024, mustDecodeHex("33a2fb16f08347c2dc4cbcecb4f97aeb")},
	4:  {"pudu", 16 * 1024 * 1024 * 1024, mustDecodeHex("0cdcb8629bef77fbe1f9b740ec3897c9")},
	5:  {"badger", 32 * 1024 * 1024 * 1024, mustDecodeHex("e9abff32aa7f74795be2aa539a079489")},
	6:  {"capybara", 64 * 1024 * 1024 * 1024, mustDecodeHex("e68678fadb1750feedfa11522270497f")},
	7:  {"llama", 128 * 1024 * 1024 * 1024, mustDecodeHex("147d379701f621ebe53f5a511bd6c380")},
	8:  {"bugbear", 256 * 1024 * 1024 * 1024, mustDecodeHex("533ae42a04a609c4b45464b1aa9e6924")},
	9:  {"hippo", 512 * 1024 * 1024 * 1024, mustDecodeHex("ecdb28f1912e2266a71e71de601117d2")},
	10: {"shai-hulud", 1024 * 1024 * 1024 * 1024, mustDecodeHex("cc883468a08f48b592a342a2cdf5bcba")},
	11: {"avanc", 2048 * 1024 * 1024 * 1024, mustDecodeHex("a2a4076f6cde947935db06e5fc5bbd14")},
}

var NetworkName = "snowblossom"

func LoadFields(dirPath string) (map[int]SnowMerkleProof, error) {
	fields := map[int]SnowMerkleProof{}

	for fieldNumber := range SnowFields {
		name := NetworkName + "." + strconv.Itoa(fieldNumber)
		fieldFolder := path.Join(dirPath, name)
		fieldFolderStat, err := os.Stat(fieldFolder)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, errors.Wrapf(err, "failed to stat field folder %v", fieldFolder)
		}

		if !fieldFolderStat.IsDir() {
			continue
		}

		field, err := LoadField(fieldFolder, name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to load field %v", name)
		}

		for i := 0; i < 64; i++ {
			if err := checkField(fieldNumber, field); err != nil {
				return nil, errors.Wrapf(err, "checkField failed for field %v on round %v", name, i)
			}
		}

		fields[fieldNumber] = field
	}

	return fields, nil
}

func checkField(fieldNumber int, field SnowMerkleProof) error {
	max := field.TotalWords()

	check := rand.Int63n(max)
	proof, err := field.GetProof(check)
	if err != nil {
		return errors.Wrapf(err, "checkField -> GetProof failed for field #%v", fieldNumber)
	}

	err = CheckProof(proof, SnowFields[fieldNumber].MerkleRootHash, SnowFields[fieldNumber].Length)
	if err != nil {
		return errors.Wrapf(err, "checkField -> CheckProof failed for field #%v", fieldNumber)
	}

	return nil
}

func LoadField(dirPath string, base string) (SnowMerkleProof, error) {
	return NewSnowMerkleProof(dirPath, base)
}
