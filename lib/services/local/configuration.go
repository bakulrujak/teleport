/*
Copyright 2017 Gravitational, Inc.

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

package local

import (
	"context"

	"github.com/gravitational/teleport/lib/backend"
	"github.com/gravitational/teleport/lib/services"

	"github.com/gravitational/trace"
)

// ClusterConfigurationService is responsible for managing cluster configuration.
type ClusterConfigurationService struct {
	backend.Backend
}

// NewClusterConfigurationService returns a new ClusterConfigurationService.
func NewClusterConfigurationService(backend backend.Backend) *ClusterConfigurationService {
	return &ClusterConfigurationService{
		Backend: backend,
	}
}

// GetClusterName gets the name of the cluster from the backend.
func (s *ClusterConfigurationService) GetClusterName() (services.ClusterName, error) {
	item, err := s.Get(context.TODO(), backend.Key(clusterConfigPrefix, namePrefix))
	if err != nil {
		if trace.IsNotFound(err) {
			return nil, trace.NotFound("cluster name not found")
		}
		return nil, trace.Wrap(err)
	}
	return services.GetClusterNameMarshaler().Unmarshal(item.Value)
}

// SetClusterName sets the name of the cluster in the backend. SetClusterName
// can only be called once on a cluster after which it will return trace.AlreadyExists.
func (s *ClusterConfigurationService) SetClusterName(c services.ClusterName) error {
	value, err := services.GetClusterNameMarshaler().Marshal(c)
	if err != nil {
		return trace.Wrap(err)
	}

	_, err = s.Create(context.TODO(), backend.Item{
		Key:     backend.Key(clusterConfigPrefix, namePrefix),
		Value:   value,
		Expires: c.Expiry(),
	})
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

// GetStaticTokens gets the list of static tokens used to provision nodes.
func (s *ClusterConfigurationService) GetStaticTokens() (services.StaticTokens, error) {
	item, err := s.Get(context.TODO(), backend.Key(clusterConfigPrefix, staticTokensPrefix))
	if err != nil {
		if trace.IsNotFound(err) {
			return nil, trace.NotFound("static tokens not found")
		}
		return nil, trace.Wrap(err)
	}
	return services.GetStaticTokensMarshaler().Unmarshal(item.Value)
}

// SetStaticTokens sets the list of static tokens used to provision nodes.
func (s *ClusterConfigurationService) SetStaticTokens(c services.StaticTokens) error {
	value, err := services.GetStaticTokensMarshaler().Marshal(c)
	if err != nil {
		return trace.Wrap(err)
	}

	_, err = s.Put(context.TODO(), backend.Item{
		Key:     backend.Key(clusterConfigPrefix, staticTokensPrefix),
		Value:   value,
		Expires: c.Expiry(),
	})
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

// GetAuthPreference fetches the cluster authentication preferences
// from the backend and return them.
func (s *ClusterConfigurationService) GetAuthPreference() (services.AuthPreference, error) {
	item, err := s.Get(context.TODO(), backend.Key(authPrefix, preferencePrefix, generalPrefix))
	if err != nil {
		if trace.IsNotFound(err) {
			return nil, trace.NotFound("authentication preference not found")
		}
		return nil, trace.Wrap(err)
	}
	return services.GetAuthPreferenceMarshaler().Unmarshal(item.Value)
}

// SetAuthPreference sets the cluster authentication preferences
// on the backend.
func (s *ClusterConfigurationService) SetAuthPreference(preferences services.AuthPreference) error {
	value, err := services.GetAuthPreferenceMarshaler().Marshal(preferences)
	if err != nil {
		return trace.Wrap(err)
	}

	item := backend.Item{
		Key:   backend.Key(authPrefix, preferencePrefix, generalPrefix),
		Value: value,
	}

	_, err = s.Put(context.TODO(), item)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

// GetClusterConfig gets services.ClusterConfig from the backend.
func (s *ClusterConfigurationService) GetClusterConfig() (services.ClusterConfig, error) {
	item, err := s.Get(context.TODO(), backend.Key(clusterConfigPrefix, generalPrefix))
	if err != nil {
		if trace.IsNotFound(err) {
			return nil, trace.NotFound("cluster configuration not found")
		}
		return nil, trace.Wrap(err)
	}
	return services.GetClusterConfigMarshaler().Unmarshal(item.Value)
}

// SetClusterConfig sets services.ClusterConfig on the backend.
func (s *ClusterConfigurationService) SetClusterConfig(c services.ClusterConfig) error {
	value, err := services.GetClusterConfigMarshaler().Marshal(c)
	if err != nil {
		return trace.Wrap(err)
	}

	item := backend.Item{
		Key:   backend.Key(clusterConfigPrefix, generalPrefix),
		Value: value,
	}

	_, err = s.Put(context.TODO(), item)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

const (
	clusterConfigPrefix = "cluster_configuration"
	namePrefix          = "name"
	staticTokensPrefix  = "static_tokens"
	authPrefix          = "authentication"
	preferencePrefix    = "preference"
	generalPrefix       = "general"
)
