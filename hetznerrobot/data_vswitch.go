package hetznerrobot

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataVSwitch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSwitchRead,
		Schema: map[string]*schema.Schema{
			"vswitch_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "VSwitch ID",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "vSwitch name",
			},
			"vlan": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "VLAN ID",
			},
			"is_cancelled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Cancellation status",
			},
			"servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Attached server list",
				Elem: ,
			},
			"subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Attached subnet list",
				Elem: ,
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

func dataSourceVSwitchRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(HetznerRobotClient)

	vSwitchID := d.Get("vswitch_id").(int)
	vSwitch, err := c.getVSwitch(vSwitchID)
	if err != nil {
		return fmt.Errorf("Unable to find VSwitch with ID %d:\n\t %q", vSwitchID, err)
	}

	d.Set("name", vSwitch.Name)
	d.Set("vlan", vSwitch.Vlan)
	d.Set("is_cancelled", vSwitch.IsCancelled)
	d.Set("servers", vSwitch.Server)
	d.Set("subnets", vSwitch.Subnet)
	d.Set("cloud_networks", vSwitch.CloudNetwork)
	d.Set("vswitch_id", vSwitchID)

	return nil
}
