package runtimecontract

import runtimecontractv1 "github.com/uvwt/agentdock-protocol/runtimecontract/v1"

// CurrentVersion is the Runtime wire contract emitted by the current AgentDock implementation.
// Consumers should bind to the concrete runtimecontract/vN package they actually support instead of
// treating this moving value as proof that they understand a newer contract.
const CurrentVersion = runtimecontractv1.Version
