package repositoriesgorm

import (
	"errors"

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
	if err := c.gorm.Create(&Tags{
		ID:            tag.ID,
		Active:        tag.Active,
		CreatedAt:     tag.CreatedAt,
		UpdatedAt:     tag.UpdatedAt,
		DeactivatedAt: tag.DeactivatedAt,
		UserID:        tag.UserID,
		Name:          tag.Name,
		Color:         tag.Color,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (c *TagRepository) DeleteTag(tag entities.Tag) error {
	if err := c.gorm.Model(&Tags{}).Where("id = ? AND user_id = ?", tag.ID, tag.UserID).Select("Active", "DeactivatedAt").Updates(map[string]interface{}{"active": tag.Active, "updated_at": tag.DeactivatedAt}).Error; err != nil {
		return err
	}

	return nil
}

func (c *TagRepository) GetTags(userID string) ([]entities.Tag, error) {
	var tagsModel []Tags
	if err := c.gorm.Where("user_id = ?", userID).Find(&tagsModel).Error; err != nil {
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
	}

	return tags, nil
}

func (c *TagRepository) GetTag(userID string, tagID string) (entities.Tag, error) {
	var tagModel Tags

	result := c.gorm.Model(&Tags{}).Where("id = ? AND user_id = ?", tagID, userID).First(&tagModel)
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

	result := c.gorm.Model(&Tags{}).Where("name = ? AND user_id = ?", tagName, userID).First(&tagModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("tag not found")
		}
		return false, errors.New(result.Error.Error())
	}

	return true, nil
}
