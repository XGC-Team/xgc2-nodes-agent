# xgc2-nodes-agent

Open XGC orchestration nodes for agents, models, MCP tools, approvals, memory
policy, and expert collaboration.

Initial node: `xgc.agent.mcp-call/v1` produces a declarative MCP-call Effect and
waits for its receipt. It cannot call MCP directly; the host must durably
prepare the Effect and dispatch it under the declared `mcp.invoke` grant.
