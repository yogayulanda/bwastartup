package campaign

type Service interface {
	FindCampaigns(UserID int) ([]Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) FindCampaigns(UserID int) ([]Campaign, error) {
	//cek ada user id atau tidak
	//kalau ada
	if UserID != 0 {
		campaigns, err := s.repository.FindByUserID(UserID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}
	// kalau userid kosong
	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}
