package e2e_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInspectionFullFlow covers the complete vehicle inspection lifecycle:
// signup → create vehicle → create inspection → add finding → assess → transitions → get report
func TestInspectionFullFlow(t *testing.T) {
	ts, cleanup := newTestServer(t)
	defer cleanup()

	// ── 1. Signup ──────────────────────────────────────────────────────────
	t.Run("signup creates org, workflow and admin user", func(t *testing.T) {
		var resp map[string]any
		r := ts.do(t, http.MethodPost, "/api/auth/signup", jsonBody(
			"org_name", "Test Leasing",
			"vertical", "VEHICLE",
			"email", "admin@test.com",
			"password", "secret123",
			"full_name", "Admin User",
		), &resp)

		require.Equal(t, http.StatusCreated, r.StatusCode)
		assert.NotEmpty(t, resp["access_token"])
		assert.NotEmpty(t, resp["organization_id"])
		assert.NotEmpty(t, resp["user_id"])

		ts.token = resp["access_token"].(string)
	})

	// ── 2. Get organization ─────────────────────────────────────────────────
	t.Run("get organization me", func(t *testing.T) {
		var resp map[string]any
		r := ts.do(t, http.MethodGet, "/api/v1/organizations/me", nil, &resp)

		require.Equal(t, http.StatusOK, r.StatusCode)
		assert.Equal(t, "Test Leasing", resp["name"])
		assert.Equal(t, "VEHICLE", resp["vertical"])
	})

	// ── 3. Get workflow ─────────────────────────────────────────────────────
	t.Run("default vehicle workflow exists after signup", func(t *testing.T) {
		var resp map[string]any
		r := ts.do(t, http.MethodGet, "/api/v1/workflow", nil, &resp)

		require.Equal(t, http.StatusOK, r.StatusCode)
		statuses := resp["statuses"].([]any)
		assert.Len(t, statuses, 5)
		transitions := resp["transitions"].([]any)
		assert.Len(t, transitions, 4)
	})

	// ── 4. Create vehicle ───────────────────────────────────────────────────
	var vehicleID string
	t.Run("create vehicle asset", func(t *testing.T) {
		var resp map[string]any
		r := ts.do(t, http.MethodPost, "/api/v1/assets/vehicles", jsonBody(
			"vin", "WBA12345678901234",
			"license_plate", "AB-123-CD",
			"brand", "BMW",
			"model", "3 Series",
			"body_type", "Sedan",
			"fuel_type", "Petrol",
			"transmission", "Automatic",
			"odometer_reading", 45000,
			"color", "Black",
			"engine_power", 184,
		), &resp)

		require.Equal(t, http.StatusCreated, r.StatusCode)
		require.NotEmpty(t, resp["id"])
		vehicleID = resp["id"].(string)
	})

	// ── 5. Create inspection ────────────────────────────────────────────────
	var inspectionID string
	t.Run("create inspection starts with initial workflow status", func(t *testing.T) {
		var resp map[string]any
		r := ts.do(t, http.MethodPost, "/api/v1/inspections", jsonBody(
			"asset_id", vehicleID,
			"contract_number", "CONTRACT-001",
		), &resp)

		require.Equal(t, http.StatusCreated, r.StatusCode)
		require.NotEmpty(t, resp["id"])
		assert.Equal(t, "new", resp["status"])
		inspectionID = resp["id"].(string)
	})

	// ── 6. Get inspection ───────────────────────────────────────────────────
	t.Run("get inspection by id", func(t *testing.T) {
		var resp map[string]any
		r := ts.do(t, http.MethodGet, "/api/v1/inspections/"+inspectionID, nil, &resp)

		require.Equal(t, http.StatusOK, r.StatusCode)
		assert.Equal(t, inspectionID, resp["id"])
		assert.Equal(t, "CONTRACT-001", resp["contract_number"])
	})

	// ── 7. List inspections ─────────────────────────────────────────────────
	t.Run("list inspections returns the created one", func(t *testing.T) {
		var resp map[string]any
		r := ts.do(t, http.MethodGet, "/api/v1/inspections?page=1&page_size=20", nil, &resp)

		require.Equal(t, http.StatusOK, r.StatusCode)
		assert.Equal(t, float64(1), resp["total"])
		items := resp["items"].([]any)
		assert.Len(t, items, 1)
	})

	// ── 8. Add finding ──────────────────────────────────────────────────────
	var findingID string
	t.Run("add finding to inspection", func(t *testing.T) {
		var resp map[string]any
		r := ts.do(t, http.MethodPost, "/api/v1/inspections/"+inspectionID+"/findings", jsonBody(
			"finding_type", "scratch",
			"description", "Deep scratch on front bumper",
			"location", map[string]any{
				"body_area":    "Front bumper",
				"coordinate_x": 0.25,
				"coordinate_y": 0.10,
			},
		), &resp)

		require.Equal(t, http.StatusCreated, r.StatusCode)
		require.NotEmpty(t, resp["id"])
		findingID = resp["id"].(string)
	})

	// ── 9. List findings ────────────────────────────────────────────────────
	t.Run("list findings for inspection", func(t *testing.T) {
		var resp map[string]any
		r := ts.do(t, http.MethodGet, "/api/v1/inspections/"+inspectionID+"/findings", nil, &resp)

		require.Equal(t, http.StatusOK, r.StatusCode)
		items := resp["items"].([]any)
		assert.Len(t, items, 1)
	})

	// ── 10. Assess finding ──────────────────────────────────────────────────
	t.Run("assess finding with severity and repair method", func(t *testing.T) {
		var resp map[string]any
		r := ts.do(t, http.MethodPut,
			"/api/v1/inspections/"+inspectionID+"/findings/"+findingID+"/assessment",
			jsonBody(
				"severity", "NOT_ACCEPTED",
				"repair_method", "REPAIR",
				"cost_breakdown", map[string]any{
					"parts": 15000,
					"labor": 8000,
					"paint": 5000,
					"other": 0,
				},
			), &resp)

		require.Equal(t, http.StatusOK, r.StatusCode)
		assert.Equal(t, "NOT_ACCEPTED", resp["severity"])
		assert.Equal(t, float64(28000), resp["total_cost"])
	})

	// ── 11. Transitions ─────────────────────────────────────────────────────
	steps := []string{"damage_entered", "damage_evaluated", "inspected", "completed"}
	for _, status := range steps {
		status := status
		t.Run("transition to "+status, func(t *testing.T) {
			var resp map[string]any
			r := ts.do(t, http.MethodPost, "/api/v1/inspections/"+inspectionID+"/transitions",
				jsonBody("status", status), &resp)

			require.Equal(t, http.StatusOK, r.StatusCode)
			assert.Equal(t, status, resp["status"])
		})
	}

	// ── 12. Invalid transition rejected ────────────────────────────────────
	t.Run("invalid transition from completed is rejected", func(t *testing.T) {
		r := ts.do(t, http.MethodPost, "/api/v1/inspections/"+inspectionID+"/transitions",
			jsonBody("status", "new"), nil)

		assert.Equal(t, http.StatusBadRequest, r.StatusCode)
		r.Body.Close()
	})

	// ── 13. Report not found (no S3) ────────────────────────────────────────
	t.Run("report endpoint returns 404 when S3 not configured", func(t *testing.T) {
		r := ts.do(t, http.MethodGet, "/api/v1/inspections/"+inspectionID+"/report", nil, nil)

		assert.Equal(t, http.StatusNotFound, r.StatusCode)
		r.Body.Close()
	})
}

