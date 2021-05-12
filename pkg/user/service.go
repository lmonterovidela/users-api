package user

type IService interface {
	Create(user *User) (int64, error)
	Update(userId int, user *User) error
	Get(userId int) (User, error)
	Find(name string, size int, offset int) (UserList, error)
	Delete(userId int) error
}

type Service struct {
	repository IRepository
}

func (s *Service) Create(user *User) (int64, error) {
	return s.repository.Insert(user)
}

func (s *Service) Update(userId int, user *User) error {
	return s.repository.Update(userId, user)
}

func (s *Service) Get(userId int) (User, error) {
	return s.repository.Get(userId)
}

func (s *Service) Find(name string, size int, offset int) (UserList, error) {
	return s.repository.Find(name, size, offset)
}

func (s *Service) Delete(userId int) error {
	return s.repository.Delete(userId)
}

func NewService() IService {
	return &Service{
		repository: NewRepository(),
	}
}
