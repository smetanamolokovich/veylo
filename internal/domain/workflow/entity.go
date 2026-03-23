package workflow

import (
	"errors"
	"fmt"
	"time"
)

type SystemStage string

const (
	StageEntry      SystemStage = "ENTRY"
	StageEvaluation SystemStage = "EVALUATION"
	StageReview     SystemStage = "REVIEW"
	StageFinal      SystemStage = "FINAL"
)

func (s SystemStage) IsValid() bool {
	switch s {
	case StageEntry, StageEvaluation, StageReview, StageFinal:
		return true
	}
	return false
}

type WorkflowStatus struct {
	name        string
	description string
	stage       SystemStage
	isInitial   bool
}

func NewWorkflowStatus(name, description string, stage SystemStage, isInitial bool) (WorkflowStatus, error) {
	if name == "" {
		return WorkflowStatus{}, errors.New("workflow: status name is required")
	}
	if !stage.IsValid() {
		return WorkflowStatus{}, fmt.Errorf("workflow: invalid system stage %q", stage)
	}
	return WorkflowStatus{
		name:        name,
		description: description,
		stage:       stage,
		isInitial:   isInitial,
	}, nil
}

func (s WorkflowStatus) Name() string        { return s.name }
func (s WorkflowStatus) Description() string { return s.description }
func (s WorkflowStatus) Stage() SystemStage  { return s.stage }
func (s WorkflowStatus) IsInitial() bool     { return s.isInitial }

type WorkflowTransition struct {
	fromStatus string
	toStatus   string
}

func NewWorkflowTransition(from, to string) (WorkflowTransition, error) {
	if from == "" || to == "" {
		return WorkflowTransition{}, errors.New("workflow: from and to statuses are required")
	}
	if from == to {
		return WorkflowTransition{}, errors.New("workflow: from and to statuses must differ")
	}
	return WorkflowTransition{fromStatus: from, toStatus: to}, nil
}

func (t WorkflowTransition) From() string { return t.fromStatus }
func (t WorkflowTransition) To() string   { return t.toStatus }

type Workflow struct {
	id             string
	organizationID string
	statuses       []WorkflowStatus
	transitions    []WorkflowTransition
	createdAt      time.Time
	updatedAt      time.Time
}

func NewWorkflow(id, organizationID string) (*Workflow, error) {
	if id == "" || organizationID == "" {
		return nil, errors.New("workflow: id and organizationID are required")
	}
	now := time.Now().UTC()
	return &Workflow{
		id:             id,
		organizationID: organizationID,
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

func ReconstitueWorkflow(
	id, organizationID string,
	statuses []WorkflowStatus,
	transitions []WorkflowTransition,
	createdAt, updatedAt time.Time,
) *Workflow {
	return &Workflow{
		id:             id,
		organizationID: organizationID,
		statuses:       statuses,
		transitions:    transitions,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

func (w *Workflow) ID() string                        { return w.id }
func (w *Workflow) OrganizationID() string            { return w.organizationID }
func (w *Workflow) Statuses() []WorkflowStatus        { return w.statuses }
func (w *Workflow) Transitions() []WorkflowTransition { return w.transitions }
func (w *Workflow) CreatedAt() time.Time              { return w.createdAt }
func (w *Workflow) UpdatedAt() time.Time              { return w.updatedAt }

func (w *Workflow) AddStatus(s WorkflowStatus) error {
	for _, existing := range w.statuses {
		if existing.name == s.name {
			return fmt.Errorf("%w: %s", ErrDuplicateStatus, s.name)
		}
	}
	if s.isInitial {
		for _, existing := range w.statuses {
			if existing.isInitial {
				return ErrInitialStatusAlreadySet
			}
		}
	}
	w.statuses = append(w.statuses, s)
	w.updatedAt = time.Now().UTC()
	return nil
}

func (w *Workflow) AddTransition(t WorkflowTransition) error {
	if !w.hasStatus(t.fromStatus) {
		return fmt.Errorf("%w: %s", ErrStatusNotFound, t.fromStatus)
	}
	if !w.hasStatus(t.toStatus) {
		return fmt.Errorf("%w: %s", ErrStatusNotFound, t.toStatus)
	}
	for _, existing := range w.transitions {
		if existing.fromStatus == t.fromStatus && existing.toStatus == t.toStatus {
			return fmt.Errorf("workflow: transition %s→%s already exists", t.fromStatus, t.toStatus)
		}
	}
	w.transitions = append(w.transitions, t)
	w.updatedAt = time.Now().UTC()
	return nil
}

func (w *Workflow) InitialStatus() (string, error) {
	for _, s := range w.statuses {
		if s.isInitial {
			return s.name, nil
		}
	}
	return "", ErrNoInitialStatus
}

// AllowedTransitions returns a map of status → allowed next statuses.
// Used by the inspection domain to validate transitions without importing this package.
func (w *Workflow) AllowedTransitions() map[string][]string {
	result := make(map[string][]string)
	for _, t := range w.transitions {
		result[t.fromStatus] = append(result[t.fromStatus], t.toStatus)
	}
	return result
}

func (w *Workflow) StageOf(statusName string) (SystemStage, error) {
	for _, s := range w.statuses {
		if s.name == statusName {
			return s.stage, nil
		}
	}
	return "", fmt.Errorf("%w: %s", ErrStatusNotFound, statusName)
}

func (w *Workflow) hasStatus(name string) bool {
	for _, s := range w.statuses {
		if s.name == name {
			return true
		}
	}
	return false
}
