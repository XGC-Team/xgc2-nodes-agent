package mcpcall

import (
	"context"
	"testing"
	"time"

	"github.com/lxk36/xgc2-orchestration-core/conformance/nodepack"
	"github.com/lxk36/xgc2-orchestration-core/kernel/canonicaljson"
	protocol "github.com/lxk36/xgc2-orchestration-core/kernel/node"
	"github.com/lxk36/xgc2-orchestration-core/sdk/go/contracts"
	nodesdk "github.com/lxk36/xgc2-orchestration-core/sdk/go/node"
)

func TestNodePackConformance(t *testing.T) {
	executor := New()
	input := map[string]any{"serverRef": "github-mcp", "tool": "issue.create", "arguments": map[string]any{"title": "marker"}}
	digest, _ := canonicaljson.DigestValue(input)
	t0 := time.Date(2026, 8, 13, 2, 0, 0, 0, time.UTC)
	request := contracts.NodeInvocationRequest{
		SchemaVersion: protocol.InvocationSchemaVersion,
		InvocationID:  "inv-1", RunID: "run-1", NodeID: "mcp", TypeRef: executor.Descriptor().TypeRef,
		DescriptorDigest: executor.Descriptor().DescriptorDigest, AttemptID: "att-1", AttemptOrdinal: 1,
		Input: input, InputDigest: digest,
		CapabilityGrants: []contracts.CapabilityGrant{{CapabilityRef: "mcp.invoke", Scope: "project", HandleRef: "grant-1", AuthorizationDigest: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", ExpiresAt: t0.Add(time.Minute)}},
		RequestedAt:      t0, Deadline: t0.Add(time.Minute),
	}
	report, err := nodepack.Validate(context.Background(), nodepack.Suite{PackageRef: "xgc2-nodes-agent", Executors: []nodesdk.Executor{executor}, Cases: []nodepack.Case{{Name: "authorized MCP effect", Executor: executor, Request: request, ExpectedStatus: contracts.NodeResultWaiting}}})
	if err != nil || report.DescriptorCount != 1 {
		t.Fatalf("report = %#v, err %v", report, err)
	}
}
