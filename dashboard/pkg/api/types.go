package api

import (
	"strings"

	applicationapiv1alpha1 "github.com/redhat-appstudio/application-api/api/v1alpha1"
	integv1alpha1 "github.com/redhat-appstudio/integration-service/api/v1alpha1"
	releasev1alpha1 "github.com/redhat-appstudio/release-service/api/v1alpha1"
)

type Row struct {
	ServiceName            string
	Version                string
	Status                 string
	UsersAdded             bool
	IntegrationTestEnabled bool
	ECEnabled              bool
	ReleasePipelineEnabled bool
}

type AppConfig struct {
	Application  applicationapiv1alpha1.Application
	Components   []applicationapiv1alpha1.Component
	Tests        []integv1alpha1.IntegrationTestScenario
	ReleasePlans []releasev1alpha1.ReleasePlan
}

func (a *AppConfig) AddComponent(comp applicationapiv1alpha1.Component) {
	a.Components = append(a.Components, comp)
}

func (a *AppConfig) AddTest(test integv1alpha1.IntegrationTestScenario) {
	a.Tests = append(a.Tests, test)
}

func (a *AppConfig) AddReleasePlan(rp releasev1alpha1.ReleasePlan) {
	a.ReleasePlans = append(a.ReleasePlans, rp)
}

func (a *AppConfig) GetTestNames() string {
	testNames := make([]string, len(a.Tests))
	for _, test := range a.Tests {
		testNames = append(testNames, test.Name)
	}
	return strings.Join(testNames, ",")
}

func (a *AppConfig) GetReleasePlanNames() string {
	rpNames := make([]string, len(a.ReleasePlans))
	for _, rp := range a.ReleasePlans {
		rpNames = append(rpNames, rp.Name)
	}
	return strings.Join(rpNames, ",")
}
