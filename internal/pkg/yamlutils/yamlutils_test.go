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
package yamlutils

import (
	"testing"

	"github.com/OpenSLO/oslo/pkg/manifest"
	v1 "github.com/OpenSLO/oslo/pkg/manifest/v1"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type args struct {
		fileContent []byte
		filename    string
	}
	tests := []struct {
		name    string
		args    args
		want    []manifest.OpenSLOKind
		wantErr bool
	}{
		{
			name: "TestParse",
			args: args{
				fileContent: []byte(`---
apiVersion: openslo/v1
kind: Service
metadata:
  name: my-rad-service
  displayName: My Rad Service
spec:
  description: This is a great description of an even better service.
---
apiVersion: openslo/v1
kind: Service
metadata:
  name: my-rad-service-deux
  displayName: My Rad Service le Deux
spec:
  description: This is a great description of an even better service.
`),
				filename: "test.yaml",
			},
			want: []manifest.OpenSLOKind{
				v1.Service{
					ObjectHeader: v1.ObjectHeader{
						ObjectHeader: manifest.ObjectHeader{
							APIVersion: "openslo/v1",
						},
						Kind: "Service",
						MetadataHolder: v1.MetadataHolder{
							Metadata: v1.Metadata{
								Name:        "my-rad-service",
								DisplayName: "My Rad Service",
							},
						},
					},
					Spec: v1.ServiceSpec{
						Description: "This is a great description of an even better service.",
					},
				},
				v1.Service{
					ObjectHeader: v1.ObjectHeader{
						ObjectHeader: manifest.ObjectHeader{
							APIVersion: "openslo/v1",
						},
						Kind: "Service",
						MetadataHolder: v1.MetadataHolder{
							Metadata: v1.Metadata{
								Name:        "my-rad-service-deux",
								DisplayName: "My Rad Service le Deux",
							},
						},
					},
					Spec: v1.ServiceSpec{
						Description: "This is a great description of an even better service.",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.fileContent, tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, 2, len(got))
			assert.Equal(t, tt.want, got)
		})
	}
}
