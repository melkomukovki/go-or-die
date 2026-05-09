package part

type service struct {
	partRepo PartRepository
}

func NewService(partRepo PartRepository) *service {
	return &service{
		partRepo: partRepo,
	}
}
