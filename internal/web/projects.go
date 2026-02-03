package web

import (
	"net/http"

	"golearning"
)

type Project struct {
	ID       string
	Title    string
	Subtitle string
	SpecMD   string
}

func (s *Server) handleProjects(w http.ResponseWriter, r *http.Request) {
	stats, _ := s.progressRepo.GetStats()

	projects := []Project{
		{
			ID:       "capstone-rest",
			Title:    "Capstone REST: сервис заказов (Gin + Postgres)",
			Subtitle: "JWT, миграции, интеграционные тесты, CI, Docker Compose, метрики/логи/трейсы, нагрузка и профили",
			SpecMD:   golearning.CapstoneRESTSpecMD,
		},
		{
			ID:       "capstone-grpc",
			Title:    "Capstone gRPC: Users/Accounts сервис (gRPC + TLS/mTLS)",
			Subtitle: "Interceptors, deadlines, безопасность, наблюдаемость; опционально grpc-gateway + OpenAPI",
			SpecMD:   golearning.CapstoneGRPCSpecMD,
		},
	}

	data := map[string]interface{}{
		"Stats":    stats,
		"Projects": projects,
	}

	s.render(w, "projects.html", data)
}
