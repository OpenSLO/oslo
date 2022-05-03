/*
Copyright © 2021 OpenSLO Team

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

	"github.com/stretchr/testify/assert"
)

func Test_readConf(t *testing.T) {
	t.Parallel()

	c, e := ReadConf("../../../test/v1alpha/valid-service.yaml")

	assert.NotNil(t, c)
	assert.Nil(t, e)

	_, e = ReadConf("../../../test/v1alpha/non-existent.yaml")

	assert.NotNil(t, e)
}

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
			name: "v1alpha multifile valid",
			args: args{
				files: []string{"../../../test/v1alpha/valid-service.yaml", "../../../test/v1alpha/valid-slos-ratio.yaml"},
			},
			wantErr: false,
		},
		{
			name: "v1alpha single service file valid",
			args: args{
				files: []string{"../../../test/v1alpha/valid-service.yaml"},
			},
			wantErr: false,
		},
		{
			name: "v1alpha single ratio SLO file valid",
			args: args{
				files: []string{"../../../test/v1alpha/valid-slos-ratio.yaml"},
			},
			wantErr: false,
		},
		{
			name: "v1alpha single threshold SLO file valid",
			args: args{
				files: []string{"../../../test/v1alpha/valid-slos-threshold.yaml"},
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
			name: "v1 single AlertCondition valid",
			args: args{
				files: []string{"../../../test/v1/valid-alert-condition.yaml"},
			},
			wantErr: false,
		},
		{
			name: "v1 single AlertNotificationTarget valid",
			args: args{
				files: []string{"../../../test/v1/valid-alert-notification-target.yaml"},
			},
			wantErr: false,
		},
		{
			name: "v1 single AlertPolicy valid",
			args: args{
				files: []string{"../../../test/v1/valid-alert-policy.yaml"},
			},
			wantErr: false,
		},
		{
			name: "v1 single DataSource valid",
			args: args{
				files: []string{"../../../test/v1/valid-data-source.yaml"},
			},
			wantErr: false,
		},
		{
			name: "v1 single Service valid",
			args: args{
				files: []string{"../../../test/v1/valid-sli.yaml"},
			},
			wantErr: false,
		},
		{
			name: "v1 single SLI valid",
			args: args{
				files: []string{"../../../test/v1/valid-sli.yaml"},
			},
			wantErr: false,
		},
		{
			name: "v1 single SLO valid",
			args: args{
				files: []string{"../../../test/v1/valid-slo.yaml"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := validateFiles(tt.args.files); (err != nil) != tt.wantErr {
				t.Errorf("validateFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
