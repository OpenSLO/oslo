/*
Package v1 contains all the types that are exported by the v1 API.
Copyright © 2022 OpenSLO Team

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package v1

// AlertCondition is a condition that is used to trigger an alert.
type AlertCondition struct {
	ObjectHeader `yaml:",inline"`
	Spec         AlertConditionSpec `yaml:"spec"`
}

// AlertConditionInline is used for inline definitions.  It is slightly
// different from the AlertCondition type because it does not have an APIVersion.
type AlertConditionInline struct {
	Kind     string             `yaml:"kind" validate:"required"`
	Metadata Metadata           `yaml:"metadata" validate:"required"`
	Spec     AlertConditionSpec `yaml:"spec" validate:"required"`
}

// AlertConditionSpec is the specification of an alert condition.
type AlertConditionSpec struct {
	Description string        `yaml:"description,omitempty" validate:"max=1050,omitempty" example:"If the CPU usage is too high for given period then it should alert"` //nolint:lll
	Severity    string        `yaml:"severity" validate:"required" example:"page"`
	Condition   ConditionType `yaml:"condition" validate:"required"`
}

// AlertConditionType is the type of an alert condition.
type AlertConditionType string

const (
	// AlertConditionTypeBurnRate is the type of a burn rate alert condition.
	AlertConditionTypeBurnRate AlertConditionType = "burnrate"
)

// ConditionType is the type of a condition to trigger an alert.
type ConditionType struct {
	Kind           *AlertConditionType `yaml:"kind" validate:"required,oneof=burnrate" example:"burnrate"`
	Threshold      int                 `yaml:"threshold" validate:"required" example:"2"`
	LookbackWindow string              `yaml:"lookbackWindow" validate:"required,validDuration" example:"1h"`
	AlertAfter     string              `yaml:"alertAfter" validate:"required,validDuration" example:"5m"`
}

// AlertNotificationTarget is a target for sending alerts.
type AlertNotificationTarget struct {
	ObjectHeader `yaml:",inline"`
	Spec         AlertNotificationTargetSpec `yaml:"spec"`
}

// AlertNotificationTargetSpec is the specification of an alert notification target.
type AlertNotificationTargetSpec struct {
	Target      string `yaml:"target" validate:"required" example:"slack"`
	Description string `yaml:"description,omitempty" validate:"max=1050,omitempty" example:"Sends P1 alert notifications to the slack channel"` //nolint:lll
}

// AlertPolicy is a policy for sending alerts.
type AlertPolicy struct {
	ObjectHeader `yaml:",inline"`
	Spec         AlertPolicySpec `yaml:"spec"`
}

// AlertPolicyCondition is a condition that is used to trigger an alert in an alert policy.  It can
// be either an inline condition or a reference to an alert condition.
type AlertPolicyCondition struct {
	*AlertPolicyConditionSpec `yaml:",inline,omitempty" validate:"required_without=AlertConditionInline"`
	*AlertConditionInline     `yaml:",inline,omitempty" validate:"required_without=AlertPolicyConditionSpec"`
}

// AlertPolicyConditionSpec is the specification of an alert policy condition.  It is
// used to reference an AlertCondition.
type AlertPolicyConditionSpec struct {
	ConditionRef string `yaml:"conditionRef" validate:"max=1050,required" example:"cpu-usage-breach"`
}

// AlertPolicyNotificationTarget is a reference to an AlertNotificationTarget.
type AlertPolicyNotificationTarget struct {
	TargetRef string `yaml:"targetRef" validate:"required" example:"OnCallDevopsMailNotification"`
}

// AlertPolicySpec is the specification of an alert policy.
type AlertPolicySpec struct {
	Description         string                          `yaml:"description,omitempty" validate:"max=1050,omitempty" example:"Alert policy for cpu usage breaches, notifies on-call devops via email"` //nolint:lll
	AlertWhenNoData     bool                            `yaml:"alertWhenNoData"`
	AlertWhenBreaching  bool                            `yaml:"alertWhenBreaching"`
	AlertWhenResolved   bool                            `yaml:"alertWhenResolved"`
	Conditions          []AlertPolicyCondition          `yaml:"conditions" validate:"required,len=1,dive"`
	NotificationTargets []AlertPolicyNotificationTarget `yaml:"notificationTargets" validate:"required,dive"`
}

// Kind returns the name of this type.
func (AlertNotificationTarget) Kind() string {
	return "AlertNotificationTarget"
}

// Kind returns the name of this type.
func (AlertCondition) Kind() string {
	return "AlertCondition"
}

// Kind returns the name of this type.
func (AlertPolicy) Kind() string {
	return "AlertPolicy"
}
