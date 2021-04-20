package internal

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
)

//Mkdir helper, creates directory if not exists
func Mkdir(path string, perm fs.FileMode) error {

	if _, errFolder := os.Stat(path); os.IsNotExist(errFolder) {
		if errMkDir := os.MkdirAll(path, perm); errMkDir != nil {
			return errMkDir
		}
	}
	return nil
}

//Hash generate sha256 of buffer
func Hash(buffer []byte) (string, error) {
	sha256 := sha256.New()
	if _, errSha256 := sha256.Write(buffer); errSha256 != nil {
		return "", errSha256
	}
	return hex.EncodeToString(sha256.Sum(nil)), nil
}

//Save Will save the state of the forest into file
func Save(object interface{}, path string, compress bool) (string, error) {

	var resultFile *os.File
	var errFile error
	var fileName string
	var hash string
	var errHash error

	//Create results folder
	if errMkdir := Mkdir(path, os.ModeDir); errMkdir != nil {
		return "", errMkdir
	}
	//forest.ToBytes(compress)
	buffer, errBUF := ToBytes(object, compress)
	if errBUF != nil {
		return "", errBUF
	}
	//Hash Data
	if hash, errHash = Hash(buffer); errHash != nil {
		return "", errHash
	}

	//filename
	fileName = path + "forest" + "-" + hash + ".bin"

	//Check if file already exists (just in case...)
	if _, errExits := os.Stat(fileName); !os.IsNotExist(errExits) {
		return "", errExits
	}
	//Create the File
	if resultFile, errFile = os.Create(fileName); errFile == nil {
		defer resultFile.Close()
	} else {
		return "", errFile
	}

	//Write into file
	if _, errWrite := resultFile.Write(buffer); errWrite != nil {
		return "", errWrite
	}

	return fileName, nil
}

//Load load forest from file
func Load(path string, forest interface{}) error {
	var forestFile *os.File = nil
	var errorForestFile error
	var buffer *bytes.Buffer = bytes.NewBuffer(nil)

	//If file does not exist
	if _, errExits := os.Stat(path); os.IsNotExist(errExits) {
		return errExits
	}

	//Open file
	if forestFile, errorForestFile = os.Open(path); errorForestFile != nil {
		return errorForestFile
	}

	//Copy file into buffer
	if _, errCopy := io.Copy(buffer, forestFile); errCopy != nil {
		return errCopy
	}

	//Check if file is compressed
	if ok, err := IsCompressed(buffer); err == nil && ok {
		if b, err := UnCompress(buffer); err == nil {

			return decode(b, forest)
		} else {
			return err
		}

	} else if err != nil {
		return err

	} else {

		return decode(buffer, forest)
	}
}

//ToBytes convert object to byte array, you can compress it with compress arg
func ToBytes(object interface{}, compress bool) ([]byte, error) {
	var buffer *bytes.Buffer = new(bytes.Buffer)

	//Encode Struct
	encoder := gob.NewEncoder(buffer)
	if errEncoding := encoder.Encode(object); errEncoding != nil {
		return nil, errEncoding
	}

	var binForest *bytes.Buffer = buffer
	if compress {
		var bufferCompressed *bytes.Buffer = new(bytes.Buffer)
		wCompress := zlib.NewWriter(bufferCompressed)
		if _, err := wCompress.Write(buffer.Bytes()); err != nil {
			return nil, err
		}
		wCompress.Close()
		binForest = bufferCompressed
	}

	return binForest.Bytes(), nil
}

//IsCompressed checks if buffer is compressed by zlib
func IsCompressed(buffer *bytes.Buffer) (bool, error) {
	if buffer == nil {
		return false, errors.New("buffer is nil")
	}
	scratch := buffer.Bytes()[:2]
	if scratch[0] == 0x78 && scratch[1] == 0x9c {
		return true, nil
	}
	return false, nil
}

//UnCompress unzips
func UnCompress(reader io.Reader) (*bytes.Buffer, error) {

	buffer := bytes.NewBuffer(nil)
	if _, errCopy := io.Copy(buffer, reader); errCopy != nil {
		return nil, errCopy
	}

	if ok, err := IsCompressed(buffer); err == nil && ok {

		if r, errZlib := zlib.NewReader(buffer); errZlib == nil {
			io.Copy(buffer, r)
			r.Close()

		} else {
			return nil, errZlib
		}

	} else if err != nil {
		return nil, err

	} else {
		return nil, errors.New("no zlib magic number")
	}

	return buffer, nil
}

//decode will decode reader to object
func decode(reader io.Reader, object interface{}) error {
	if reader == nil {
		return errors.New("reader nil")
	}
	if object == nil {
		return errors.New("forest nil")
	}

	decoder := gob.NewDecoder(reader)
	if errDecoder := decoder.Decode(object); errDecoder != nil {
		log.Println(errDecoder)
		return errDecoder
	}

	return nil
}

//ByteToForest Converts byte array representation of object
func ByteToForest(byteArray []byte, object interface{}) error {
	if byteArray == nil {
		return errors.New("forestbyte nil")
	}
	if object == nil {
		return errors.New("forest nil")
	}

	buffer := bytes.NewBuffer(byteArray)

	if ok, err := IsCompressed(buffer); err == nil && ok {

		if unCompressBuff, errUnComp := UnCompress(buffer); errUnComp == nil {

			return decode(unCompressBuff, object)
		} else {
			return errUnComp
		}

	} else if err != nil {
		return err
	} else {
		return decode(buffer, object)
	}
}
