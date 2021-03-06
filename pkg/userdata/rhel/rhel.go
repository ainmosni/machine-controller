/*
Copyright 2019 The Machine Controller Authors.

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

package rhel

import (
	"encoding/json"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
)

// Config contains specific configuration for RHEL.
type Config struct {
	DistUpgradeOnBoot               bool   `json:"distUpgradeOnBoot"`
	RHELSubscriptionManagerUser     string `json:"rhelSubscriptionManagerUser,omitempty"`
	RHELSubscriptionManagerPassword string `json:"rhelSubscriptionManagerPassword,omitempty"`
}

// LoadConfig retrieves the RHEL configuration from raw data.
func LoadConfig(r runtime.RawExtension) (*Config, error) {
	cfg := &Config{}
	if len(r.Raw) == 0 {
		return cfg, nil
	}
	if err := json.Unmarshal(r.Raw, cfg); err != nil {
		return nil, err
	}

	if cfg.RHELSubscriptionManagerUser == "" {
		subUser, ok := os.LookupEnv("RHEL_SUBSCRIPTION_MANAGER_USER")
		if !ok {
			return nil, fmt.Errorf("RHEL_SUBSCRIPTION_MANAGER_USER env variable is not found")
		}
		cfg.RHELSubscriptionManagerUser = subUser
	}

	if cfg.RHELSubscriptionManagerPassword == "" {
		subPassword, ok := os.LookupEnv("RHEL_SUBSCRIPTION_MANAGER_PASSWORD")
		if !ok {
			return nil, fmt.Errorf("RHEL_SUBSCRIPTION_MANAGER_PASSWORD env variable is not found")
		}
		cfg.RHELSubscriptionManagerPassword = subPassword
	}

	return cfg, nil
}

// Spec return the configuration as raw data.
func (cfg *Config) Spec() (*runtime.RawExtension, error) {
	ext := &runtime.RawExtension{}
	b, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	ext.Raw = b
	return ext, nil
}
