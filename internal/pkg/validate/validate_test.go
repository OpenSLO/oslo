/*
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
package validate

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func Test_validateFiles(t *testing.T) {
	t.Parallel()
	type args struct {
		files []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid apiVersion",
			args: args{
				files: []string{"../../../test/invalid-apiversion.yaml"},
			},
			wantErr: true,
		},
		{
			name: "v1alpha gets v1 Kind",
			args: args{
				files: []string{"../../../test/v1/invalid-apiversion.yaml"},
			},
			wantErr: true,
		},
		{
			name: "v1alpha",
			args: args{
				files: []string{
					"../../../test/v1alpha/valid-service.yaml",
					"../../../test/v1alpha/valid-slos-ratio.yaml",
					"../../../test/v1alpha/valid-slos-threshold.yaml",
				},
			},
			wantErr: false,
		},
		{
			name: "v1alpha single file invalid",
			args: args{
				files: []string{"../../../test/v1alpha/invalid-service.yaml"},
			},
			wantErr: true,
		},
		{
			name: "v1 AlertCondition",
			args: args{
				files: []string{
					"../../../test/v1/alert-condition/alert-condition.yaml",
					"../../../test/v1/alert-condition/alert-condition-no-description.yaml",
				},
			},
			wantErr: false,
		},
		{
			name: "v1 AlertCondition invalid",
			args: args{
				files: []string{
					"../../../test/v1/alert-condition/alert-condition-no-condition.yaml",
					"../../../test/v1/alert-condition/alert-condition-no-sev.yaml",
				},
			},
			wantErr: true,
		},
		{
			name: "v1 AlertNotificationTarget",
			args: args{
				files: []string{
					"../../../test/v1/alert-notification-target/alert-notification-target.yaml",
					"../../../test/v1/alert-notification-target/alert-notification-target-no-description.yaml",
				},
			},
			wantErr: false,
		},
		{
			name: "v1 AlertNotificationTarget invalid",
			args: args{
				files: []string{"../../../test/v1/alert-notification-target/alert-notification-target-no-target.yaml"},
			},
			wantErr: true,
		},
		{
			name: "v1 AlertPolicy",
			args: args{
				files: []string{
					"../../../test/v1/alert-policy/alert-policy.yaml",
					"../../../test/v1/alert-policy/alert-policy-inline-cond.yaml",
					"../../../test/v1/alert-policy/alert-policy-many-notificationref.yaml",
				},
			},
			wantErr: false,
		},
		{
			name: "v1 AlertPolicy invalid",
			args: args{
				files: []string{
					"../../../test/v1/alert-policy/alert-policy-malformed-cond.yaml",
					"../../../test/v1/alert-policy/alert-policy-malformed-targetref.yaml",
					"../../../test/v1/alert-policy/alert-policy-many-cond.yaml",
					"../../../test/v1/alert-policy/alert-policy-no-cond.yaml",
					"../../../test/v1/alert-policy/alert-policy-no-notification.yaml",
				},
			},
			wantErr: true,
		},
		{
			name: "v1 single DataSource valid",
			args: args{
				files: []string{"../../../test/v1/data-source/data-source.yaml"},
			},
			wantErr: false,
		},
		{
			name: "v1 Service",
			args: args{
				files: []string{
					"../../../test/v1/service/service.yaml",
					"../../../test/v1/service/service-no-displayname.yaml",
				},
			},
			wantErr: false,
		},
		{
			name: "v1 Service long description",
			args: args{
				files: []string{"../../../test/v1/service/service-long-description.yaml"},
			},
			wantErr: true,
		},
		{
			name: "v1 SLI",
			args: args{
				files: []string{
					"../../../test/v1/sli/sli-description-ratio-bad-inline-metricsource.yaml",
					"../../../test/v1/sli/sli-description-ratio-bad-metricsourceref.yaml",
					"../../../test/v1/sli/sli-description-ratio-good-inline-metricsource.yaml",
					"../../../test/v1/sli/sli-description-ratio-good-metricsourceref.yaml",
					"../../../test/v1/sli/sli-description-threshold-inline-metricsource.yaml",
					"../../../test/v1/sli/sli-description-threshold-metricsourceref.yaml",
					"../../../test/v1/sli/sli-no-description-ratio-bad-inline-metricsource.yaml",
					"../../../test/v1/sli/sli-no-description-ratio-bad-metricsourceref.yaml",
					"../../../test/v1/sli/sli-no-description-ratio-good-inline-metricsource.yaml",
					"../../../test/v1/sli/sli-no-description-ratio-good-metricsourceref.yaml",
					"../../../test/v1/sli/sli-no-description-threshold-inline-metricsource.yaml",
					"../../../test/v1/sli/sli-no-description-threshold-metricsourceref.yaml",
				},
			},
			wantErr: false,
		},
		{
			name: "v1 SLO",
			args: args{
				files: []string{
					"../../../test/v1/slo/slo-indicatorref-calendar-alerts.yaml",
					"../../../test/v1/slo/slo-indicatorref-calendar-no-alerts.yaml",
					"../../../test/v1/slo/slo-indicatorref-rolling-alerts.yaml",
					"../../../test/v1/slo/slo-indicatorref-rolling-no-alerts.yaml",
					"../../../test/v1/slo/slo-no-indicatorref-calendar-alerts.yaml",
					"../../../test/v1/slo/slo-no-indicatorref-calendar-no-alerts.yaml",
					"../../../test/v1/slo/slo-no-indicatorref-rolling-alerts.yaml",
					"../../../test/v1/slo/slo-no-indicatorref-rolling-no-alerts.yaml",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := Files(tt.args.files); (err != nil) != tt.wantErr {
				t.Errorf("Files() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isValidDurationString(t *testing.T) {
	t.Parallel()
	type args struct {
		durStr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Seconds",
			args: args{
				durStr: "1s",
			},
			wantErr: false,
		},
		{
			name: "Minutes",
			args: args{
				durStr: "5m",
			},
			wantErr: false,
		},
		{
			name: "Hours",
			args: args{
				durStr: "2h",
			},
			wantErr: false,
		},
		{
			name: "Weeks",
			args: args{
				durStr: "3w",
			},
			wantErr: false,
		},
		{
			name: "Months",
			args: args{
				durStr: "6M",
			},
			wantErr: false,
		},
		{
			name: "Quarters",
			args: args{
				durStr: "7Q",
			},
			wantErr: false,
		},
		{
			name: "Years",
			args: args{
				durStr: "8Y",
			},
			wantErr: false,
		},
		{
			name: "Invalid",
			args: args{
				durStr: "8y",
			},
			wantErr: true,
		},
	}

	validate := validator.New()
	if err := validate.RegisterValidation("isValidDurationString", isValidDurationString); err != nil {
		t.Errorf("unexpected error registering validation: %v", err)
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := validate.Var(tt.args.durStr, "isValidDurationString"); (err != nil) != tt.wantErr {
				t.Errorf("isValidDurationString() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
