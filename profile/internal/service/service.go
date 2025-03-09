package service

type ProfileRepo interface{}

type profileService struct {
	repo ProfileRepo
}

func NewProfileService(repo ProfileRepo) *profileService {
	return &profileService{repo: repo}
}
