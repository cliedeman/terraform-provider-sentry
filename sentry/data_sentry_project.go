package sentry

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func dataSourceSentryProject() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Project data source.",

		ReadContext: dataSourceSentryProjectRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the project belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The name of this project.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"slug": {
				Description: "The slug for this project.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"platform": {
				Description: "The optional platform for this project.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"internal_id": {
				Description: "The internal ID for this project.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"color": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"features": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"digests_min_delay": {
				Description: "The minimum amount of time (in seconds) to wait between scheduling digests for delivery after the initial scheduling.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"digests_max_delay": {
				Description: "The maximum amount of time (in seconds) to wait between scheduling digests for delivery.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"resolve_age": {
				Description: "Hours in which an issue is automatically resolve if not seen after this amount of time.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceSentryProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	slug := d.Get("slug").(string)
	org := d.Get("organization").(string)

	tflog.Debug(ctx, "Reading Sentry project", map[string]interface{}{
		"projectSlug": slug,
		"org":         org,
	})
	proj, resp, err := client.Projects.Get(ctx, org, slug)
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read Sentry project", map[string]interface{}{
		"projectSlug": proj.Slug,
		"projectID":   proj.ID,
		"org":         org,
	})

	d.SetId(proj.Slug)
	retErr := multierror.Append(
		d.Set("organization", proj.Organization.Slug),
		d.Set("name", proj.Name),
		d.Set("slug", proj.Slug),
		d.Set("platform", proj.Platform),
		d.Set("internal_id", proj.ID),
		d.Set("is_public", proj.IsPublic),
		d.Set("color", proj.Color),
		d.Set("features", proj.Features),
		d.Set("status", proj.Status),
		d.Set("digests_min_delay", proj.DigestsMinDelay),
		d.Set("digests_max_delay", proj.DigestsMaxDelay),
		d.Set("resolve_age", proj.ResolveAge),
	)

	teams := make([]string, 0, len(proj.Teams))
	for _, team := range proj.Teams {
		teams = append(teams, *team.Slug)
	}
	retErr = multierror.Append(retErr, d.Set("teams", flattenStringSet(teams)))

	// TODO: Project options

	return diag.FromErr(retErr.ErrorOrNil())
}
