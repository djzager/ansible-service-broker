package apb

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	ft "github.com/openshift/ansible-service-broker/pkg/fusortest"
	yaml "gopkg.in/yaml.v2"
)

const PlanName = "dev"
const PlanDescription = "Mediawiki123 apb implementation"

var PlanMetadata = map[string]interface{}{
	"displayName":     "Development",
	"longDescription": "Basic development plan",
	"cost":            "$0.00",
}

const PlanFree = true
const PlanBindable = true

var expectedPlanParameters = []ParameterDescriptor{
	ParameterDescriptor{
		Name:     "mediawiki_db_schema",
		Title:    "Mediawiki DB Schema",
		Type:     "string",
		Default:  "mediawiki",
		Required: true},
	ParameterDescriptor{
		Name:     "mediawiki_site_name",
		Title:    "Mediawiki Site Name",
		Type:     "string",
		Default:  "MediaWiki",
		Required: true},
	ParameterDescriptor{
		Name:     "mediawiki_site_lang",
		Title:    "Mediawiki Site Language",
		Type:     "string",
		Default:  "en",
		Required: true},
	ParameterDescriptor{
		Name:     "mediawiki_admin_user",
		Title:    "Mediawiki Admin User",
		Type:     "string",
		Default:  "admin",
		Required: true},
	ParameterDescriptor{
		Name:     "mediawiki_admin_pass",
		Title:    "Mediawiki Admin User Password",
		Type:     "string",
		Required: true},
}

var p = Plan{
	Name:        PlanName,
	Description: PlanDescription,
	Metadata:    PlanMetadata,
	Free:        PlanFree,
	Bindable:    PlanBindable,
	Parameters:  expectedPlanParameters,
}

const SpecName = "mediawiki123-apb"
const SpecImage = "ansibleplaybookbundle/mediawiki123-apb"
const SpecBindable = false
const SpecAsync = "optional"
const SpecDescription = "Mediawiki123 apb implementation"
const SpecPlans = `
[
	{
		"name": "dev",
		"description": "Mediawiki123 apb implementation",
		"free": true,
		"bindable": true,
		"metadata": {
			"displayName": "Development",
			"longDescription": "Basic development plan",
			"cost": "$0.00"
		},
		"parameters": [
		{
			"name": "mediawiki_db_schema",
			"title": "Mediawiki DB Schema",
			"type": "string",
			"default": "mediawiki",
			"required": true
		},
		{
			"name": "mediawiki_site_name",
			"title": "Mediawiki Site Name",
			"type": "string",
			"default": "MediaWiki",
			"required": true
		},
		{
			"name": "mediawiki_site_lang",
			"title": "Mediawiki Site Language",
			"type": "string",
			"default": "en",
			"required": true
		},
		{
			"name": "mediawiki_admin_user",
			"title": "Mediawiki Admin User",
			"type": "string",
			"default": "admin",
			"required": true
		},
		{
			"name": "mediawiki_admin_pass",
			"title": "Mediawiki Admin User Password",
			"type": "string",
			"required": true
		}
		]
	}
]
`

var SpecJSON = fmt.Sprintf(`
{
	"id": "",
	"tags": null,
	"description": "%s",
	"name": "%s",
	"image": "%s",
	"bindable": %t,
	"async": "%s",
	"plans": %s
}
`, SpecDescription, SpecName, SpecImage, SpecBindable, SpecAsync, SpecPlans)

func TestSpecLoadJSON(t *testing.T) {

	s := Spec{}
	err := LoadJSON(SpecJSON, &s)
	if err != nil {
		panic(err)
	}

	ft.AssertEqual(t, s.Description, SpecDescription)
	ft.AssertEqual(t, s.FQName, SpecName)
	ft.AssertEqual(t, s.Image, SpecImage)
	ft.AssertEqual(t, s.Bindable, SpecBindable)
	ft.AssertEqual(t, s.Async, SpecAsync)
	ft.AssertTrue(t, reflect.DeepEqual(s.Plans[0].Parameters, expectedPlanParameters))
}

func TestSpecDumpJSON(t *testing.T) {
	s := Spec{
		Description: SpecDescription,
		FQName:      SpecName,
		Image:       SpecImage,
		Bindable:    SpecBindable,
		Async:       SpecAsync,
		Plans:       []Plan{p},
	}

	var knownMap interface{}
	var subjectMap interface{}

	raw, err := DumpJSON(&s)
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(SpecJSON), &knownMap)
	json.Unmarshal([]byte(raw), &subjectMap)
	ft.AssertTrue(t, reflect.DeepEqual(knownMap, subjectMap))
}

func TestEncodedParameters(t *testing.T) {
	decodedyaml, err := base64.StdEncoding.DecodeString(ft.EncodedApb())
	if err != nil {
		t.Fatal(err)
	}

	spec := &Spec{}
	if err = yaml.Unmarshal(decodedyaml, spec); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v", spec)
	ft.AssertEqual(t, spec.FQName, "mediawiki123-apb")
	ft.AssertEqual(t, len(spec.Plans[0].Parameters), 5)

	// picking something other than the first one
	sitelang := spec.Plans[0].Parameters[2] // mediawiki_site_lang

	ft.AssertEqual(t, sitelang.Name, "mediawiki_site_lang")
	ft.AssertEqual(t, sitelang.Title, "Mediawiki Site Language")
	ft.AssertEqual(t, sitelang.Type, "string")
	ft.AssertEqual(t, sitelang.Description, "")
	ft.AssertEqual(t, sitelang.Default, "en")
	ft.AssertEqual(t, sitelang.Maxlength, 0)
	ft.AssertEqual(t, sitelang.Pattern, "")
	ft.AssertEqual(t, len(sitelang.Enum), 0)
}
