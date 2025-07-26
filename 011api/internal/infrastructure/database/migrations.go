package database

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

type Migration struct {
	Version int
	Name	string
	UpSQL	string
	DownSQL	string
}

type Migrator struct {
	db *DB
}

func NewMigrator(db *DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) CreateMigrationTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := m.db.Exec(query)
	return err
}

func (m *Migrator) GetAppliedMigrations() ([]int, error) {
	query := "SELECT version FROM schema_migrations ORDER BY version"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []int
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, rows.Err()
}

func (m *Migrator) ApplyMigration(migration Migration) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(migration.UpSQL); err != nil {
		return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
	}

	query := "INSERT INTO schema_migrations (version, name) VALUES ($1, $2)"
	if _, err := tx.Exec(query, migration.Version, migration.Name); err != nil {
		return fmt.Errorf("failed to record migration %d: %s", migration.Version, err)
	}

	return tx.Commit()
}

func (m *Migrator) RollbackMigration(migration Migration) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(migration.DownSQL); err != nil {
		return fmt.Errorf("failed to rollback migration %d: %w", migration.Version, err)
	}

	query := "DELETE FROM schema_migrations WHERE version = $1"
	if _, err := tx.Exec(query, migration.Version); err != nil {
		return fmt.Errorf("failed to remove migration record %d: %w", migration.Version, err)
	}

	return tx.Commit()
}

func LoadMigrationsFromFS(migrationFS fs.FS) ([]Migration, error) {
	var migrations []Migration
	migrationMap := make(map[int]*Migration)

	err := fs.WalkDir(migrationFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".sql") {
			return nil
		}

		//parse filename: 001_create_users_table.up.sql
		filename := filepath.Base(path)
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			return nil
		}

		var version int
		if _, err := fmt.Sscanf(parts[0], "%d", &version); err != nil {
			return nil
		}

		content, err := fs.ReadFile(migrationFS, path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", path, err)
		}

		migration, exists := migrationMap[version]
		if !exists {
			migration = &Migration{Version: version}
			migrationMap[version] = migration
		}

		if strings.Contains(filename, ".up.sql") {
			migration.UpSQL = string(content)
			//extract name from filename
			nameParts := strings.Split(filename, ".")
			if len(nameParts) >= 2 {
				migration.Name = strings.Join(strings.Split(nameParts[0], "_")[1:], "_")
			}
		} else if strings.Contains(filename, ".down.sql") {
			migration.DownSQL = string(content)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, migration := range migrationMap {
		migrations = append(migrations, *migration)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}