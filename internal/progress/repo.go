package progress

import (
	"database/sql"
	"fmt"
	"time"
)

// Status — статус прохождения урока.
type Status string

const (
	StatusNew     Status = "new"
	StatusReading Status = "reading"
	StatusDone    Status = "done"
)

// Progress — прогресс по уроку.
type Progress struct {
	LessonID     int64
	Status       Status
	PracticeDone bool
	PointsEarned int
	UpdatedAt    time.Time
}

// Note — заметка к уроку.
type Note struct {
	LessonID  int64
	NoteMD    string
	UpdatedAt time.Time
}

// Submission — отправка решения.
type Submission struct {
	ID        int64
	TaskID    int64
	Code      string
	Status    string // pending, success, error, timeout
	Stdout    string
	Stderr    string
	CreatedAt time.Time
}

// Stats — общая статистика.
type Stats struct {
	TotalLessons    int
	CompletedCount  int
	InProgressCount int
	TotalPoints     int
	EarnedPoints    int
}

// Repository — репозиторий для работы с прогрессом.
type Repository struct {
	db *sql.DB
}

// NewRepository создаёт новый репозиторий.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// --- Progress ---

// GetProgress возвращает прогресс по уроку.
func (r *Repository) GetProgress(lessonID int64) (*Progress, error) {
	p := &Progress{}
	err := r.db.QueryRow(
		`SELECT lesson_id, status, practice_done, points_earned, updated_at 
		 FROM progress WHERE lesson_id = ?`,
		lessonID,
	).Scan(&p.LessonID, &p.Status, &p.PracticeDone, &p.PointsEarned, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		// Возвращаем дефолтный прогресс
		return &Progress{
			LessonID:     lessonID,
			Status:       StatusNew,
			PracticeDone: false,
			PointsEarned: 0,
			UpdatedAt:    time.Now(),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get progress: %w", err)
	}

	return p, nil
}

// UpdateProgress обновляет прогресс по уроку.
func (r *Repository) UpdateProgress(p *Progress) error {
	_, err := r.db.Exec(
		`INSERT INTO progress (lesson_id, status, practice_done, points_earned, updated_at)
		 VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		 ON CONFLICT(lesson_id) DO UPDATE SET 
		   status = excluded.status,
		   practice_done = excluded.practice_done,
		   points_earned = excluded.points_earned,
		   updated_at = CURRENT_TIMESTAMP`,
		p.LessonID, p.Status, p.PracticeDone, p.PointsEarned,
	)
	if err != nil {
		return fmt.Errorf("update progress: %w", err)
	}
	return nil
}

// SetStatus устанавливает статус урока.
func (r *Repository) SetStatus(lessonID int64, status Status) error {
	_, err := r.db.Exec(
		`INSERT INTO progress (lesson_id, status, updated_at)
		 VALUES (?, ?, CURRENT_TIMESTAMP)
		 ON CONFLICT(lesson_id) DO UPDATE SET 
		   status = excluded.status,
		   updated_at = CURRENT_TIMESTAMP`,
		lessonID, status,
	)
	return err
}

// SetPracticeDone отмечает практику как выполненную.
func (r *Repository) SetPracticeDone(lessonID int64, points int) error {
	_, err := r.db.Exec(
		`INSERT INTO progress (lesson_id, practice_done, points_earned, updated_at)
		 VALUES (?, 1, ?, CURRENT_TIMESTAMP)
		 ON CONFLICT(lesson_id) DO UPDATE SET 
		   practice_done = 1,
		   points_earned = points_earned + excluded.points_earned,
		   updated_at = CURRENT_TIMESTAMP`,
		lessonID, points,
	)
	return err
}

// GetAllProgress возвращает прогресс по всем урокам.
func (r *Repository) GetAllProgress() (map[int64]*Progress, error) {
	rows, err := r.db.Query(
		`SELECT lesson_id, status, practice_done, points_earned, updated_at FROM progress`,
	)
	if err != nil {
		return nil, fmt.Errorf("get all progress: %w", err)
	}
	defer rows.Close()

	result := make(map[int64]*Progress)
	for rows.Next() {
		p := &Progress{}
		if err := rows.Scan(&p.LessonID, &p.Status, &p.PracticeDone, &p.PointsEarned, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan progress: %w", err)
		}
		result[p.LessonID] = p
	}

	return result, rows.Err()
}

// --- Notes ---

// GetNote возвращает заметку к уроку.
func (r *Repository) GetNote(lessonID int64) (*Note, error) {
	n := &Note{}
	err := r.db.QueryRow(
		`SELECT lesson_id, note_md, updated_at FROM notes WHERE lesson_id = ?`,
		lessonID,
	).Scan(&n.LessonID, &n.NoteMD, &n.UpdatedAt)

	if err == sql.ErrNoRows {
		return &Note{LessonID: lessonID, NoteMD: "", UpdatedAt: time.Now()}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get note: %w", err)
	}

	return n, nil
}

// SaveNote сохраняет заметку к уроку.
func (r *Repository) SaveNote(lessonID int64, noteMD string) error {
	_, err := r.db.Exec(
		`INSERT INTO notes (lesson_id, note_md, updated_at)
		 VALUES (?, ?, CURRENT_TIMESTAMP)
		 ON CONFLICT(lesson_id) DO UPDATE SET 
		   note_md = excluded.note_md,
		   updated_at = CURRENT_TIMESTAMP`,
		lessonID, noteMD,
	)
	return err
}

// --- Submissions ---

// CreateSubmission создаёт запись об отправке решения.
func (r *Repository) CreateSubmission(s *Submission) error {
	result, err := r.db.Exec(
		`INSERT INTO submissions (task_id, code, status, stdout, stderr)
		 VALUES (?, ?, ?, ?, ?)`,
		s.TaskID, s.Code, s.Status, s.Stdout, s.Stderr,
	)
	if err != nil {
		return fmt.Errorf("create submission: %w", err)
	}
	s.ID, _ = result.LastInsertId()
	return nil
}

// UpdateSubmission обновляет статус отправки.
func (r *Repository) UpdateSubmission(s *Submission) error {
	_, err := r.db.Exec(
		`UPDATE submissions SET status = ?, stdout = ?, stderr = ? WHERE id = ?`,
		s.Status, s.Stdout, s.Stderr, s.ID,
	)
	return err
}

// IsTaskSolvedSuccessfully проверяет, было ли задание уже успешно решено.
func (r *Repository) IsTaskSolvedSuccessfully(taskID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM submissions WHERE task_id = ? AND status = 'success'`,
		taskID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check task solved: %w", err)
	}
	return count > 0, nil
}

// GetSubmissionsByTaskID возвращает отправки по заданию.
func (r *Repository) GetSubmissionsByTaskID(taskID int64, limit int) ([]Submission, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := r.db.Query(
		`SELECT id, task_id, code, status, stdout, stderr, created_at 
		 FROM submissions WHERE task_id = ? ORDER BY created_at DESC LIMIT ?`,
		taskID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("get submissions: %w", err)
	}
	defer rows.Close()

	var submissions []Submission
	for rows.Next() {
		var s Submission
		if err := rows.Scan(&s.ID, &s.TaskID, &s.Code, &s.Status, &s.Stdout, &s.Stderr, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan submission: %w", err)
		}
		submissions = append(submissions, s)
	}

	return submissions, rows.Err()
}

// --- Stats ---

// ResetAllProgress сбрасывает весь прогресс (очки, статусы, отправки).
func (r *Repository) ResetAllProgress() error {
	// Удаляем все отправки
	if _, err := r.db.Exec(`DELETE FROM submissions`); err != nil {
		return fmt.Errorf("delete submissions: %w", err)
	}
	// Удаляем весь прогресс
	if _, err := r.db.Exec(`DELETE FROM progress`); err != nil {
		return fmt.Errorf("delete progress: %w", err)
	}
	// Заметки оставляем — они полезны
	return nil
}

// GetStats возвращает общую статистику.
func (r *Repository) GetStats() (*Stats, error) {
	stats := &Stats{}

	// Общее количество уроков
	err := r.db.QueryRow(`SELECT COUNT(*) FROM lessons`).Scan(&stats.TotalLessons)
	if err != nil {
		return nil, fmt.Errorf("count lessons: %w", err)
	}

	// Завершённые уроки
	err = r.db.QueryRow(`SELECT COUNT(*) FROM progress WHERE status = 'done'`).Scan(&stats.CompletedCount)
	if err != nil {
		return nil, fmt.Errorf("count completed: %w", err)
	}

	// В процессе
	err = r.db.QueryRow(`SELECT COUNT(*) FROM progress WHERE status = 'reading'`).Scan(&stats.InProgressCount)
	if err != nil {
		return nil, fmt.Errorf("count in progress: %w", err)
	}

	// Общее количество очков
	err = r.db.QueryRow(`SELECT COALESCE(SUM(points), 0) FROM tasks`).Scan(&stats.TotalPoints)
	if err != nil {
		return nil, fmt.Errorf("sum total points: %w", err)
	}

	// Заработанные очки
	err = r.db.QueryRow(`SELECT COALESCE(SUM(points_earned), 0) FROM progress`).Scan(&stats.EarnedPoints)
	if err != nil {
		return nil, fmt.Errorf("sum earned points: %w", err)
	}

	return stats, nil
}
