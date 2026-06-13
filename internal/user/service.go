package user

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) FindAll() ([]User, error) {
	return s.repo.FindAll()
}

func (s *Service) FindById(id uint) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Create(user *User) error {
	return s.repo.Create(user)
}

func (s *Service) Update(user *User) error {

	return s.repo.Update(user)
}

func (s *Service) Delete(id uint) error {

	return s.repo.Delete(id)
}
