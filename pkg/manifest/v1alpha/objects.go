package v1alpha

import (
	"fmt"

	"github.com/openslo/oslo/pkg/manifest"

	"gopkg.in/yaml.v3"
)

// APIVersion is a value of valid apiVersions
const (
	APIVersion = "openslo/v1alpha"
)

// Possible values of field kind for valid Objects.
const (
	KindSLO     = "SLO"
	KindService = "Service"
)

// Service struct which mapped one to one with kind: service yaml definition
type Service struct {
	manifest.ObjectHeader `yaml:",inline"`
	Spec                  ServiceSpec `yaml:"spec"`
}

// ServiceSpec represents content of Spec typical for Service Object
type ServiceSpec struct {
	Description string `yaml:"description" validate:"max=1050" example:"Bleeding edge web app"`
}

// SLO struct which mapped one to one with kind: slo yaml definition, external usage
type SLO struct {
	manifest.ObjectHeader `yaml:",inline"`
	Spec                  SLOSpec `yaml:"spec"`
}

// SLOSpec represents content of Spec typical for SLO Object
type SLOSpec struct {
	// TimeWindows     []TimeWindow `yaml:"timeWindows" validate:"required,len=1,dive"`
	BudgetingMethod string      `yaml:"budgetingMethod" validate:"required,oneof=Occurrences Timeslices" example:"Occurrences"`
	Description     string      `yaml:"description" validate:"max=1050" example:"Total count of server requests"`
	Indicator       *Indicator  `yaml:"indicator"`
	Service         string      `yaml:"service" validate:"required" example:"webapp-service"`
	Objectives      []Objective `json:"objectives" validate:"required,dive"`
}

// Indicator represents integration with metric source
type Indicator struct {
	ThresholdMetric MetricSourceSpec `yaml:"thresholdMetric" validate:"required"`
}

// MetricSourceSpec represents the metric source
type MetricSourceSpec struct {
	Source    string `yaml:"source" validate:"required,alpha"`
	QueryType string `yaml:"queryType" validate:"required,alpha"`
	Query     string `yaml:"query" validate:"required"`
}

// Objective represents single threshold for SLO, for internal usage
type Objective struct {
	ObjectiveBase
	BudgetTarget    *float64 `yaml:"target" validate:"required,numeric,gte=0,lt=1" example:"0.9"`
	TimeSliceTarget *float64 `yaml:"timeSliceTarget,omitempty" example:"0.9"`
	Operator        *string  `yaml:"op,omitempty" example:"lte"`
}

// ObjectiveBase base structure representing a threshold
type ObjectiveBase struct {
	DisplayName string  `yaml:"displayName" validate:"max=1050" example:"Good"`
	Value       float64 `yaml:"value" validate:"numeric" example:"100"`
}

// TimeWindow represents content of time window
type TimeWindow struct {
	Unit      string    `yaml:"unit" validate:"required,oneof=Second Quarter Month Week Day" example:"Week"`
	Count     int       `yaml:"count" validate:"required,gt=0" example:"1"`
	IsRolling bool      `yaml:"isRolling" example:"true"`
	Calendar  *Calendar `yaml:"calendar,omitempty"`
}

// Calendar struct represents calendar time window
type Calendar struct {
	StartTime string `yaml:"startTime" validate:"required,dateWithTime,minDateTime" example:"2020-01-21 12:30:00"`
	TimeZone  string `yaml:"timeZone" validate:"required,timeZone" example:"America/New_York"`
}

// Parse is responsible for parsing all structs in this apiVersion
func Parse(fileContent []byte, m manifest.ObjectGeneric, filename string) (interface{}, error) {
	var allErrors []string

	switch m.Kind {
	case KindService:
		var content Service
		if err := yaml.Unmarshal(fileContent, &content); err != nil {
			allErrors = append(allErrors, err.Error())
		}
		return content, nil
	case KindSLO:
		var content SLO
		if err := yaml.Unmarshal(fileContent, &content); err != nil {
			allErrors = append(allErrors, err.Error())
		}
		return content, nil
	default:
		return nil, fmt.Errorf("Unsupported kind: %s", m.Kind)
	}
}
