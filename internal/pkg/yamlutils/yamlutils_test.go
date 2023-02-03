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
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/OpenSLO/oslo/pkg/manifest"
	v1 "github.com/OpenSLO/oslo/pkg/manifest/v1"
)

//go:embed test-input
var testInput string

func TestParse(t *testing.T) {
	t.Parallel()
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
				fileContent: []byte(testInput),
				filename:    "test.yaml",
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, _, err := Parse(tt.args.fileContent, tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, 2, len(got))
			assert.Equal(t, tt.want, got)
		})
	}
}
