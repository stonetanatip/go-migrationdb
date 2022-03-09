package migrate

import "database/sql"

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Up(req *Request, filePath map[string]string) error {
	return nil
}

func (s *Service) Down(req *Request, filePath map[string]string) error {
	return nil
}
