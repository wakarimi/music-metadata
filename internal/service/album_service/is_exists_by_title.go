package album_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s Service) IsExistsByTitle(tx *sqlx.Tx, title string) (exists bool, err error) {
	log.Debug().Str("title", title).Msg("Checking album existence")

	exists, err = s.AlbumRepo.IsExistsByTitle(tx, title)
	if err != nil {
		log.Error().Err(err).Str("title", title).Msg("Failed to check album existence")
		return false, err
	}

	log.Debug().Str("title", title).Bool("exists", exists).Msg("Album existence checked successfully")
	return exists, err
}
