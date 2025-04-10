package processing

import (
	"fmt"
	"image"
	"time"

	"github.com/charmbracelet/log"
	"github.com/itspacchu/anilist-chart/anilist"
)

type Anime struct {
	Name  string
	Cover string
	Count int
}

func (a *Anime) CountUp() {
	a.Count += 1
}

func ProcessChart(username string, timeDays int, activity_type string) *image.RGBA {
	chartMap := make(map[int64]Anime, 10)
	if timeDays == 0 {
		timeDays = 7
	}
	userID := anilist.FetchIDFromUsername(username)
	epoch := time.Now().Add(-time.Duration(timeDays) * time.Hour * 24).Unix()
	var activity anilist.Response = anilist.FetchActivitiesDetails(userID, epoch, activity_type)
	for _, activity := range activity.Data.Page.Activities {
		if anime, ok := chartMap[activity.Media.ID]; ok {
			anime.CountUp()
			chartMap[activity.Media.ID] = anime
		} else {
			log.Debugf("Got %s", activity.Media.Title.Romaji)
			chartMap[activity.Media.ID] = Anime{
				Name:  activity.Media.Title.Romaji,
				Cover: activity.Media.CoverImage.Large,
			}
		}
	}
	return GenerateAnimeGridImage(chartMap, 500, 10, fmt.Sprintf("%s.jpeg", username))
}
