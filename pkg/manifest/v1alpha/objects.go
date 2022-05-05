package v1alpha

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/OpenSLO/oslo/pkg/manifest"
)

// APIVersion is a value of valid apiVersions.
const (
	APIVersion = "openslo/v1alpha"
)

// Possible values of field kind for valid Objects.
const (
	KindSLO     = "SLO"
	KindService = "Service"
)

// ObjectHeader is a header for all objects.
type ObjectHeader struct {
	manifest.ObjectHeader `yaml:",inline"`
	Kind                  string `yaml:"kind" validate:"required,oneof=Service SLO AlertNotificationTarget" example:"kind"`
	MetadataHolder        `yaml:",inline"`
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

// SLO struct which mapped one to one with kind: slo yaml definition, external usage.
type SLO struct {
	ObjectHeader `yaml:",inline"`
	Spec         SLOSpec `yaml:"spec"`
}

// Kind returns the name of this type.
func (SLO) Kind() string {
	return "SLO"
}

// SLOSpec represents content of Spec typical for SLO Object.
type SLOSpec struct {
	TimeWindows     []TimeWindow `yaml:"timeWindows" validate:"required,len=1,dive"`
	BudgetingMethod string       `yaml:"budgetingMethod" validate:"required,oneof=Occurrences Timeslices" example:"Occurrences"` //nolint: lll
	Description     string       `yaml:"description" validate:"max=1050" example:"Total count of server requests"`
	Indicator       *Indicator   `yaml:"indicator"`
	Service         string       `yaml:"service" validate:"required" example:"webapp-service"`
	Objectives      []Objective  `json:"objectives" validate:"required,dive"`
}

// Indicator represents integration with metric source.
type Indicator struct {
	ThresholdMetric MetricSourceSpec `yaml:"thresholdMetric" validate:"required"`
}

// MetricSourceSpec represents the metric source.
type MetricSourceSpec struct {
	Source    string `yaml:"source" validate:"required,alpha"`
	QueryType string `yaml:"queryType" validate:"required,alpha"`
	Query     string `yaml:"query" validate:"required"`
}

// Objective represents single threshold for SLO, for internal usage.
type Objective struct {
	ObjectiveBase   `yaml:",inline"`
	RatioMetrics    *RatioMetrics `yaml:"ratioMetrics"`
	BudgetTarget    *float64      `yaml:"target" validate:"required,numeric,gte=0,lt=1" example:"0.9"`
	TimeSliceTarget *float64      `yaml:"timeSliceTarget,omitempty" example:"0.9"`
	Operator        *string       `yaml:"op,omitempty" example:"lte"`
}

// RatioMetrics base struct for ratio metrics.
type RatioMetrics struct {
	Good    MetricSourceSpec `yaml:"good" validate:"required"`
	Total   MetricSourceSpec `yaml:"total" validate:"required"`
	Counter bool             `yaml:"counter" example:"true"`
}

// ObjectiveBase base structure representing a threshold.
type ObjectiveBase struct {
	DisplayName string  `yaml:"displayName" validate:"max=1050" example:"Good"`
	Value       float64 `yaml:"value" validate:"numeric" example:"100"`
}

// TimeWindow represents content of time window.
type TimeWindow struct {
	Unit      string    `yaml:"unit" validate:"required,oneof=Second Quarter Month Week Day" example:"Week"`
	Count     int       `yaml:"count" validate:"required,gt=0" example:"1"`
	IsRolling bool      `yaml:"isRolling" example:"true"`
	Calendar  *Calendar `yaml:"calendar,omitempty"`
}

// Calendar struct represents calendar time window.
type Calendar struct {
	StartTime string `yaml:"startTime" validate:"required,dateWithTime" example:"2020-01-21 12:30:00"`
	TimeZone  string `yaml:"timeZone" validate:"required,timeZone" example:"America/New_York"`
}

// Metadata represents part of object which is is common for all available Objects, for internal usage.
type Metadata struct {
	Name        string `yaml:"name" validate:"required" example:"name"`
	DisplayName string `yaml:"displayName,omitempty" validate:"omitempty,min=0,max=63" example:"Prometheus Source"`
}

// MetadataHolder is an intermediate structure that can provides metadata related
// field to other structures.
type MetadataHolder struct {
	Metadata Metadata `yaml:"metadata"`
}

// ObjectGeneric represents struct to which every Objects is parsable
// Specific types of Object have different structures as Spec.
type ObjectGeneric struct {
	ObjectHeader `yaml:",inline"`
}

// Parse is responsible for parsing all structs in this apiVersion.
func Parse(fileContent []byte, m ObjectGeneric, filename string) (manifest.OpenSLOKind, error) {
	switch m.Kind {
	case KindService:
		var content Service
		err := yaml.Unmarshal(fileContent, &content)
		return content, err
	case KindSLO:
		var content SLO
		err := yaml.Unmarshal(fileContent, &content)
		return content, err
	default:
		return nil, fmt.Errorf("unsupported kind: %s", m.Kind)
	}
}
