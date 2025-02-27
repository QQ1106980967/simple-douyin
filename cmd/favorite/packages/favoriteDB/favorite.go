package favoriteDB

import (
	// "log"

	"errors"

	// "github.com/karry-almond/model"
	"github.com/AgSword/simpleDouyin/cmd/favorite/packages/model"
	"gorm.io/gorm"
	// "golang.org/x/tools/go/analysis/passes/nilfunc"
)

// status 返回0——成功，返回1——失败
// err 返回nil——成功，返回其他——失败原因
func NewFavorite(user_id int64, video_id int64) (status int32, err error) {

	// 创建一条favorite数据
	favorite := model.Favorite{
		//TODO:ID这里不是逐主键
		ID:      1,
		UserId:  user_id,
		VideoId: video_id}

	//新建喜欢、新增喜欢为同一事务
	err = Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&favorite).Error; err != nil {
			return err
		}
		//更改对应video的favorite_count
		var video model.Video

		if err := tx.Select("*").First(&model.Video{ID: video_id}).Scan(&video).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Video{ID: video_id}).Update("favorite_count", video.FavoriteCount+1).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 1, err
	}

	return 0, nil

}

func CancelFavorite(user_id int64, video_id int64) (status int32, err error) {
	//先根据user_id和video_id寻找到id，再根据id软删除
	var favorite model.Favorite

	err = Db.Transaction(func(tx *gorm.DB) error {
		//通过user_id和video_id找到要删除的favorite记录
		if err := tx.Select("*").First(&model.Favorite{UserId: user_id, VideoId: video_id}).Scan(&favorite).Error; err != nil {
			return err
		}
		//删除这条favorite记录
		if err := tx.Delete(&favorite).Error; err != nil {
			return err
		}
		//更改对应video的favorite_count
		var video model.Video

		if err := tx.Select("*").First(&model.Video{ID: video_id}).Scan(&video).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Video{ID: video_id}).Update("favorite_count", video.FavoriteCount+1).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 1, err
	}

	return 0, nil
}

func GetFavoriteList(user_id int64) (status int32, videoList []model.Video, err error) {
	var favoriteList []model.Favorite

	//根据user_id，在favorite表中找到他的所有favorite记录
	if err = Db.Select("video_id").Where(&model.Favorite{UserId: user_id}).Find(&favoriteList).Error; err != nil {
		return 1, nil, err
	}
	//根据favorite记录找到所有对应video_id的video
	var video model.Video
	for _, favoritelog := range favoriteList {
		if err = Db.Where(&model.Video{ID: favoritelog.VideoId}).First(&video).Error; err != nil {
			return 1, videoList, errors.New("failed finding all the favorite videos")
		}
		videoList = append(videoList, video)
	}

	return 0, videoList, err
}
