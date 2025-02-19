package openstack

import (
	"testing"

	"github.com/aquasecurity/defsec/adapters/terraform/testutil"
	"github.com/aquasecurity/defsec/parsers/types"
	"github.com/aquasecurity/defsec/providers/openstack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFields(t *testing.T) {
	tests := []struct {
		name      string
		terraform string
		expected  openstack.OpenStack
	}{
		{
			name: "Plaintext password",
			terraform: `
			resource "openstack_compute_instance_v2" "my-instance" {
			  admin_pass      = "N0tSoS3cretP4ssw0rd"

			}`,
			expected: openstack.OpenStack{
				Metadata: types.NewTestMetadata(),
				Compute: openstack.Compute{
					Metadata: types.NewTestMetadata(),
					Instances: []openstack.Instance{
						{
							Metadata:      types.NewTestMetadata(),
							AdminPassword: types.String("N0tSoS3cretP4ssw0rd", types.NewTestMetadata()),
						},
					},
				},
			},
		},
		{
			name: "No plaintext password",
			terraform: `
			resource "openstack_compute_instance_v2" "my-instance" {
			}`,
			expected: openstack.OpenStack{
				Metadata: types.NewTestMetadata(),
				Compute: openstack.Compute{
					Metadata: types.NewTestMetadata(),
					Instances: []openstack.Instance{
						{
							Metadata:      types.NewTestMetadata(),
							AdminPassword: types.String("", types.NewTestMetadata()),
						},
					},
				},
			},
		},
		{
			name: "Firewall rule",
			terraform: `
			resource "openstack_fw_rule_v1" "rule_1" {
				action                 = "allow"
				protocol               = "tcp"
				destination_port       = "22"
				destination_ip_address = "10.10.10.1"
				source_ip_address      = "10.10.10.2"
				enabled                = "true"
			}`,
			expected: openstack.OpenStack{
				Metadata: types.NewTestMetadata(),
				Compute: openstack.Compute{
					Metadata: types.NewTestMetadata(),
					Firewall: openstack.Firewall{
						Metadata: types.NewTestMetadata(),
						AllowRules: []openstack.FirewallRule{
							{
								Metadata:        types.NewTestMetadata(),
								Enabled:         types.Bool(true, types.NewTestMetadata()),
								Destination:     types.String("10.10.10.1", types.NewTestMetadata()),
								Source:          types.String("10.10.10.2", types.NewTestMetadata()),
								DestinationPort: types.String("22", types.NewTestMetadata()),
								SourcePort:      types.String("", types.NewTestMetadata()),
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			modules := testutil.CreateModulesFromSource(test.terraform, ".tf", t)
			adapted := Adapt(modules)
			testutil.AssertDefsecEqual(t, test.expected, adapted)
		})
	}
}

func TestLines(t *testing.T) {
	src := `
	resource "openstack_compute_instance_v2" "my-instance" {
		admin_pass      = "N0tSoS3cretP4ssw0rd"
	}

	resource "openstack_fw_rule_v1" "rule_1" {
		action                 = "allow"
		protocol               = "tcp"
		destination_port       = "22"
		destination_ip_address = "10.10.10.1"
		source_ip_address      = "10.10.10.2"
		enabled                = "true"
	}`

	modules := testutil.CreateModulesFromSource(src, ".tf", t)
	adapted := Adapt(modules)

	require.Len(t, adapted.Compute.Instances, 1)
	instance := adapted.Compute.Instances[0]

	require.Len(t, adapted.Compute.Firewall.AllowRules, 1)
	rule := adapted.Compute.Firewall.AllowRules[0]

	assert.Equal(t, 3, instance.AdminPassword.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 3, instance.AdminPassword.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 9, rule.DestinationPort.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 9, rule.DestinationPort.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 10, rule.Destination.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 10, rule.Destination.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 11, rule.Source.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 11, rule.Source.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 12, rule.Enabled.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 12, rule.Enabled.GetMetadata().Range().GetEndLine())
}
