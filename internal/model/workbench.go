package model

import (
	"context"
	"fmt"
	"strings"

	"terraform-provider-plural/internal/client"
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

func (in *Workbench) Attributes(client *client.Client, ctx context.Context, d *diag.Diagnostics) (*gqlclient.WorkbenchAttributes, error) {
	agentRuntimeID, err := in.agentRuntimeAttribute(client, ctx)
	if err != nil {
		return nil, err
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
		ToolAssociations: in.toolsAttribute(ctx, d),
	}, nil
}

func (in *Workbench) toolsAttribute(ctx context.Context, d *diag.Diagnostics) []*gqlclient.WorkbenchToolAssociationAttributes {
	if in.Tools.IsNull() {
		return nil
	}

	toolIDs := make([]types.String, len(in.Tools.Elements()))
	d.Append(in.Tools.ElementsAs(ctx, &toolIDs, false)...)

	result := make([]*gqlclient.WorkbenchToolAssociationAttributes, 0, len(toolIDs))
	for _, toolID := range toolIDs {
		if toolID.IsNull() {
			continue
		}

		result = append(result, &gqlclient.WorkbenchToolAssociationAttributes{ToolID: toolID.ValueString()})
	}

	return result
}

func (in *Workbench) agentRuntimeAttribute(client *client.Client, ctx context.Context) (*string, error) {
	ref := in.AgentRuntime.ValueString()

	if lo.IsEmpty(ref) {
		return nil, nil
	}

	split := strings.Split(ref, "/")
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid agent runtime reference: %s", ref)
	}

	clusterHandle, runtimeName := split[0], split[1]
	cluster, err := client.GetClusterByHandle(ctx, lo.ToPtr(clusterHandle))
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster: %s", err.Error())
	}
	if cluster == nil {
		return nil, fmt.Errorf("cluster not found: %s", clusterHandle)
	}

	agentRuntime, err := client.GetAgentRuntimeByName(ctx, runtimeName, cluster.Cluster.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent runtime: %s", err.Error())
	}
	if agentRuntime == nil {
		return nil, fmt.Errorf("agent runtime not found: %s", runtimeName)
	}

	return &agentRuntime.AgentRuntime.ID, nil
}

func (in *Workbench) From(response *gqlclient.WorkbenchFragment, ctx context.Context, d *diag.Diagnostics) {
	if in == nil {
		return
	}
	if response == nil {
		return
	}

	in.Id = types.StringValue(response.ID)
	in.Name = types.StringValue(response.Name)
	in.Description = types.StringPointerValue(response.Description)
	in.SystemPrompt = types.StringPointerValue(response.SystemPrompt)
	in.ProjectID = common.ProjectFrom(response.Project)

	if response.Repository != nil {
		in.RepositoryID = types.StringValue(response.Repository.ID)
	} else {
		in.RepositoryID = types.StringNull()
	}

	if response.AgentRuntime != nil && response.AgentRuntime.Cluster != nil && response.AgentRuntime.Cluster.Handle != nil {
		in.AgentRuntime = types.StringValue(fmt.Sprintf("%s/%s", *response.AgentRuntime.Cluster.Handle, response.AgentRuntime.ID))
	}

	in.Configuration.From(response.Configuration, ctx, d)
	in.Skills.From(response.Skills, ctx, d)

	in.Tools = in.toolsFrom(response.Tools, in.Tools, ctx, d)
}

func (in *Workbench) toolsFrom(tools []*gqlclient.WorkbenchToolFragment, config types.Set, ctx context.Context, d *diag.Diagnostics) types.Set {
	if len(tools) == 0 {
		// Rewriting config to state to avoid inconsistent result errors.
		// This could happen, for example, when sending "nil" to API and "[]" is returned as a result.
		return config
	}

	toolIDs := lo.Map(tools, func(tool *gqlclient.WorkbenchToolFragment, _ int) *string { return &tool.ID })

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

func (in *WorkbenchConfiguration) From(configuration *gqlclient.WorkbenchFragment_Configuration, ctx context.Context, d *diag.Diagnostics) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.Coding.From(configuration.Coding, ctx, d)
	in.Infrastructure.From(configuration.Infrastructure)
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

func (in *WorkbenchInfrastructure) From(configuration *gqlclient.WorkbenchFragment_Configuration_Infrastructure) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.Services = types.BoolPointerValue(configuration.Services)
	in.Stacks = types.BoolPointerValue(configuration.Stacks)
	in.Kubernetes = types.BoolPointerValue(configuration.Kubernetes)
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

func (in *WorkbenchCoding) From(configuration *gqlclient.WorkbenchFragment_Configuration_Coding, ctx context.Context, d *diag.Diagnostics) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	var mode *string
	if configuration.Mode != nil {
		mode = lo.ToPtr(string(*configuration.Mode))
	}

	in.Mode = types.StringPointerValue(mode)
	in.Repositories = common.SetFrom(configuration.Repositories, in.Repositories, ctx, d)
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
		Ref:   in.Ref.Attributes(ctx),
		Files: lo.Map(files, func(v types.String, _ int) *string { return v.ValueStringPointer() }),
	}
}

func (in *WorkbenchSkills) From(response *gqlclient.WorkbenchFragment_Skills, ctx context.Context, d *diag.Diagnostics) {
	if in == nil {
		return
	}
	if response == nil {
		return
	}

	in.Ref.From(response.Ref, ctx, d)
	in.Files = common.SetFrom(response.Files, in.Files, ctx, d)
}

type WorkbenchGitRef struct {
	Ref    types.String `tfsdk:"ref"`
	Folder types.String `tfsdk:"folder"`
	Files  types.Set    `tfsdk:"files"`
}

func (in *WorkbenchGitRef) Attributes(ctx context.Context) *gqlclient.GitRefAttributes {
	if in == nil {
		return nil
	}

	files := make([]types.String, len(in.Files.Elements()))
	in.Files.ElementsAs(ctx, &files, false)

	return &gqlclient.GitRefAttributes{
		Ref:    in.Ref.ValueString(),
		Folder: in.Folder.ValueString(),
		Files:  lo.Map(files, func(v types.String, _ int) string { return v.ValueString() }),
	}
}

func (in *WorkbenchGitRef) From(response *gqlclient.WorkbenchFragment_Skills_Ref, ctx context.Context, d *diag.Diagnostics) {
	if in == nil {
		return
	}
	if response == nil {
		return
	}

	in.Ref = types.StringValue(response.Ref)
	in.Folder = types.StringValue(response.Folder)
	in.Files = common.SetFrom(lo.ToSlicePtr(response.Files), in.Files, ctx, d)
}
