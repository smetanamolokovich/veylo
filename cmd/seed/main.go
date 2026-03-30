package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	orgID       = "seed-org-veylo-demo"
	workflowID  = "seed-wfl-001"
	userAdmin   = "seed-usr-admin"
	userManager = "seed-usr-manager"
	userInsp    = "seed-usr-inspector"
	userEval    = "seed-usr-evaluator"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:6543/veylo?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	for i := range 10 {
		if err := db.Ping(); err != nil {
			if i == 9 {
				log.Fatalf("db not reachable: %v", err)
			}
			time.Sleep(time.Second)
			continue
		}
		break
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}
	pwdHash := string(hash)

	steps := []struct {
		name string
		fn   func(*sql.DB) error
	}{
		{"organization", seedOrg},
		{"workflow", seedWorkflow},
		{"users", func(db *sql.DB) error { return seedUsers(db, pwdHash) }},
		{"assets", seedAssets},
		{"inspections", seedInspections},
	}

	for _, s := range steps {
		if err := s.fn(db); err != nil {
			log.Fatalf("seed %s: %v", s.name, err)
		}
		fmt.Printf("✓ %s\n", s.name)
	}

	fmt.Println("\nSeed complete. Credentials:")
	fmt.Println("  admin@veylo.demo    / password123  (ADMIN)")
	fmt.Println("  manager@veylo.demo  / password123  (MANAGER)")
	fmt.Println("  inspector@veylo.demo / password123 (INSPECTOR)")
	fmt.Println("  evaluator@veylo.demo / password123 (EVALUATOR)")
}

func seedOrg(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO organizations (id, name, vertical, onboarding_completed_at)
		VALUES ($1, 'Veylo Demo', 'VEHICLE', NOW())
		ON CONFLICT (id) DO NOTHING
	`, orgID)
	return err
}

func seedWorkflow(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO workflows (id, organization_id)
		VALUES ($1, $2)
		ON CONFLICT (id) DO NOTHING
	`, workflowID, orgID)
	if err != nil {
		return err
	}

	statuses := []struct {
		name        string
		description string
		stage       string
		isInitial   bool
	}{
		{"new", "Inspection created", "ENTRY", true},
		{"damage_entered", "Damages recorded", "ENTRY", false},
		{"damage_evaluated", "Damages assessed", "EVALUATION", false},
		{"inspected", "Pending manager review", "REVIEW", false},
		{"completed", "Inspection closed", "FINAL", false},
	}

	for _, s := range statuses {
		_, err := db.Exec(`
			INSERT INTO workflow_statuses (workflow_id, name, description, stage, is_initial)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (workflow_id, name) DO NOTHING
		`, workflowID, s.name, s.description, s.stage, s.isInitial)
		if err != nil {
			return fmt.Errorf("status %s: %w", s.name, err)
		}
	}

	transitions := [][2]string{
		{"new", "damage_entered"},
		{"damage_entered", "damage_evaluated"},
		{"damage_evaluated", "inspected"},
		{"inspected", "completed"},
	}

	for _, t := range transitions {
		_, err := db.Exec(`
			INSERT INTO workflow_transitions (workflow_id, from_status, to_status)
			VALUES ($1, $2, $3)
			ON CONFLICT (workflow_id, from_status, to_status) DO NOTHING
		`, workflowID, t[0], t[1])
		if err != nil {
			return fmt.Errorf("transition %s→%s: %w", t[0], t[1], err)
		}
	}

	return nil
}

func seedUsers(db *sql.DB, pwdHash string) error {
	users := []struct {
		id, email, name, role string
	}{
		{userAdmin, "admin@veylo.demo", "Alice Admin", "ADMIN"},
		{userManager, "manager@veylo.demo", "Bob Manager", "MANAGER"},
		{userInsp, "inspector@veylo.demo", "Carl Inspector", "INSPECTOR"},
		{userEval, "evaluator@veylo.demo", "Dana Evaluator", "EVALUATOR"},
	}

	for _, u := range users {
		_, err := db.Exec(`
			INSERT INTO users (id, organization_id, email, password_hash, full_name, role, status)
			VALUES ($1, $2, $3, $4, $5, $6, 'ACTIVE')
			ON CONFLICT (id) DO NOTHING
		`, u.id, orgID, u.email, pwdHash, u.name, u.role)
		if err != nil {
			return fmt.Errorf("user %s: %w", u.email, err)
		}
	}
	return nil
}

