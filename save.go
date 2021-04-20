package randomforest

import (
	"errors"

	internal "github.com/MandelV/randomForest/v2/internal"
)

//Save Will save the state of the forest into file
func (forest *Forest) Save(folder string, compress bool) (string, error) {
	return internal.Save(&forest, folder, compress)
}

func (deepforest *DeepForest) Save(folder string, compress bool) (string, error) {
	return internal.Save(&deepforest, folder, compress)
}

func (forest *Forest) ToBytes(compress bool) ([]byte, error) {
	return internal.ToBytes(&forest, compress)
}

func (deepForest *DeepForest) ToBytes(compress bool) ([]byte, error) {
	return internal.ToBytes(&deepForest, compress)
}

func BytesToForest(bytes []byte) (*Forest, error) {
	var forest *Forest = &Forest{}
	var err error
	if err = internal.ByteToForest(bytes, forest); err == nil {
		return forest, nil
	}
	return nil, err
}

func BytesToDeepForest(bytes []byte) (*DeepForest, error) {
	var deepForest *DeepForest = &DeepForest{}
	var err error
	if err = internal.ByteToForest(bytes, deepForest); err == nil {
		return deepForest, nil
	}
	return nil, err
}

func LoadForest(path string) (*Forest, error) {
	var forest *Forest = &Forest{}
	var err error
	if err = internal.Load(path, forest); err == nil && forest != nil {
		return forest, nil
	}
	if forest == nil {
		return nil, errors.New("Forest nil")
	}
	return nil, err
}

func LoadDeepForest(path string) (*DeepForest, error) {
	var deepForest *DeepForest = &DeepForest{}
	var err error
	if err = internal.Load(path, deepForest); err == nil {
		return deepForest, nil
	}
	return nil, err
}
