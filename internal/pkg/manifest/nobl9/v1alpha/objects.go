package nobl9v1alpha

import (
	nobl9manifest "github.com/OpenSLO/oslo/internal/pkg/manifest/nobl9"
)

// apiVersion: n9/v1alpha
// kind: Service
// metadata:
//   name: # string
//   displayName: # string
//   project: # string
// spec:
//   description: # string
//   serviceType: # string
const (
	APIVersion = "n9/v1alpha"
)

// Possible values of field kind for valid Objects.
const (
	KindSLO     = "SLO"
	KindService = "Service"
)

// Metadata represents part of object which is is common for all available Objects, for internal usage.
type Metadata struct {
	Name        string `yaml:"name" validate:"required" example:"name"`
	DisplayName string `yaml:"displayName,omitempty" validate:"omitempty,min=0,max=63" example:"Prometheus Source"`
	Project     string `json:"project,omitempty" validate:"objectName" example:"default"`
	// TODO Come back to Labels
}

// MetadataHolder is an intermediate structure that can provides metadata related
// field to other structures.
type MetadataHolder struct {
	Metadata Metadata `yaml:"metadata"`
}

// ObjectHeader is a header for all objects.
type ObjectHeader struct {
	nobl9manifest.ObjectHeader `yaml:",inline"`
	Kind                       string `yaml:"kind" validate:"required,oneof=Service SLO AlertNotificationTarget" example:"kind"` //nolint:lll
	MetadataHolder             `yaml:",inline"`
}

// Service struct which mapped one to one with kind: service yaml definition.
type Service struct {
	ObjectHeader `yaml:",inline"`
	Spec         ServiceSpec `yaml:"spec"`
}

// Kind returns the name of this type.
func (Service) Kind() string {
	return "Service"
}

// ServiceSpec represents content of Spec typical for Service Object.
type ServiceSpec struct {
	Description string `yaml:"description" validate:"max=1050" example:"Bleeding edge web app"`
}

// AlertPolicy represents a set of conditions that can trigger an alert.
type AlertPolicy struct {
	ObjectHeader `yaml:",inline"`
	Spec         AlertPolicySpec `yaml:"spec"`
}

// AlertPolicySpec represents content of AlertPolicy's Spec.
type AlertPolicySpec struct {
	Description      string           `yaml:"description" validate:"description" example:"Error budget is at risk"`
	Severity         string           `yaml:"severity" validate:"required,severity" example:"High"`
	CoolDownDuration string           `yaml:"coolDown,omitempty" validate:"omitempty,validDuration,nonNegativeDuration,durationAtLeast=5m" example:"5m"` //nolint:lll
	Conditions       []AlertCondition `yaml:"conditions" validate:"required,min=1,dive"`
	AlertMethods     []string         `yaml:"alertMethods"`
}

// AlertCondition represents a condition to meet to trigger an alert.
type AlertCondition struct {
	Measurement      string      `yaml:"measurement" validate:"required,alertPolicyMeasurement" example:"BurnedBudget"`
	Value            interface{} `yaml:"value" validate:"required" swaggertype:"string" example:"0.97"`
	LastsForDuration string      `yaml:"lastsFor,omitempty" validate:"omitempty,validDuration,nonNegativeDuration" example:"15m"` //nolint:lll
	Operation        string      `yaml:"op" validate:"required,alertOperation" example:"lt"`
}

// AlertMethodAssignment represents an AlertMethod assigned to AlertPolicy.
type AlertMethodAssignment struct {
	Project string `json:"project,omitempty" validate:"omitempty,objectName" example:"default"`
	Name    string `json:"name" validate:"required,objectName" example:"webhook-alertmethod"`
}
