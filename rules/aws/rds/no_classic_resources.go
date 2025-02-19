package rds

import (
	"github.com/aquasecurity/defsec/providers"
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/severity"
	"github.com/aquasecurity/defsec/state"
)

var CheckNoClassicResources = rules.Register(
	rules.Rule{
		AVDID:      "AVD-AWS-0081",
		Provider:   providers.AWSProvider,
		Service:    "rds",
		ShortCode:  "no-classic-resources",
		Summary:    "AWS Classic resource usage.",
		Impact:     "Classic resources are running in a shared environment with other customers",
		Resolution: "Switch to VPC resources",
		Explanation: `AWS Classic resources run in a shared environment with infrastructure owned by other AWS customers. You should run
resources in a VPC instead.`,
		Links: []string{
			"https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-classic-platform.html",
		},
		Terraform: &rules.EngineMetadata{
			GoodExamples:        terraformNoClassicResourcesGoodExamples,
			BadExamples:         terraformNoClassicResourcesBadExamples,
			Links:               terraformNoClassicResourcesLinks,
			RemediationMarkdown: terraformNoClassicResourcesRemediationMarkdown,
		},
		CloudFormation: &rules.EngineMetadata{
			GoodExamples:        cloudFormationNoClassicResourcesGoodExamples,
			BadExamples:         cloudFormationNoClassicResourcesBadExamples,
			Links:               cloudFormationNoClassicResourcesLinks,
			RemediationMarkdown: cloudFormationNoClassicResourcesRemediationMarkdown,
		},
		Severity: severity.Critical,
	},
	func(s *state.State) (results rules.Results) {
		for _, group := range s.AWS.RDS.Classic.DBSecurityGroups {
			results.Add(
				"Classic resources should not be used.",
				&group,
			)
		}
		return
	},
)
