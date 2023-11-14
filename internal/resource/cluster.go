package resource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	consoleClient "github.com/pluralsh/console-client-go"
	"github.com/samber/lo"
)

func Cluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: clusterCreate,
		ReadContext:   clusterRead,
		UpdateContext: clusterUpdate,
		DeleteContext: clusterDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"handle": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cloud": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"byok"}, true),
			},
		},
	}
}

func clusterCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*consoleClient.Client)

	clusterAttrs := consoleClient.ClusterAttributes{Name: d.Get("name").(string)}

	if handle := d.Get("handle").(string); handle != "" {
		clusterAttrs.Handle = &handle
	}

	cluster, err := client.CreateCluster(ctx, clusterAttrs)
	if err != nil {
		return diag.FromErr(err)
	}

	cloud := d.Get("cloud")
	if cloud == "byok" {
		if cluster.CreateCluster.DeployToken == nil {
			return diag.Errorf("could not fetch deploy token from cluster")
		}

		// deployToken := *cluster.CreateCluster.DeployToken
		// url := fmt.Sprintf("%s/ext/gql", p.ConsoleClient.Url())
		// p.doInstallOperator(url, deployToken)
	}

	d.SetId(cluster.CreateCluster.ID)

	return clusterRead(ctx, d, m)
}

func clusterRead(_ context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*consoleClient.Client)
	cluster, err := client.GetCluster(context.Background(), lo.ToPtr(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", cluster.Cluster.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("handle", cluster.Cluster.Handle)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func clusterUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	return clusterRead(ctx, d, m)
}

func clusterDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*consoleClient.Client)
	_, err := client.DeleteCluster(ctx, d.Id())
	return diag.FromErr(err)
}
