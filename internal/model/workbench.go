package model

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
)

type Workbench struct {
	Id            types.String            `tfsdk:"id"`
	Name          types.String            `tfsdk:"name"`
	Description   types.String            `tfsdk:"description"`
	SystemPrompt  types.String            `tfsdk:"system_prompt"`
	ProjectID     types.String            `tfsdk:"project_id"`
	RepositoryID  types.String            `tfsdk:"repository_id"`
	AgentRuntime  types.String            `tfsdk:"agent_runtime"`
	Configuration *WorkbenchConfiguration `tfsdk:"configuration"`
	Skills        *WorkbenchSkills        `tfsdk:"skills"`
	Tools         types.Set               `tfsdk:"tool_ids"`
}

func (in *Workbench) Attributes(agentRuntimeID *string, ctx context.Context, d *diag.Diagnostics) (*gqlclient.WorkbenchAttributes, error) {
	tools := make([]*gqlclient.WorkbenchToolAssociationAttributes, len(in.Tools.Elements()))
	elements := make([]WorkbenchTool, len(in.Tools.Elements()))
	d.Append(in.Tools.ElementsAs(ctx, &elements, false)...)
	for i, tool := range elements {
		tools[i] = &gqlclient.WorkbenchToolAssociationAttributes{ToolID: tool.Id.ValueString()}
	}

	return &gqlclient.WorkbenchAttributes{
		Name:             in.Name.ValueStringPointer(),
		Description:      in.Description.ValueStringPointer(),
		SystemPrompt:     in.SystemPrompt.ValueStringPointer(),
		ProjectID:        in.ProjectID.ValueStringPointer(),
		RepositoryID:     in.RepositoryID.ValueStringPointer(),
		AgentRuntimeID:   agentRuntimeID,
		Configuration:    in.Configuration.Attributes(ctx),
		Skills:           in.Skills.Attributes(ctx),
		ToolAssociations: tools,
	}, nil
}

func (in *Workbench) From(response *gqlclient.WorkbenchFragment, ctx context.Context, d *diag.Diagnostics) {
	in.Id = types.StringValue(response.ID)
	in.Name = types.StringValue(response.Name)
	in.Description = types.StringPointerValue(response.Description)
	in.SystemPrompt = types.StringPointerValue(response.SystemPrompt)
	in.ProjectID = types.StringValue(response.Project.ID)
	in.RepositoryID = types.StringValue(response.Repository.ID)

	if response.AgentRuntime != nil && response.AgentRuntime.Cluster != nil && response.AgentRuntime.Cluster.Handle != nil {
		in.AgentRuntime = types.StringValue(fmt.Sprintf("%s/%s", *response.AgentRuntime.Cluster.Handle, response.AgentRuntime.ID))
	}

	// TODO Configuration and skills.

	in.Tools = in.toolsFrom(response.Tools, in.Tools, ctx, d)
}

func (in *Workbench) toolsFrom(tools []*gqlclient.WorkbenchToolFragment, config types.Set, ctx context.Context, d *diag.Diagnostics) types.Set {
	if len(tools) == 0 {
		// Rewriting config to state to avoid inconsistent result errors.
		// This could happen, for example, when sending "nil" to API and "[]" is returned as a result.
		return config
	}

	toolIDs := lo.Map(tools, func(tool *gqlclient.WorkbenchToolFragment, _ int) *string {
		return &tool.ID
	})

	return common.SetFrom(toolIDs, config, ctx, d)
}

type WorkbenchConfiguration struct {
	Infrastructure *WorkbenchInfrastructure `tfsdk:"infrastructure"`
	Coding         *WorkbenchCoding         `tfsdk:"coding"`
}

func (in *WorkbenchConfiguration) Attributes(ctx context.Context) *gqlclient.WorkbenchConfigurationAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchConfigurationAttributes{
		Infrastructure: in.Infrastructure.Attributes(),
		Coding:         in.Coding.Attributes(ctx),
	}
}

type WorkbenchInfrastructure struct {
	Services   types.Bool `tfsdk:"services"`
	Stacks     types.Bool `tfsdk:"stacks"`
	Kubernetes types.Bool `tfsdk:"kubernetes"`
}

func (in *WorkbenchInfrastructure) Attributes() *gqlclient.WorkbenchInfrastructureAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchInfrastructureAttributes{
		Services:   in.Services.ValueBoolPointer(),
		Stacks:     in.Stacks.ValueBoolPointer(),
		Kubernetes: in.Kubernetes.ValueBoolPointer(),
	}
}

type WorkbenchCoding struct {
	Mode         types.String `tfsdk:"mode"`
	Repositories types.Set    `tfsdk:"repositories"`
}

func (in *WorkbenchCoding) Attributes(ctx context.Context) *gqlclient.WorkbenchCodingAttributes {
	if in == nil {
		return nil
	}

	repositories := make([]types.String, len(in.Repositories.Elements()))
	in.Repositories.ElementsAs(ctx, &repositories, false)

	return &gqlclient.WorkbenchCodingAttributes{
		Mode:         lo.ToPtr(gqlclient.AgentRunMode(in.Mode.ValueString())),
		Repositories: lo.Map(repositories, func(v types.String, _ int) *string { return lo.ToPtr(v.ValueString()) }),
	}
}

type WorkbenchSkills struct {
	Ref   *WorkbenchGitRef `tfsdk:"ref"`
	Files types.Set        `tfsdk:"files"`
}

func (in *WorkbenchSkills) Attributes(ctx context.Context) *gqlclient.WorkbenchSkillsAttributes {
	if in == nil {
		return nil
	}

	files := make([]types.String, len(in.Files.Elements()))
	in.Files.ElementsAs(ctx, &files, false)

	return &gqlclient.WorkbenchSkillsAttributes{
		Ref:   in.Ref.Attributes(),
		Files: lo.Map(files, func(v types.String, _ int) *string { return v.ValueStringPointer() }),
	}
}

type WorkbenchGitRef struct {
	Ref    types.String `tfsdk:"ref"`
	Folder types.String `tfsdk:"folder"`
	Files  types.Set    `tfsdk:"files"`
}

func (in *WorkbenchGitRef) Attributes() *gqlclient.GitRefAttributes {
	if in == nil {
		return nil
	}

	files := make([]types.String, len(in.Files.Elements()))
	in.Files.ElementsAs(context.Background(), &files, false)

	return &gqlclient.GitRefAttributes{
		Ref:    in.Ref.ValueString(),
		Folder: in.Folder.ValueString(),
		Files:  lo.Map(files, func(v types.String, _ int) string { return v.ValueString() }),
	}
}
