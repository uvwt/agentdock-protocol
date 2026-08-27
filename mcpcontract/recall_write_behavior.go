package mcpcontract

type RecallWriteOutcome string

const (
	RecallWritePreview  RecallWriteOutcome = "preview"
	RecallWriteMutation RecallWriteOutcome = "mutation"
	RecallWriteError    RecallWriteOutcome = "error"
	RecallWriteReadOnly RecallWriteOutcome = "read_only"
)

type RecallWriteBehaviorCase struct {
	Name              string
	Target            string
	Action            string
	Path              string
	Confirmed         bool
	DryRun            bool
	Expected          RecallWriteOutcome
	Existing          bool
	OverwriteSemantic bool
}

// RecallWriteBehaviorCases defines cross-entrypoint observable semantics. DryRun is
// an internal safety switch and is intentionally not part of the public MCP schema.
func RecallWriteBehaviorCases() []RecallWriteBehaviorCase {
	return []RecallWriteBehaviorCase{
		{Name: "card plan previews", Target: "card", Action: "plan", Expected: RecallWritePreview},
		{Name: "card create unconfirmed previews", Target: "card", Action: "create", Expected: RecallWritePreview},
		{Name: "card create confirmed mutates", Target: "card", Action: "create", Confirmed: true, Expected: RecallWriteMutation},
		{Name: "card create confirmed dry run previews", Target: "card", Action: "create", Confirmed: true, DryRun: true, Expected: RecallWritePreview},
		{Name: "markdown plan previews", Target: "markdown", Action: "plan", Path: "recall/docs/inbox/plan.md", Confirmed: true, Expected: RecallWritePreview},
		{Name: "markdown inbox create unconfirmed mutates", Target: "markdown", Action: "create", Path: "recall/docs/inbox/create.md", Expected: RecallWriteMutation},
		{Name: "markdown create never overwrites existing", Target: "markdown", Action: "create", Path: "recall/docs/inbox/existing.md", Confirmed: true, Existing: true, Expected: RecallWriteError},
		{Name: "markdown protected create unconfirmed errors", Target: "markdown", Action: "create", Path: "profile.md", Expected: RecallWriteError},
		{Name: "markdown protected create confirmed mutates", Target: "markdown", Action: "create", Path: "profile.md", Confirmed: true, Expected: RecallWriteMutation},
		{Name: "markdown replace unconfirmed previews", Target: "markdown", Action: "replace", Path: "recall/docs/inbox/existing.md", Existing: true, Expected: RecallWritePreview, OverwriteSemantic: true},
		{Name: "markdown replace confirmed mutates", Target: "markdown", Action: "replace", Path: "recall/docs/inbox/existing.md", Confirmed: true, Existing: true, Expected: RecallWriteMutation, OverwriteSemantic: true},
		{Name: "markdown append unconfirmed previews", Target: "markdown", Action: "append", Path: "recall/docs/inbox/existing.md", Existing: true, Expected: RecallWritePreview, OverwriteSemantic: true},
		{Name: "markdown append confirmed mutates", Target: "markdown", Action: "append", Path: "recall/docs/inbox/existing.md", Confirmed: true, Existing: true, Expected: RecallWriteMutation, OverwriteSemantic: true},
		{Name: "markdown patch unconfirmed previews", Target: "markdown", Action: "patch", Path: "recall/docs/inbox/existing.md", Existing: true, Expected: RecallWritePreview, OverwriteSemantic: true},
		{Name: "markdown patch confirmed mutates", Target: "markdown", Action: "patch", Path: "recall/docs/inbox/existing.md", Confirmed: true, Existing: true, Expected: RecallWriteMutation, OverwriteSemantic: true},
		{Name: "markdown update fact unconfirmed previews", Target: "markdown", Action: "update_fact", Path: "recall/docs/inbox/existing.md", Existing: true, Expected: RecallWritePreview, OverwriteSemantic: true},
		{Name: "markdown update fact confirmed mutates", Target: "markdown", Action: "update_fact", Path: "recall/docs/inbox/existing.md", Confirmed: true, Existing: true, Expected: RecallWriteMutation, OverwriteSemantic: true},
		{Name: "markdown diff is read only", Target: "markdown", Action: "diff", Path: "recall/docs/inbox/existing.md", Existing: true, Expected: RecallWriteReadOnly},
		{Name: "markdown delete unconfirmed errors", Target: "markdown", Action: "delete", Path: "recall/docs/inbox/existing.md", Existing: true, Expected: RecallWriteError},
		{Name: "markdown delete confirmed mutates", Target: "markdown", Action: "delete", Path: "recall/docs/inbox/existing.md", Confirmed: true, Existing: true, Expected: RecallWriteMutation},
		{Name: "markdown replace dry run previews", Target: "markdown", Action: "replace", Path: "recall/docs/inbox/existing.md", Confirmed: true, DryRun: true, Existing: true, Expected: RecallWritePreview, OverwriteSemantic: true},
		{Name: "markdown create dry run previews", Target: "markdown", Action: "create", Path: "recall/docs/inbox/dry-run.md", Confirmed: true, DryRun: true, Expected: RecallWritePreview},
	}
}