func seedAssets(db *sql.DB) error {
	vehicles := []struct {
		id, vin, plate, brand, model, bodyType, fuelType, color string
		odometer, power                                          int
	}{
		{"seed-ast-001", "WBA5A5C50FD520469", "A001BC", "BMW", "5 Series", "Sedan", "Diesel", "Black", 45000, 190},
		{"seed-ast-002", "WDDGF8AB1EA943456", "B002DE", "Mercedes-Benz", "C 220", "Sedan", "Diesel", "Silver", 62000, 170},
		{"seed-ast-003", "WAUZZZ8K5BA012345", "C003FG", "Audi", "A4", "Sedan", "Petrol", "White", 38000, 150},
		{"seed-ast-004", "4T1BF1FK2EU789012", "D004HI", "Toyota", "Camry", "Sedan", "Hybrid", "Blue", 91000, 160},
		{"seed-ast-005", "WVWZZZ3CZDE345678", "E005JK", "Volkswagen", "Passat", "Wagon", "Diesel", "Grey", 73000, 140},
	}

	for _, v := range vehicles {
		_, err := db.Exec(`
			INSERT INTO assets (id, organization_id, type)
			VALUES ($1, $2, 'vehicle')
			ON CONFLICT (id) DO NOTHING
		`, v.id, orgID)
		if err != nil {
			return fmt.Errorf("asset %s: %w", v.id, err)
		}

		_, err = db.Exec(`
			INSERT INTO vehicle_attributes
				(asset_id, vin, license_plate, brand, model, body_type, fuel_type, color, odometer_reading, engine_power)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
			ON CONFLICT (asset_id) DO NOTHING
		`, v.id, v.vin, v.plate, v.brand, v.model, v.bodyType, v.fuelType, v.color, v.odometer, v.power)
		if err != nil {
			return fmt.Errorf("vehicle_attributes %s: %w", v.id, err)
		}
	}
	return nil
}

func seedInspections(db *sql.DB) error {
	now := time.Now().UTC()

	inspections := []struct {
		id, assetID, status, contractNumber string
		createdAt                           time.Time
	}{
		{"seed-ins-001", "seed-ast-001", "new", "CTR-2025-001", now.Add(-1 * 24 * time.Hour)},
		{"seed-ins-002", "seed-ast-002", "new", "CTR-2025-002", now.Add(-2 * 24 * time.Hour)},
		{"seed-ins-003", "seed-ast-003", "damage_entered", "CTR-2025-003", now.Add(-3 * 24 * time.Hour)},
		{"seed-ins-004", "seed-ast-004", "damage_entered", "CTR-2025-004", now.Add(-4 * 24 * time.Hour)},
		{"seed-ins-005", "seed-ast-001", "damage_evaluated", "CTR-2025-005", now.Add(-5 * 24 * time.Hour)},
		{"seed-ins-006", "seed-ast-005", "damage_evaluated", "CTR-2025-006", now.Add(-6 * 24 * time.Hour)},
		{"seed-ins-007", "seed-ast-002", "inspected", "CTR-2025-007", now.Add(-7 * 24 * time.Hour)},
		{"seed-ins-008", "seed-ast-003", "inspected", "CTR-2025-008", now.Add(-8 * 24 * time.Hour)},
		{"seed-ins-009", "seed-ast-004", "completed", "CTR-2025-009", now.Add(-9 * 24 * time.Hour)},
		{"seed-ins-010", "seed-ast-005", "completed", "CTR-2025-010", now.Add(-10 * 24 * time.Hour)},
	}

	for _, ins := range inspections {
		_, err := db.Exec(`
			INSERT INTO inspections (id, organization_id, asset_id, contract_number, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $6)
			ON CONFLICT (id) DO NOTHING
		`, ins.id, orgID, ins.assetID, ins.contractNumber, ins.status, ins.createdAt)
		if err != nil {
			return fmt.Errorf("inspection %s: %w", ins.id, err)
		}
	}
	return nil
}
