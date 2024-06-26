// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMain(t *testing.T) {
	tc := testutil.SystemTest(t)
	env := map[string]string{"GOOGLE_CLOUD_PROJECT": tc.ProjectID}
	scope := fmt.Sprintf("projects/%s", tc.ProjectID)
	fullResourceName := fmt.Sprintf("//cloudresourcemanager.googleapis.com/projects/%s", tc.ProjectID)

	m := testutil.BuildMain(t)
	defer m.Cleanup()

	if !m.Built() {
		t.Errorf("failed to build app")
	}

	// Create a bucket in GCS.
	ctx := context.Background()
	bucketName := testutil.TestBucket(ctx, t, tc.ProjectID, "for-assets")
	uri := fmt.Sprintf("gs://%s/client_library_obj", bucketName)

	stdOut, stdErr, err := m.Run(env, 2*time.Minute, fmt.Sprintf("--scope=%s", scope), fmt.Sprintf("--fullResourceName=%s", fullResourceName), fmt.Sprintf("--uri=%s", uri))

	if err != nil {
		t.Errorf("execution failed: %v", err)
	}
	if len(stdErr) > 0 {
		t.Errorf("did not expect stderr output, got %d bytes: %s", len(stdErr), string(stdErr))
	}
	got := string(stdOut)
	if !strings.Contains(got, "operation completed successfully") {
		t.Errorf("stdout returned %s, wanted to contain %s", got, uri)
	}
}
