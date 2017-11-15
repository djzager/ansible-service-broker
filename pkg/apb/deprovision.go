//
// Copyright (c) 2017 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package apb

import (
	logging "github.com/op/go-logging"
	"github.com/openshift/ansible-service-broker/pkg/metrics"
	"github.com/openshift/ansible-service-broker/pkg/runtime"
	"github.com/pkg/errors"
)

// Deprovision - runs the abp with the deprovision action.
func Deprovision(
	instance *ServiceInstance, clusterConfig ClusterConfig, log *logging.Logger,
) (string, error) {
	log.Notice("============================================================")
	log.Notice("                      DEPROVISIONING                        ")
	log.Notice("============================================================")
	log.Noticef("ServiceInstance.Id: %s", instance.Spec.ID)
	log.Noticef("ServiceInstance.Name: %v", instance.Spec.FQName)
	log.Noticef("ServiceInstance.Image: %s", instance.Spec.Image)
	log.Noticef("ServiceInstance.Description: %s", instance.Spec.Description)
	log.Notice("============================================================")

	// Explicitly error out if image field is missing from instance.Spec
	// was introduced as a change to the apb instance.Spec to support integration
	// with the broker and still allow for providing an img path
	// Legacy ansibleapps will hit this.
	// TODO: Move this validation to a Spec creation function (yet to be created)
	if instance.Spec.Image == "" {
		log.Error("No image field found on the apb instance.Spec (apb.yaml)")
		log.Error("apb instance.Spec requires [name] and [image] fields to be separate")
		log.Error("Are you trying to run a legacy ansibleapp without an image field?")
		return "", errors.New("No image field found on instance.Spec")
	}

	// Might need to change up this interface to feed in instance ids
	metrics.ActionStarted("deprovision")
	executionContext, err := ExecuteApb(
		"deprovision",
		clusterConfig,
		instance.Spec,
		instance.Context,
		instance.Parameters,
		log,
	)
	defer runtime.Provider.DestroySandbox(executionContext.PodName, executionContext.Namespace, executionContext.Targets, clusterConfig.Namespace, clusterConfig.KeepNamespace, clusterConfig.KeepNamespaceOnError)
	if err != nil {
		log.Errorf("Problem executing apb [%s] deprovision", executionContext.PodName)
		return executionContext.PodName, err
	}

	return executionContext.PodName, err
}
