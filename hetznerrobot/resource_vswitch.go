package hetznerrobot

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVSwitchCreate,
		ReadContext:   resourceVSwitchRead,
		UpdateContext: resourceVSwitchUpdate,
		DeleteContext: resourceVSwitchDelete,

		Importer: &schema.ResourceImporter{
			State: resourceVSwitchImportState,
		},

		Schema: map[string]*schema.Schema{
			"vswitch_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "VSwitch ID",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "vSwitch name",
			},
			"vlan": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "VLAN ID",
			},
			// computed / read-only fields
			"is_cancelled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Cancellation status",
			},
			"servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Attached server list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"server_ipv6_net": {
							Type:     schema.TypeString,
							Required: true,
						},
						"server_number": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"status": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Attached subnet list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mask": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"gateway": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"cloud_networks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Attached cloud network list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mask": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"gateway": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceVSwitchImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	c := meta.(HetznerRobotClient)

	vSwitchID, _ := d.Get("vswitch_id").(int)
	vSwitch, err := c.getVSwitch(vSwitchID)
	if err != nil {
		return nil, fmt.Errorf("Unable to find VSwitch with ID %d:\n\t %q", vSwitchID, err)
	}

	d.Set("name", vSwitch.Name)
	d.Set("vlan", vSwitch.Vlan)
	d.Set("is_cancelled", vSwitch.IsCancelled)
	d.Set("servers", vSwitch.Server)
	d.Set("subnets", vSwitch.Subnet)
	d.Set("cloud_networks", vSwitch.CloudNetwork)
	d.Set("vswitch_id", vSwitchID)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceVSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	name := d.Get("name").(string)
	vlan := d.Get("vlan").(int)
	vSwitch, err := c.createVSwitch(name, vlan)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Unable to create VSwitch :\n\t %q", err))
	}

	d.Set("is_cancelled", vSwitch.IsCancelled)
	d.Set("servers", vSwitch.Server)
	d.Set("subnets", vSwitch.Subnet)
	d.Set("cloud_networks", vSwitch.CloudNetwork)
	d.Set("vswitch_id", vSwitch.Id)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceVSwitchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	vSwitchID, _ := d.Get("vswitch_id").(int)
	vSwitch, err := c.getVSwitch(vSwitchID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Unable to find VSwitch with ID %d:\n\t %q", vSwitchID, err))
	}

	d.Set("name", vSwitch.Name)
	d.Set("vlan", vSwitch.Vlan)
	d.Set("cancelled", vSwitch.IsCancelled)
	d.Set("servers", vSwitch.Server)
	d.Set("subnets", vSwitch.Subnet)
	d.Set("cloud_networks", vSwitch.CloudNetwork)
	d.Set("vswitch_id", vSwitchID)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceVSwitchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	vSwitchID := d.Get("vswitch_id").(int)
	name := d.Get("name").(string)
	vlan := d.Get("vlan").(int)
	vSwitch, err := c.updateVSwitch(vSwitchID, name, vlan)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Unable to update VSwitch :\n\t %q", err))
	}

	d.Set("is_cancelled", vSwitch.IsCancelled)
	d.Set("servers", vSwitch.Server)
	d.Set("subnets", vSwitch.Subnet)
	d.Set("cloud_networks", vSwitch.CloudNetwork)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceVSwitchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	vSwitchID, _ := d.Get("vswitch_id").(int)
	err := c.deleteVSwitch(vSwitchID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Unable to find VSwitch with ID %d:\n\t %q", vSwitchID, err))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
