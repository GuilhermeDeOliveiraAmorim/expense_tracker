package repositoriesgorm

import (
	"errors"
	"sort"

	"github.com/GuilhermeDeOliveiraAmorim/expense-tracker/internal/entities"
	"gorm.io/gorm"
)

type TagRepository struct {
	gorm *gorm.DB
}

func NewTagRepository(gorm *gorm.DB) *TagRepository {
	return &TagRepository{
		gorm: gorm,
	}
}

func (c *TagRepository) CreateTag(tag entities.Tag) error {
	tx := c.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Create(&Tags{
		ID:            tag.ID,
		Active:        tag.Active,
		CreatedAt:     tag.CreatedAt,
		UpdatedAt:     tag.UpdatedAt,
		DeactivatedAt: tag.DeactivatedAt,
		UserID:        tag.UserID,
		Name:          tag.Name,
		Color:         tag.Color,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (c *TagRepository) DeleteTag(tag entities.Tag) error {
	tx := c.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	result := tx.Model(&Tags{}).Where("id = ? AND user_id = ? AND active = ?", tag.ID, tag.UserID, true).
		Select("Active", "DeactivatedAt", "UpdatedAt").Updates(Tags{
		Active:        tag.Active,
		DeactivatedAt: tag.DeactivatedAt,
		UpdatedAt:     tag.UpdatedAt,
	})

	if result.Error != nil {
		tx.Rollback()
		return errors.New(result.Error.Error())
	}

	return tx.Commit().Error
}

func (c *TagRepository) GetTags(userID string) ([]entities.Tag, error) {
	var tagsModel []Tags
	if err := c.gorm.Where("user_id = ? AND active = ?", userID, true).Find(&tagsModel).Order("created_at DESC").Error; err != nil {
		return nil, err
	}

	var tags []entities.Tag

	if len(tagsModel) > 0 {
		for _, tagModel := range tagsModel {
			tag := entities.Tag{
				SharedEntity: entities.SharedEntity{
					ID:            tagModel.ID,
					Active:        tagModel.Active,
					CreatedAt:     tagModel.CreatedAt,
					UpdatedAt:     tagModel.UpdatedAt,
					DeactivatedAt: tagModel.DeactivatedAt,
				},
				UserID: tagModel.UserID,
				Name:   tagModel.Name,
				Color:  tagModel.Color,
			}

			tags = append(tags, tag)
		}

		sort.Slice(tags, func(i, j int) bool {
			return tags[i].CreatedAt.After(tags[j].CreatedAt)
		})
	} else {
		tags = []entities.Tag{}
	}

	return tags, nil
}

func (c *TagRepository) GetTag(userID string, tagID string) (entities.Tag, error) {
	var tagModel Tags

	result := c.gorm.Model(&Tags{}).Where("id = ? AND user_id = ? AND active = ?", tagID, userID, true).First(&tagModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.Tag{}, errors.New("tag not found")
		}
		return entities.Tag{}, errors.New(result.Error.Error())
	}

	tag := entities.Tag{
		SharedEntity: entities.SharedEntity{
			ID:            tagModel.ID,
			Active:        tagModel.Active,
			CreatedAt:     tagModel.CreatedAt,
			UpdatedAt:     tagModel.UpdatedAt,
			DeactivatedAt: tagModel.DeactivatedAt,
		},
		UserID: tagModel.UserID,
		Name:   tagModel.Name,
		Color:  tagModel.Color,
	}

	return tag, nil
}

func (c *TagRepository) ThisTagExists(userID string, tagName string) (bool, error) {
	var tagModel Tags

	result := c.gorm.Model(&Tags{}).Where("name = ? AND user_id = ? AND active = ?", tagName, userID, true).First(&tagModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("tag not found")
		}
		return false, errors.New(result.Error.Error())
	}

	return true, nil
}

func (c *TagRepository) UpdateTag(tag entities.Tag) error {
	tx := c.gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	result := tx.Model(&Tags{}).Where("id = ? AND user_id = ? AND active = ?", tag.ID, tag.UserID, true).Updates(Tags{
		Name:      tag.Name,
		Color:     tag.Color,
		UpdatedAt: tag.UpdatedAt,
	})

	if result.Error != nil {
		tx.Rollback()
		return errors.New(result.Error.Error())
	}

	return tx.Commit().Error
}
