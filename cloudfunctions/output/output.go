package output

// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/googlecloudplatform/security-response-automation/services"
)

var topics = map[string]struct{ Topic string }{
	"turbinia": {Topic: "notify-turbinia"},
}

// Services contains the services needed for this function.
type Services struct {
	Logger *services.Logger
	PubSub *services.PubSub
}

// Output contains the output of this function.
type Output struct {
	// DiskNames optionally contains the names of the disks copied to a target project.
	DiskNames []string
	// Project references which project the remediation was taken, it can be used to send emails or create issues on pager duty
	Project string
	// Zone references which project the remediation was taken
	Zone string
	// Reason
	Reason string
}

// Values are requirements for this function.
type Values struct {
	Name    string
	Message []byte
}

// Execute will route & publish the incoming message to the appropriate output function.
func Execute(ctx context.Context, v *Values, services *Services) error {
	log.Printf("executing output %q", v.Name)
	if topic, ok := topics[v.Name]; ok {
		if _, err := services.PubSub.Publish(ctx, topic.Topic, &pubsub.Message{Data: v.Message}); err != nil {
			services.Logger.Error("failed to publish to %q for %q - %q", topic, v.Name, err)
			return err
		}

		log.Printf("sent to pubsub topic: %q", topic.Topic)
		return nil
	}

	services.Logger.Error("Invalid output option")
	return errors.New("invalid output option")
}