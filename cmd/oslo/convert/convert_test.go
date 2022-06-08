package convert

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConvertCmd(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    []string
		wantOut string
		wantErr bool
	}{
		{
			name: "Single file - Service",
			args: []string{
				"-o", "nobl9",
				"-f", "../../../test/v1/service/service.yaml",
			},
			wantOut: `---
apiVersion: n9/v1alpha
kind: Service
metadata:
    name: my-rad-service
    displayName: My Rad Service
    project: default
spec:
    description: This is a great description of an even better service.
`,
			wantErr: false,
		},
		{
			name: "Single file - Service - non-default project",
			args: []string{
				"-o", "nobl9",
				"-p", "my-project",
				"-f", "../../../test/v1/service/service.yaml",
			},
			wantOut: `---
apiVersion: n9/v1alpha
kind: Service
metadata:
    name: my-rad-service
    displayName: My Rad Service
    project: my-project
spec:
    description: This is a great description of an even better service.
`,
			wantErr: false,
		},
		{
			name: "Single file - Alert Policy",
			args: []string{
				"-o", "nobl9",
				"-f", "../../../test/v1/alert-policy/alert-policy.yaml",
			},
			wantOut: ``,
			wantErr: true,
		},
		{
			name: "Alert Policy - Separate Condition",
			args: []string{
				"-o", "nobl9",
				"-f", "../../../test/v1/alert-policy/alert-policy.yaml",
				"-f", "../../../test/v1/alert-condition/alert-condition.yaml",
			},
			wantOut: `---
apiVersion: n9/v1alpha
kind: AlertPolicy
metadata:
    name: AlertPolicy
    displayName: Alert Policy
    project: default
spec:
    description: Alert policy for cpu usage breaches, notifies on-call devops via email
    severity: Medium
    coolDown: 5m
    conditions:
        - measurement: averageBurnRate
          value: 24
          lastsFor: 3m
          op: gt
    alertMethods: []
`,
			wantErr: false,
		},
		{
			name: "Alert Policy - Inline Condition",
			args: []string{
				"-o", "nobl9",
				"-f", "../../../test/v1/alert-policy/alert-policy-inline-cond.yaml",
			},
			wantOut: `---
apiVersion: n9/v1alpha
kind: AlertPolicy
metadata:
    name: AlertPolicy
    displayName: Alert Policy
    project: default
spec:
    description: Alert policy for cpu usage breaches, notifies on-call devops via email
    severity: Medium
    coolDown: 5m
    conditions:
        - measurement: averageBurnRate
          value: 2
          lastsFor: 5m
          op: gt
    alertMethods: []
`,
			wantErr: false,
		},
		{
			name: "Alert Policy - Multiple Condition",
			args: []string{
				"-o", "nobl9",
				"-f", "../../../test/v1/alert-policy/alert-policy-many-cond.yaml",
				"-f", "../../../test/v1/alert-condition/alert-condition.yaml",
			},
			wantOut: `---
apiVersion: n9/v1alpha
kind: AlertPolicy
metadata:
    name: AlertPolicy
    displayName: Alert Policy
    project: default
spec:
    description: Alert policy for cpu usage breaches, notifies on-call devops via email
    severity: Medium
    coolDown: 5m
    conditions:
        - measurement: averageBurnRate
          value: 24
          lastsFor: 3m
          op: gt
        - measurement: averageBurnRate
          value: 24
          lastsFor: 3m
          op: gt
        - measurement: averageBurnRate
          value: 24
          lastsFor: 3m
          op: gt
        - measurement: averageBurnRate
          value: 24
          lastsFor: 3m
          op: gt
        - measurement: averageBurnRate
          value: 24
          lastsFor: 3m
          op: gt
        - measurement: averageBurnRate
          value: 24
          lastsFor: 3m
          op: gt
        - measurement: averageBurnRate
          value: 24
          lastsFor: 3m
          op: gt
        - measurement: averageBurnRate
          value: 24
          lastsFor: 3m
          op: gt
        - measurement: averageBurnRate
          value: 24
          lastsFor: 3m
          op: gt
    alertMethods: []
`,
			wantErr: false,
		},
		{
			name: "Alert Policy - No Matching Condition",
			args: []string{
				"-o", "nobl9",
				"-f", "../../../test/v1/alert-policy/alert-policy.yaml",
				"-f", "../../../test/v1/alert-condition/alert-condition-invalid-name.yaml",
			},
			wantOut: ``,
			wantErr: true,
		},
		{
			name: "Duplicate file",
			args: []string{
				"-o", "nobl9",
				"-f", "../../../test/v1/service/service.yaml",
				"-f", "../../../test/v1/service/service.yaml",
			},
			wantOut: `---
apiVersion: n9/v1alpha
kind: Service
metadata:
    name: my-rad-service
    displayName: My Rad Service
    project: default
spec:
    description: This is a great description of an even better service.
`,
			wantErr: false,
		},
		{
			name: "Single SLO",
			args: []string{
				"-o", "nobl9",
				"-f", "../../../test/v1/slo/slo-no-indicatorref-rolling-alerts.yaml",
			},
			wantOut: `---
apiVersion: n9/v1alpha
kind: SLO
metadata:
    name: TestSLO
    displayName: Test SLO
    project: default
spec:
    description: This is a great description
    indicator:
        metricSource:
            name: ChangeMe
            kind: ""
    budgetingMethod: Occurrences
    objectives:
        - displayName: Foo Total Errors
          value: 1
          target: 0.98
          countMetrics:
            incremental: true
            good:
                datadog:
                    query: sum:trace.http.request.hits.by_http_status{http.status_code:200}.as_count()
            total:
                datadog:
                    query: sum:trace.http.request.hits.by_http_status{*}.as_count()
    service: TheServiceName
    timeWindows:
        - unit: Month
          count: 1
          isRolling: true
    alertPolicies:
        - FooAlertPolicy
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt // Parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := new(bytes.Buffer)
			root := NewConvertCmd()
			root.SetOut(actual)
			root.SetErr(actual)
			root.SetArgs(tt.args)

			if err := root.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("Error executing root command: %s", err)
				return
			}

			// if tt.wantErr is false assert that the output is correct
			if !tt.wantErr {
				assert.Equal(t, tt.wantOut, actual.String())
			}
		})
	}
}