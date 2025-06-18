package services

import "Netlfy/repositories"

type TagService interface {
	GetTags() ([]string, error)
}

type TagServiceImpl struct {
	tagRepo repositories.TagRepository
}

func NewTagService(tagRepo repositories.TagRepository) TagService {
	return &TagServiceImpl{tagRepo}
}

func (s *TagServiceImpl) GetTags() ([]string, error) {
	tags, err := s.tagRepo.GetAllTags()
	if err != nil {
		return nil, err
	}
	var result []string
	for _, tag := range tags {
		result = append(result, tag.Name)
	}
	return result, nil
}
