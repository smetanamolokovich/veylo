package workflow

import "time"

// DefaultVehicleWorkflow returns the standard vehicle inspection workflow.
// Used when a new organization with the VEHICLE vertical is created.
func DefaultVehicleWorkflow(id, organizationID string) *Workflow {
	now := time.Now().UTC()
	statuses := []WorkflowStatus{
		{name: "new", description: "Inspection created", stage: StageEntry, isInitial: true},
		{name: "damage_entered", description: "Damages recorded", stage: StageEntry},
		{name: "damage_evaluated", description: "Damages assessed", stage: StageEvaluation},
		{name: "inspected", description: "Pending manager review", stage: StageReview},
		{name: "completed", description: "Inspection closed", stage: StageFinal},
	}
	transitions := []WorkflowTransition{
		{fromStatus: "new", toStatus: "damage_entered"},
		{fromStatus: "damage_entered", toStatus: "damage_evaluated"},
		{fromStatus: "damage_evaluated", toStatus: "inspected"},
		{fromStatus: "inspected", toStatus: "completed"},
	}
	return &Workflow{
		id:             id,
		organizationID: organizationID,
		statuses:       statuses,
		transitions:    transitions,
		createdAt:      now,
		updatedAt:      now,
	}
}
