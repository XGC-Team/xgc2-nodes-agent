package mcpcall

import (
	"context"
	"errors"

	"github.com/lxk36/xgc2-orchestration-core/kernel/canonicaljson"
	protocol "github.com/lxk36/xgc2-orchestration-core/kernel/node"
	"github.com/lxk36/xgc2-orchestration-core/sdk/go/contracts"
)

const packageDigest = "sha256:85fa768ca860a45810117f0157e41068f07db2f18ebc1f7f461f201824559e30"

type Executor struct{ descriptor contracts.NodeDescriptor }

func New() *Executor {
	descriptor := contracts.NodeDescriptor{
		SchemaVersion: protocol.DescriptorSchemaVersion, TypeRef: "xgc.agent.mcp-call/v1", DisplayName: "Authorized MCP call",
		PackageRef: "xgc2-nodes-agent", PackageDigest: packageDigest,
		InputSchema: contracts.Schema{Type: contracts.TypeObject, Properties: map[string]contracts.Schema{
			"serverRef": {Type: contracts.TypeString}, "tool": {Type: contracts.TypeString}, "arguments": {Type: contracts.TypeObject, AdditionalProperties: true},
		}, Required: []string{"serverRef", "tool", "arguments"}},
		OutputSchema: contracts.Schema{Type: contracts.TypeObject, AdditionalProperties: true},
		Mode:         contracts.NodeEffectful, Determinism: contracts.NodeDeterministic,
		RequiredCapabilities: []contracts.CapabilityRequirement{{CapabilityRef: "mcp.invoke", Scope: "project"}},
		AllowedEffectKinds:   []string{"xgc.mcp-call/v1"}, MaxInputBytes: 262144, MaxOutputBytes: 262144,
	}
	descriptor.DescriptorDigest, _ = protocol.DescriptorDigest(descriptor)
	return &Executor{descriptor: descriptor}
}

func (executor *Executor) Descriptor() contracts.NodeDescriptor { return executor.descriptor }

func (executor *Executor) Execute(_ context.Context, request contracts.NodeInvocationRequest) (contracts.NodeResult, error) {
	server, serverOK := request.Input["serverRef"].(string)
	tool, toolOK := request.Input["tool"].(string)
	arguments, argumentsOK := request.Input["arguments"].(map[string]any)
	if !serverOK || !toolOK || !argumentsOK || server == "" || tool == "" {
		return contracts.NodeResult{}, errors.New("MCP server, tool, or arguments are invalid")
	}
	intent := map[string]any{"serverRef": server, "tool": tool, "arguments": arguments}
	intentDigest, err := canonicaljson.DigestValue(intent)
	if err != nil {
		return contracts.NodeResult{}, err
	}
	proposal := contracts.EffectProposal{
		EffectKey: "mcp-call", Kind: "xgc.mcp-call/v1", TargetRef: server,
		IntentSchemaDigest: packageDigest, IntentDigest: intentDigest, Ownership: contracts.EffectAttached,
		CompensationPolicy: contracts.CompensationNone, RequiredCapabilityRefs: []string{"mcp.invoke"},
		PolicyDigest: request.CapabilityGrants[0].AuthorizationDigest, Deadline: request.Deadline,
	}
	evidence, err := canonicaljson.DigestValue(proposal)
	if err != nil {
		return contracts.NodeResult{}, err
	}
	return contracts.NodeResult{
		Status: contracts.NodeResultWaiting, Effects: []contracts.EffectProposal{proposal},
		Wait:           &contracts.NodeWait{Kind: contracts.NodeWaitEffect, SubjectRef: proposal.EffectKey, ConditionDigest: intentDigest},
		EvidenceDigest: evidence,
	}, nil
}
