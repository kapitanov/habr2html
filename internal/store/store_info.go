package store

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
)

const RootInfoFileName = "info.json"

type rootInfoJSON struct {
	Count    int                         `json:"count"`
	Articles map[string]*articleInfoJSON `json:"articles"`
}

type articleInfoJSON struct {
}

func loadRootInfoJSON(infoFilePath string) (*rootInfoJSON, error) {
	var rootInfo = rootInfoJSON{
		Articles: make(map[string]*articleInfoJSON),
	}

	f, err := os.Open(infoFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Debug().Str("path", infoFilePath).Msg("missing info.json")
			return &rootInfo, nil
		}

		log.Error().Str("path", infoFilePath).Err(err).Msg("unable to open info.json")
		return nil, err
	}

	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Error().Str("path", infoFilePath).Err(err).Msg("unable to open info.json")
		return nil, err
	}

	err = json.Unmarshal(bytes, &rootInfo)
	if err != nil {
		log.Error().Str("path", infoFilePath).Err(err).Msg("unable to open info.json")
		return nil, err
	}

	return &rootInfo, nil
}

func (info *rootInfoJSON) Save(infoFilePath string) error {
	info.Count = len(info.Articles)

	bytes, err := json.MarshalIndent(info, "", "    ")
	if err != nil {
		log.Error().Err(err).Msg("unable to save info.json")
		return err
	}

	f, err := os.Create(infoFilePath)
	if err != nil {
		log.Error().Str("path", infoFilePath).Err(err).Msg("unable to save info.json")
		return err
	}

	defer f.Close()
	_, err = f.Write(bytes)
	if err != nil {
		log.Error().Str("path", infoFilePath).Err(err).Msg("unable to save info.json")
		return err
	}

	return nil
}
