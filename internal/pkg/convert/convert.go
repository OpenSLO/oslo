/*
Package convert provides a command to convert from openslo to other formats.

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
package convert

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"

	nobl9manifest "github.com/OpenSLO/oslo/internal/pkg/manifest/nobl9"
	nobl9v1alpha "github.com/OpenSLO/oslo/internal/pkg/manifest/nobl9/v1alpha"
	"github.com/OpenSLO/oslo/internal/pkg/yamlutils"
	"github.com/OpenSLO/oslo/pkg/manifest"
	v1 "github.com/OpenSLO/oslo/pkg/manifest/v1"
)

// RemoveDuplicates to remove duplicate string from a slice.
func RemoveDuplicates(s []string) []string {
	result := make([]string, 0, len(s))
	m := make(map[string]bool)
	for _, v := range s {
		if _, ok := m[v]; !ok {
			m[v] = true
		}
	}
	for k := range m {
		result = append(result, k)
	}
	return result
}

// Files converts the provided file.
func Files(out io.Writer, filenames []string) error {
	var rval []interface{}
	parsed, err := getParsedObjects(filenames)
	if err != nil {
		return fmt.Errorf("issue parsing content: %w", err)
	}

	// Get the service objects.
	if err := getServiceObjects(parsed, &rval); err != nil {
		return fmt.Errorf("issue getting service objects: %w", err)
	}

	// Get the alertPolicy objects.
	if err := getAlertPolicyObjects(parsed, &rval); err != nil {
		return fmt.Errorf("issue getting alertPolicy objects: %w", err)
	}

	// Print out all of our objects.
	for _, s := range rval {
		err := printYaml(out, s)
		if err != nil {
			return fmt.Errorf("issue printing content: %w", err)
		}
	}
	// foo, err := getObjectByKind("AlertPolicy", parsed)
	// if err != nil {
	//   return fmt.Errorf("issue getting alert policy: %w", err)
	// }
	// alertPolicy, ok := foo.(*v1.AlertPolicy)
	// if !ok {
	//   return fmt.Errorf("issue casting to alert policy")
	// }
	//
	// // Find which objects we have
	// for _, p := range parsed {
	//   switch pp := p.(type) {
	//   case v1.AlertPolicy:
	//     object := nobl9v1alpha.AlertPolicy{
	//       ObjectHeader: getObjectHeader("Service", pp.Metadata.Name, pp.Metadata.DisplayName, "default"),
	//       Spec: nobl9v1alpha.AlertPolicySpec{
	//         Description: pp.Spec.Description,
	//       },
	//     }
	//     printYaml(out, object)
	//
	//   case v1.Service:
	//     printYaml(out, object)
	//   }
	// }

	// For each Nobl9 kind, try and create it from the OpenSLOKind
	return nil
}

// Constructs Nobl9 AlertPolicy objects from our list of OpenSLOKinds.
func getAlertPolicyObjects(parsed []manifest.OpenSLOKind, rval *[]interface{}) error {
	// Get the alert policy object.
	ap, err := getObjectByKind("AlertPolicy", parsed)
	if err != nil {
		return fmt.Errorf("issue getting alert policy from parsed list: %w", err)
	}

	// Return if ap is empty.
	if len(ap) == 0 {
		return nil
	}

	// AlertCondition is required so get any from our parsed list.
	ac, err := getObjectByKind("AlertCondition", parsed)
	if err != nil {
		return fmt.Errorf("issue getting alert condition from parsed list: %w", err)
	}

	// For each AlertPolicy
	for _, o := range ap {
		// Cast to OpenSLO service objects.
		apObj, ok := o.(v1.AlertPolicy)
		if !ok {
			return fmt.Errorf("issue casting to AlertPolicy")
		}

		// Gather the alert conditions.
		var conditions []nobl9v1alpha.AlertCondition
		for _, a := range apObj.Spec.Conditions {
			err := getAlertCondition(a, &conditions, ac)
			if err != nil {
				return fmt.Errorf("issue getting alert condition: %w", err)
			}
		}

		// Construct the nobl9 AlertPolicy object from the OpenSLO AlertPolicy object.
		*rval = append(*rval, nobl9v1alpha.AlertPolicy{
			ObjectHeader: getObjectHeader("AlertPolicy", apObj.Metadata.Name, apObj.Metadata.DisplayName, "default"),
			Spec: nobl9v1alpha.AlertPolicySpec{
				Description:      apObj.Spec.Description,
				Conditions:       conditions,
				Severity:         "high",
				CoolDownDuration: "5m", // default
			},
		})
	}
	return nil
}

// returns an nobl9v1alpha.AlertCondition from an OpenSLO.AlertPolicyCondition
func getAlertCondition(apc v1.AlertPolicyCondition, conditions *[]nobl9v1alpha.AlertCondition, ac []manifest.OpenSLOKind) error {
	// If we have an inline condition, we can use it.
	if apc.AlertConditionInline != nil {
		*conditions = append(*conditions, nobl9v1alpha.AlertCondition{
			Measurement:      "averageBurnRate", // TODO: add other measurements in OpenSLO
			Value:            apc.AlertConditionInline.Spec.Condition.Threshold,
			LastsForDuration: apc.AlertConditionInline.Spec.Condition.AlertAfter,
			Operation:        "gt", // TODO add this to OpenSLO
		})
	} else {
		// Error if we don't have any, since we need at least one.
		if len(ac) == 0 {
			return fmt.Errorf("no alert conditions found. Required for alert policy")
		}

		// If we don't have an inline condition, we need to get the AlertCondition.
		for _, c := range ac {
			// Get the AlertCondition that matches the name.
			acObj, ok := c.(v1.AlertCondition)
			if !ok {
				return fmt.Errorf("issue casting to AlertCondition")
			}
			if apc.AlertPolicyConditionSpec.ConditionRef == acObj.Metadata.Name {
				*conditions = append(*conditions, nobl9v1alpha.AlertCondition{
					Measurement:      "averageBurnRate",
					Value:            acObj.Spec.Condition.Threshold,
					LastsForDuration: acObj.Spec.Condition.AlertAfter,
					Operation:        "gt",
				})
			} else {
				return fmt.Errorf("alert condition %s not found", apc.AlertPolicyConditionSpec.ConditionRef)
			}
		}
	}
	return nil
}

func getParsedObjects(filenames []string) ([]manifest.OpenSLOKind, error) {
	var parsed []manifest.OpenSLOKind
	for _, filename := range filenames {
		// Get the file contents.
		content, err := yamlutils.ReadConf(filename)
		if err != nil {
			return nil, fmt.Errorf("issue reading content: %w", err)
		}

		// Parse the byte arrays to OpenSLOKind objects.
		p, err := yamlutils.Parse(content, filename)
		if err != nil {
			return nil, fmt.Errorf("issue parsing content: %w", err)
		}

		parsed = append(parsed, p...)
	}
	return parsed, nil
}

// function that that returns an object by Kind from a list of OpenSLOKinds.
func getObjectByKind(kind string, objects []manifest.OpenSLOKind) ([]manifest.OpenSLOKind, error) {
	var found []manifest.OpenSLOKind
	for _, o := range objects {
		if o.Kind() == kind {
			found = append(found, o)
		}
	}
	// TODO: I dont think that we care about the length here.
	// if len(found) == 0 {
	//   return nil, fmt.Errorf("no %s found", kind)
	// }
	return found, nil
}

// function that takes a manifest.OpenSLOKind and returns a nobl9v1alpha.ObjectHeader.
func getObjectHeader(kind, name, displayName, project string) nobl9v1alpha.ObjectHeader {
	return nobl9v1alpha.ObjectHeader{
		ObjectHeader: nobl9manifest.ObjectHeader{
			APIVersion: nobl9v1alpha.APIVersion,
		},
		Kind: kind,
		MetadataHolder: nobl9v1alpha.MetadataHolder{
			Metadata: nobl9v1alpha.Metadata{
				Name:        name,
				DisplayName: displayName,
				Project:     project,
			},
		},
	}
}

// Constructs Nobl9 Service objects from our list of OpenSLOKinds.
func getServiceObjects(parsed []manifest.OpenSLOKind, rval *[]interface{}) error {
	// Get the service object.
	obj, err := getObjectByKind("Service", parsed)
	if err != nil {
		return fmt.Errorf("issue getting service from parsed list: %w", err)
	}

	for _, o := range obj {
		// Cast to OpenSLO service objects.
		srvObj, ok := o.(v1.Service)
		if !ok {
			return fmt.Errorf("issue casting to service")
		}
		// Construct the nobl9 service object from the OpenSLO service object.
		*rval = append(*rval, nobl9v1alpha.Service{
			ObjectHeader: getObjectHeader("Service", srvObj.Metadata.Name, srvObj.Metadata.DisplayName, "default"),
			Spec: nobl9v1alpha.ServiceSpec{
				Description: srvObj.Spec.Description,
			},
		})
	}
	return nil
}

func printYaml(out io.Writer, object interface{}) error {
	// Convert parsed to yaml and print to out.
	yml, err := yaml.Marshal(object)
	if err != nil {
		return fmt.Errorf("issue marshaling content: %w", err)
	}
	_, err = out.Write(yml)
	if err != nil {
		return fmt.Errorf("issue writing content: %w", err)
	}

	return nil
}
