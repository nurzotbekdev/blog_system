package models

import "gorm.io/gorm"

type PlaylistVideo struct {
	gorm.Model
	PlaylistID uint     `json:"playlist_id" gorm:"index:idx_playlist_video,unique"`
	Playlist   Playlist `json:"-" gorm:"foreignKey:PlaylistID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VideoID    uint     `json:"video_id" gorm:"index:idx_playlist_video,unique"`
	Video      Video    `json:"-" gorm:"foreignKey:VideoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
