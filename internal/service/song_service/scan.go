package song_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-metadata/internal/client/music_files_client/audio_file_client"
	"music-metadata/internal/model"
)

func (s *Service) Scan(tx *sqlx.Tx) (err error) {
	log.Debug().Msg("Scanning songs")

	audioFiles, err := s.AudioFileClient.GetAll()
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch audio files")
		return err
	}

	songs, err := s.SongRepo.ReadAll(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get songs")
		return err
	}

	onlyAudioFileExists := make([]audio_file_client.GetAllResponseItem, 0)
	onlySongExists := make([]model.Song, 0)
	audioFilesWithChangedId := make([]audio_file_client.GetAllResponseItem, 0)
	songsWithChangedContent := make([]model.Song, 0)

	for _, audioFile := range audioFiles {
		processed := false
		for _, song := range songs {
			if (audioFile.AudioFileId == song.AudioFileId) && (audioFile.Sha256 == song.Sha256) {
				processed = true
				break
			} else if (audioFile.AudioFileId == song.AudioFileId) && (audioFile.Sha256 != song.Sha256) {
				songsWithChangedContent = append(songsWithChangedContent, song)
				processed = true
				break
			} else if (audioFile.AudioFileId != song.AudioFileId) && (audioFile.Sha256 == song.Sha256) {
				audioFilesWithChangedId = append(audioFilesWithChangedId, audioFile)
				processed = true
				break
			}
		}
		if processed {
			continue
		}

		onlyAudioFileExists = append(onlyAudioFileExists, audioFile)
	}
	for _, song := range songs {
		processed := false
		for _, audioFile := range audioFiles {
			if (audioFile.AudioFileId == song.AudioFileId) || (audioFile.Sha256 == song.Sha256) {
				processed = true
				break
			}
		}
		if !processed {
			onlySongExists = append(onlySongExists, song)
		}
	}

	err = s.createMissedSongs(tx, onlyAudioFileExists)
	if err != nil {
		log.Error().Err(err).Msg("Failed to created missed songs")
		return err
	}

	err = s.removeObsoleteSongs(tx, onlySongExists)
	if err != nil {
		log.Error().Err(err).Msg("Failed to remove obsolete songs")
		return err
	}

	err = s.updateSongsWithChangedAudioFileId(tx, audioFilesWithChangedId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update songs with changed id")
		return err
	}

	err = s.updateSongsWithChangedContent(tx, songsWithChangedContent)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update songs with changed content")
		return err
	}

	err = s.AlbumService.RemoveUnnecessaryItems(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to remove unnecessary albums")
		return err
	}

	err = s.ArtistService.RemoveUnnecessaryItems(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to remove unnecessary artists")
		return err
	}

	err = s.GenreService.RemoveUnnecessaryItems(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to remove unnecessary genres")
		return err
	}

	log.Debug().Msg("Songs scanned successfully")
	return nil
}

func (s *Service) createMissedSongs(tx *sqlx.Tx, audioFiles []audio_file_client.GetAllResponseItem) (err error) {
	for _, audioFile := range audioFiles {
		song, err := s.SongByAudioFileWithoutSha(tx, audioFile.AudioFileId)
		if err != nil {
			log.Error().Err(err).Int("audioFileId", audioFile.AudioFileId).Msg("Failed to prepare song")
			return err
		}
		song.Sha256 = audioFile.Sha256
		_, err = s.SongRepo.Create(tx, song)
		if err != nil {
			log.Error().Err(err).Int("audioFileId", audioFile.AudioFileId).Msg("Failed to create song")
			return err
		}
	}
	return nil
}

func (s *Service) removeObsoleteSongs(tx *sqlx.Tx, songs []model.Song) (err error) {
	for _, song := range songs {
		err = s.SongRepo.Delete(tx, song.SongId)
		if err != nil {
			log.Error().Err(err).Int("songId", song.SongId).Msg("Failed to delete song")
			return err
		}
	}
	return nil
}

func (s *Service) updateSongsWithChangedAudioFileId(tx *sqlx.Tx, audioFiles []audio_file_client.GetAllResponseItem) (err error) {
	songs, err := s.SongRepo.ReadAll(tx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all songs")
		return err
	}

	for _, audioFile := range audioFiles {
		songId := -1
		for _, song := range songs {
			if song.Sha256 == audioFile.Sha256 {
				songId = song.SongId
				break
			}
		}
		err := s.SongRepo.UpdateAudioFileId(tx, songId, audioFile.AudioFileId)
		if err != nil {
			log.Error().Err(err).Int("songId", songId).Msg("Failed to update audio file id")
			return err
		}
	}
	return nil
}

func (s *Service) updateSongsWithChangedContent(tx *sqlx.Tx, songs []model.Song) (err error) {
	for _, song := range songs {
		audioFile, err := s.AudioFileClient.Get(song.AudioFileId)
		if err != nil {
			log.Error().Err(err).Int("songId", song.SongId).Int("audioFileId", song.AudioFileId).Msg("Failed to fetch audio file")
			return err
		}

		newSong, err := s.SongByAudioFileWithoutSha(tx, audioFile.AudioFileId)
		if err != nil {
			log.Error().Err(err).Int("audioFileId", audioFile.AudioFileId).Msg("Failed to prepare newSong")
			return err
		}
		newSong.Sha256 = audioFile.Sha256
		err = s.SongRepo.Update(tx, song.SongId, newSong)
		if err != nil {
			log.Error().Err(err).Int("audioFileId", audioFile.AudioFileId).Msg("Failed to create newSong")
			return err
		}
	}
	return nil
}
