// Copyright 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package consumer

import (
	"fmt"
	"sync"

	"github.com/vmware/hamlet/pkg/server/state"
)

// Registry maintains an active set of consumers and provides a mechanism to
// interact with them.
type Registry interface {
	// Register creates a new entry for the given consumer identified by id.
	Register(id string) (Consumer, error)

	// Deregister deregisters the consumer identified by id.
	Deregister(id string) error

	// GetAll returns all the registered consumer instances.
	GetAll() []Consumer
}

// registry is a concrete implementation of the registry interface.
type registry struct {
	Registry

	// stateProvider provides the mechanism to query the federated service
	// mesh owner implementation for the current state of a particular type
	// of resources.
	stateProvider state.StateProvider

	// consumers holds a set of registered consumers.
	consumers map[string]Consumer

	// mutex synchronizes access to the consumer listing.
	mutex *sync.Mutex
}

// NewRegistry returns a new instance of the registry.
func NewRegistry(stateProvider state.StateProvider) Registry {
	return &registry{
		stateProvider: stateProvider,
		consumers:     make(map[string]Consumer),
		mutex:         &sync.Mutex{},
	}
}

func (r *registry) Register(id string) (Consumer, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, found := r.consumers[id]; found {
		return nil, fmt.Errorf("Consumer with id %s already exists", id)
	}

	r.consumers[id] = newConsumer(id, r.stateProvider)
	return r.consumers[id], nil
}

func (r *registry) Deregister(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.consumers, id)
	return nil
}

func (r *registry) GetAll() []Consumer {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	consumers := make([]Consumer, 0, len(r.consumers))
	for _, consumer := range r.consumers {
		consumers = append(consumers, consumer)
	}
	return consumers
}
