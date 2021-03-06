package types

import (
	"strings"

	youtubeapiinter "github.com/vidur2/blindOsuBackend/youtubeApiInter"
)

type VideoReq struct {
	VideoId string `json:"video_id"`
}

func (v *VideoReq) TranslateVideoId() error {
	if len(v.VideoId) != 11 || len(strings.Split(v.VideoId, " ")) > 1 {
		id, err := youtubeapiinter.GetVideoUrl(v.VideoId)
		if err != nil {
			return err
		}
		v.VideoId = id
	}

	return nil
}

type VideoRes struct {
	Base64Url string               `json:"base64_url"`
	Game      []RelativeModelCoord `json:"model_coord"`
}

type ErrorRes struct {
	Err string `json:"err"`
}

type PitchCoordinate struct {
	Time  float32 `json:"time"`
	Pitch float32 `json:"pitch"`
}

type VideoResYin struct {
	PitchCoords []PitchCoordinate `json:"pitch_coordinate"`
	Base64Url   string            `json:"base64_url"`
}

func (v *VideoResYin) TranslateToRelative() {
	maxVal := v.findMax()

	for idx, coord := range v.PitchCoords {
		v.PitchCoords[idx] = PitchCoordinate{
			Time:  coord.Time,
			Pitch: coord.Pitch / maxVal,
		}
	}
}

func (v *VideoResYin) findMax() float32 {
	max := float32(0.)

	for _, coord := range v.PitchCoords {
		if coord.Pitch > max {
			max = coord.Pitch
		}
	}

	return max
}
