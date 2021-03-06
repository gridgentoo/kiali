package models

// IstioTypeValidations represents a set of IstioNameValidations grouper per Istio ObjectType.
// It is possible that different object types have same name, but names per ObjectType are unique.
type IstioTypeValidations map[string]*IstioNameValidations

// IstioNameValidations represents a set of checks grouped per Istio object.
// The key of the map represents the name of the Istio object.
type IstioNameValidations map[string]*IstioValidation

type IstioValidation struct {
	Name       string        `json:"name"`
	ObjectType string        `json:"object_type"`
	Valid      bool          `json:"valid"`
	Checks     []*IstioCheck `json:"checks"`
}

type IstioCheck struct {
	Message  string `json:"message"`
	Severity string `json:"severity"`
	Path     string `json:"path"`
}

func BuildCheck(message, severity, path string) IstioCheck {
	return IstioCheck{message, severity, path}
}

func (in IstioTypeValidations) MergeValidations(typeValidations *IstioTypeValidations) IstioTypeValidations {
	for objectType, nameValidations := range *typeValidations {
		if in[objectType] != nil {
			in[objectType].MergeNameValidations(nameValidations)
		} else {
			in[objectType] = nameValidations
		}
	}

	return in
}

func (in IstioNameValidations) MergeNameValidations(nameValidations *IstioNameValidations) IstioNameValidations {
	for name, validation := range *nameValidations {
		if in[name] != nil {
			in[name].Checks = append(in[name].Checks, validation.Checks...)
			in[name].Valid = validation.Valid && in[name].Valid
		} else {
			in[name] = validation
		}
	}

	return in
}