// TestAuthFlow covers login, duplicate signup, and protected route access.
func TestAuthFlow(t *testing.T) {
	ts, cleanup := newTestServer(t)
	defer cleanup()

	// Signup
	var signupResp map[string]any
	r := ts.do(t, http.MethodPost, "/api/auth/signup", jsonBody(
		"org_name", "Auth Test Org",
		"vertical", "VEHICLE",
		"email", "user@auth.com",
		"password", "pass123",
		"full_name", "Test User",
	), &signupResp)
	require.Equal(t, http.StatusCreated, r.StatusCode)
	ts.token = signupResp["access_token"].(string)

	t.Run("protected route requires token", func(t *testing.T) {
		savedToken := ts.token
		ts.token = ""
		r := ts.do(t, http.MethodGet, "/api/v1/organizations/me", nil, nil)
		assert.Equal(t, http.StatusUnauthorized, r.StatusCode)
		r.Body.Close()
		ts.token = savedToken
	})

	t.Run("login returns tokens", func(t *testing.T) {
		var resp map[string]any
		r := ts.do(t, http.MethodPost, "/api/auth/login", jsonBody(
			"email", "user@auth.com",
			"password", "pass123",
			"organization_id", signupResp["organization_id"],
		), &resp)

		require.Equal(t, http.StatusOK, r.StatusCode)
		assert.NotEmpty(t, resp["access_token"])
		assert.NotEmpty(t, resp["refresh_token"])
	})

	t.Run("login with wrong password returns 401", func(t *testing.T) {
		r := ts.do(t, http.MethodPost, "/api/auth/login", jsonBody(
			"email", "user@auth.com",
			"password", "wrong",
			"organization_id", signupResp["organization_id"],
		), nil)

		assert.Equal(t, http.StatusUnauthorized, r.StatusCode)
		r.Body.Close()
	})
}

// TestWorkflowCustomization verifies orgs can add custom statuses and transitions.
func TestWorkflowCustomization(t *testing.T) {
	ts, cleanup := newTestServer(t)
	defer cleanup()

	var signupResp map[string]any
	r := ts.do(t, http.MethodPost, "/api/auth/signup", jsonBody(
		"org_name", "Custom Workflow Org",
		"vertical", "VEHICLE",
		"email", "admin@custom.com",
		"password", "pass123",
		"full_name", "Admin",
	), &signupResp)
	require.Equal(t, http.StatusCreated, r.StatusCode)
	ts.token = signupResp["access_token"].(string)

	t.Run("add custom status to workflow", func(t *testing.T) {
		var resp map[string]any
		r := ts.do(t, http.MethodPost, "/api/v1/workflow/statuses", jsonBody(
			"name", "photos_taken",
			"description", "Photos of damages uploaded",
			"stage", "ENTRY",
			"is_initial", false,
		), &resp)

		require.Equal(t, http.StatusCreated, r.StatusCode)
		assert.Equal(t, "photos_taken", resp["name"])
		assert.Equal(t, "ENTRY", resp["stage"])
	})

	t.Run("add transition for custom status", func(t *testing.T) {
		var resp map[string]any
		r := ts.do(t, http.MethodPost, "/api/v1/workflow/transitions", jsonBody(
			"from_status", "damage_entered",
			"to_status", "photos_taken",
		), &resp)

		require.Equal(t, http.StatusCreated, r.StatusCode)
		assert.Equal(t, "damage_entered", resp["from_status"])
		assert.Equal(t, "photos_taken", resp["to_status"])
	})

	t.Run("workflow now has 6 statuses and 5 transitions", func(t *testing.T) {
		var resp map[string]any
		r := ts.do(t, http.MethodGet, "/api/v1/workflow", nil, &resp)

		require.Equal(t, http.StatusOK, r.StatusCode)
		assert.Len(t, resp["statuses"].([]any), 6)
		assert.Len(t, resp["transitions"].([]any), 5)
	})
}
