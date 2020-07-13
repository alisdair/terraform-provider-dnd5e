package provider

import (
	"math/rand"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceRoll() *schema.Resource {
	return &schema.Resource{
		Create: resourceRollCreate,
		Read:   schema.Noop,
		Delete: schema.RemoveFromState,

		Schema: map[string]*schema.Schema{
			"reroll": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"number": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"sides": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      20,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(2, 20),
			},
			"modifier": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  0,
			},
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceRollCreate(d *schema.ResourceData, meta interface{}) error {
	number := d.Get("number").(int)
	sides := d.Get("sides").(int)
	modifier := d.Get("modifier").(int)

	values := make([]int, number)
	total := modifier
	for i := 0; i < number; i++ {
		values[i] = rand.Intn(sides) + 1
		total += values[i]
	}

	d.Set("values", values)
	d.Set("total", total)

	d.SetId(strconv.Itoa(total))

	return nil
}
