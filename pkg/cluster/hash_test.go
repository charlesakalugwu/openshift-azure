package cluster

import (
	"context"
	"encoding/hex"
	"io/ioutil"
	"testing"

	"github.com/ghodss/yaml"

	"github.com/openshift/openshift-azure/pkg/api"
	pluginapi "github.com/openshift/openshift-azure/pkg/api/plugin"
	"github.com/openshift/openshift-azure/pkg/arm"
	"github.com/openshift/openshift-azure/pkg/metrics"
	"github.com/openshift/openshift-azure/pkg/startup"
	"github.com/openshift/openshift-azure/pkg/sync"
	"github.com/openshift/openshift-azure/test/util/populate"
)

// These tests would be stronger were we to populate all the hash inputs rather
// than leaving most of them at the zero value.  However this approach is
// hampered by the fact that we do not version the internal representation. Were
// we not to use zero values, future refactors of the internal representation
// could move fields around and easily provoke false positives here.

func TestHashScaleSetStability(t *testing.T) {
	// IMPORTANT: hashes are free to change for the version under development,
	// but should be fixed thereafter.  If this test is failing against a
	// released version, do not change the hash here, but fix the underlying
	// problem, otherwise it will result in unwanted rotations in production.

	tests := map[string][]struct {
		role         api.AgentPoolProfileRole
		expectedHash string
	}{
		"v4.2": {
			{
				role: api.AgentPoolProfileRoleMaster,
				// this value should not change
				expectedHash: "9a4b049bda1bda0bf394b6e5a21c89bb8b219346261de0e6b01d2efb100f3e01",
			},
			{
				role: api.AgentPoolProfileRoleInfra,
				// this value should not change
				expectedHash: "a3dd5cb39b243e4e67a2d59cdda59d9179e39a87d6b6201ae7e6ce78c6317a72",
			},
			{
				role: api.AgentPoolProfileRoleCompute,
				// this value should not change
				expectedHash: "0b9ada416dcffb9b0ecc5e96bc21b2a9ceaa78e74af82819040cdc14854a3a12",
			},
		},
		"v4.3": {
			{
				role: api.AgentPoolProfileRoleMaster,
				// this value should not change
				expectedHash: "87d0e3642432a78bd3b65d753116e2988324a22a22468c54e7b59fc31185524c",
			},
			{
				role: api.AgentPoolProfileRoleInfra,
				// this value should not change
				expectedHash: "afa8a3e0cbe7264b811f231ca7577f213022ae7a8a7a3c92da715c5c813213c3",
			},
			{
				role: api.AgentPoolProfileRoleCompute,
				// this value should not change
				expectedHash: "f6f60ca20f23036807e7c54686c304a653830028864830dd31b58de7358e3fee",
			},
		},
		"v4.4": {
			{
				role: api.AgentPoolProfileRoleMaster,
				// this value should not change
				expectedHash: "87d0e3642432a78bd3b65d753116e2988324a22a22468c54e7b59fc31185524c",
			},
			{
				role: api.AgentPoolProfileRoleInfra,
				// this value should not change
				expectedHash: "afa8a3e0cbe7264b811f231ca7577f213022ae7a8a7a3c92da715c5c813213c3",
			},
			{
				role: api.AgentPoolProfileRoleCompute,
				// this value should not change
				expectedHash: "f6f60ca20f23036807e7c54686c304a653830028864830dd31b58de7358e3fee",
			},
		},
		"v5.1": {
			{
				role: api.AgentPoolProfileRoleMaster,
				// this value should not change
				expectedHash: "4d4d1c053cba07da9641a817d7ff5712a40d6f43ba10306dc802a4d3d8b71932",
			},
			{
				role: api.AgentPoolProfileRoleInfra,
				// this value should not change
				expectedHash: "d954d5572f9a67de96bbaf97fc89c62c5625812890a6cc39debd9c6fb7e0dd47",
			},
			{
				role: api.AgentPoolProfileRoleCompute,
				// this value should not change
				expectedHash: "c43d513804031e9b14f076a53c5f02b8ae121926ca7e2f33cbaad1b22cde89b8",
			},
		},
	}

	// check we're testing all versions in our pluginconfig
	b, err := ioutil.ReadFile("../../pluginconfig/pluginconfig-311.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var template *pluginapi.Config
	err = yaml.Unmarshal(b, &template)
	if err != nil {
		t.Fatal(err)
	}

	for version := range template.Versions {
		if _, found := tests[version]; !found && version != template.PluginVersion {
			t.Errorf("update tests to include version %s", version)
		}
	}

	cs := &api.OpenShiftManagedCluster{
		Properties: api.Properties{
			RouterProfiles: []api.RouterProfile{
				{},
			},
			AgentPoolProfiles: []api.AgentPoolProfile{
				{
					Role: api.AgentPoolProfileRoleMaster,
				},
				{
					Role: api.AgentPoolProfileRoleInfra,
				},
				{
					Role:   api.AgentPoolProfileRoleCompute,
					VMSize: api.StandardD2sV3,
				},
			},
			AuthProfile: api.AuthProfile{
				IdentityProviders: []api.IdentityProvider{
					{
						Provider: &api.AADIdentityProvider{},
					},
				},
			},
		},
	}
	populate.DummyCertsAndKeys(cs)

	for version, tests := range tests {
		for _, tt := range tests {
			cs.Config.PluginVersion = version

			arm, err := arm.New(context.Background(), nil, cs, api.TestConfig{})
			if err != nil {
				t.Fatal(err)
			}

			hasher := Hash{
				StartupFactory: startup.New,
				Arm:            arm,
			}

			b, err := hasher.HashScaleSet(cs, &api.AgentPoolProfile{Role: tt.role})
			if err != nil {
				t.Fatal(err)
			}

			h := hex.EncodeToString(b)
			if h != tt.expectedHash {
				t.Errorf("%s: %s: hash changed to %s", version, tt.role, h)
			}
		}
	}
}

func TestHashSyncPodStability(t *testing.T) {
	// IMPORTANT: hashes are free to change for the version under development,
	// but should be fixed thereafter.  If this test is failing against a
	// released version, do not change the hash here, but fix the underlying
	// problem, otherwise it will result in unwanted rotations in production.

	tests := map[string]struct {
		expectedHash string
	}{
		"v4.2": {
			// this value should not change
			expectedHash: "f731d62dc04f805314b912c3720189679cc404a35ee41ee3569fb3f642c8bd29",
		},
		"v4.3": {
			// this value should not change
			expectedHash: "f731d62dc04f805314b912c3720189679cc404a35ee41ee3569fb3f642c8bd29",
		},
		"v4.4": {
			// this value should not change
			expectedHash: "f731d62dc04f805314b912c3720189679cc404a35ee41ee3569fb3f642c8bd29",
		},
		"v5.1": {
			// this value should not change
			expectedHash: "fac2ba9ab18ad8d26f99b9f2d694529463b024974efdfd2a964a12df7a1e545a",
		},
	}

	// check we're testing all versions in our pluginconfig
	b, err := ioutil.ReadFile("../../pluginconfig/pluginconfig-311.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var template *pluginapi.Config
	err = yaml.Unmarshal(b, &template)
	if err != nil {
		t.Fatal(err)
	}

	for version := range template.Versions {
		if _, found := tests[version]; !found && version != template.PluginVersion {
			t.Errorf("update tests to include version %s", version)
		}
	}

	cs := &api.OpenShiftManagedCluster{
		Properties: api.Properties{
			RouterProfiles: []api.RouterProfile{
				{},
			},
			AuthProfile: api.AuthProfile{
				IdentityProviders: []api.IdentityProvider{
					{
						Provider: &api.AADIdentityProvider{},
					},
				},
			},
		},
		Config: api.Config{
			ImageVersion: "311.0.0",
			Images: api.ImageConfig{
				AlertManager:             ":",
				ConfigReloader:           ":",
				Grafana:                  ":",
				KubeRbacProxy:            ":",
				KubeStateMetrics:         ":",
				NodeExporter:             ":",
				OAuthProxy:               ":",
				Prometheus:               ":",
				PrometheusConfigReloader: ":",
				PrometheusOperator:       ":",
			},
		},
	}
	populate.DummyCertsAndKeys(cs)

	for version, tt := range tests {
		cs.Config.PluginVersion = version

		s, err := sync.New(nil, cs, false, metrics.DefaultCollector())
		if err != nil {
			t.Fatal(err)
		}

		b, err := s.Hash()
		if err != nil {
			t.Fatal(err)
		}

		h := hex.EncodeToString(b)
		if h != tt.expectedHash {
			t.Errorf("%s: hash changed to %s", version, h)
		}
	}
}
