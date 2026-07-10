package db

import (
	"database/sql"
	"time"
)

type Rule struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	ConditionJSON  []byte    `json:"condition_json"`
	Action         string    `json:"action"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

type Stats struct {
	TotalEvents  int64   `json:"total_events"`
	BlockedCount int64   `json:"blocked_count"`
	WarnCount    int64   `json:"warn_count"`
	PassCount    int64   `json:"pass_count"`
	AvgScore     float64 `json:"avg_score"`
}

func GetRules(db *sql.DB) ([]Rule, error) {
	rows, err := db.Query("SELECT id, name, condition_json, action, is_active, created_at FROM detection_rules ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []Rule
	for rows.Next() {
		var r Rule
		if err := rows.Scan(&r.ID, &r.Name, &r.ConditionJSON, &r.Action, &r.IsActive, &r.CreatedAt); err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}

func CreateRule(db *sql.DB, name string, conditionJSON []byte, action string) error {
	_, err := db.Exec(
		"INSERT INTO detection_rules (name, condition_json, action, is_active) VALUES ($1, $2, $3, true)",
		name, conditionJSON, action,
	)
	return err
}

func GetStats(db *sql.DB, from, to time.Time) (*Stats, error) {
	stats := &Stats{}
	err := db.QueryRow(`
		SELECT 
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN verdict = 'BLOCK' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN verdict = 'WARN' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN verdict = 'PASS' THEN 1 ELSE 0 END), 0),
			COALESCE(AVG(score), 0)
		FROM events 
		WHERE timestamp BETWEEN $1 AND $2
	`, from, to).Scan(&stats.TotalEvents, &stats.BlockedCount, &stats.WarnCount, &stats.PassCount, &stats.AvgScore)
	return stats, err
}

func UpdateEventVerdict(db *sql.DB, eventID string, verdict string) error {
	_, err := db.Exec("UPDATE events SET verdict = $1 WHERE id = $2", verdict, eventID)
	return err
}
