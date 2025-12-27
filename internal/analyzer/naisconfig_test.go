package analyzer

import "testing"

func TestNaisConfig_IsStateful_WithDatabase(t *testing.T) {
	cfg := &NaisConfig{
		Spec: NaisSpec{
			GCP: &NaisGCP{
				SQLInstances: []NaisSQLInstance{{Type: "POSTGRES_14"}},
			},
		},
	}
	if !cfg.IsStateful() {
		t.Error("expected IsStateful() to return true when database exists")
	}
}

func TestNaisConfig_IsStateful_Nil(t *testing.T) {
	var cfg *NaisConfig
	if cfg.IsStateful() {
		t.Error("expected IsStateful() to return false for nil config")
	}
}

func TestNaisConfig_GetInboundApps(t *testing.T) {
	cfg := &NaisConfig{
		Spec: NaisSpec{
			AccessPolicy: &NaisAccessPolicy{
				Inbound: &NaisAccessRules{
					Rules: []NaisAccessRule{
						{Application: "app-a"},
						{Application: "app-b"},
					},
				},
			},
		},
	}
	apps := cfg.GetInboundApps()
	if len(apps) != 2 || apps[0] != "app-a" || apps[1] != "app-b" {
		t.Errorf("unexpected inbound apps: %v", apps)
	}
}
