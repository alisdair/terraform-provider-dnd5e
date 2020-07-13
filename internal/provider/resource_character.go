package provider

import (
	"errors"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var xpLevels = []int{
	0,
	300,
	900,
	2700,
	6500,
	14000,
	23000,
	34000,
	48000,
	64000,
	85000,
	100000,
	120000,
	140000,
	165000,
	195000,
	225000,
	265000,
	305000,
	355000,
}

var abilities = []string{
	"strength",
	"dexterity",
	"constitution",
	"intelligence",
	"wisdom",
	"charisma",
}

func resourceCharacter() *schema.Resource {
	return &schema.Resource{
		Create: resourceCharacterCreate,
		Read:   resourceCharacterRead,
		Update: resourceCharacterUpdate,
		Delete: schema.RemoveFromState,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"class": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"barbarian",
						"bard",
						"cleric",
						"druid",
						"fighter",
						"monk",
						"paladin",
						"ranger",
						"rogue",
						"sorcerer",
						"warlock",
						"wizard",
					},
					true,
				),
			},
			"level": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"alignment": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"lawful good",
						"neutral good",
						"chaotic good",
						"lawful neutral",
						"neutral",
						"chaotic neutral",
						"lawful evil",
						"neutral evil",
						"chaotic evil",
					},
					true,
				),
			},
			"experience_points": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"strength": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 30),
			},
			"strength_modifier": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dexterity": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 30),
			},
			"dexterity_modifier": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"constitution": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 30),
			},
			"constitution_modifier": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"intelligence": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 30),
			},
			"intelligence_modifier": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"wisdom": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 30),
			},
			"wisdom_modifier": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"charisma": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 30),
			},
			"charisma_modifier": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"proficiency_bonus": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"armor_class": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"speed": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"inventory_item": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"armor_class": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"equipped": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"number": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
					},
				},
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.ComputedIf("level", diffHasChange("experience_points")),
			customdiff.ComputedIf("proficiency_bonus", diffHasChange("experience_points")),
			customdiff.ComputedIf("strength_modifier", diffHasChange("strength")),
			customdiff.ComputedIf("dexterity_modifier", diffHasChange("dexterity")),
			customdiff.ComputedIf("constitution_modifier", diffHasChange("constitution")),
			customdiff.ComputedIf("intelligence_modifier", diffHasChange("intelligence")),
			customdiff.ComputedIf("wisdom_modifier", diffHasChange("wisdom")),
			customdiff.ComputedIf("charisma_modifier", diffHasChange("charisma")),
			customdiff.ComputedIf("armor_class", func(d *schema.ResourceDiff, meta interface{}) bool {
				return d.HasChange("dexterity") || d.HasChange("inventory_item")
			}),
			customdiff.ValidateChange("experience_points", func(old, new, meta interface{}) error {
				if new.(int) < old.(int) {
					return errors.New("experience points cannot decrease")
				}
				return nil
			}),
		),
	}
}

func resourceCharacterRead(d *schema.ResourceData, meta interface{}) error {
	// This seems weird but we're only updating the computed things
	return resourceCharacterUpdate(d, meta)
}

func resourceCharacterCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	d.SetId(name)

	return resourceCharacterUpdate(d, meta)
}

func resourceCharacterUpdate(d *schema.ResourceData, meta interface{}) error {
	xp := d.Get("experience_points").(int)
	level := 0
	for _, xpLimit := range xpLevels {
		level++
		if xp < xpLimit {
			break
		}
	}
	d.Set("level", level)

	proficiencyBonus := int64(math.Ceil(float64(level) * 1.25))
	d.Set("proficiency_bonus", proficiencyBonus)

	for _, ability := range abilities {
		value := d.Get(ability).(int)
		modifier := int64(math.Floor(float64(value-10) / 2))
		d.Set(ability+"_modifier", modifier)
	}

	armorClass := 0
	if inventory, ok := d.GetOk("inventory_item"); ok {
		inventoryList := inventory.([]interface{})
		for i := range inventoryList {
			item := inventoryList[i].(map[string]interface{})
			if equipped, ok := item["equipped"].(bool); ok && !equipped {
				continue
			}
			if armor, ok := item["armor_class"].(int); ok {
				armorClass += armor
			}
		}
	}
	if armorClass == 0 {
		dexterity := d.Get("dexterity_modifier").(int)
		armorClass = dexterity + 10
	}
	if armorClass < 0 {
		armorClass = 0
	}
	d.Set("armor_class", armorClass)

	return nil
}

func diffHasChange(name string) func(*schema.ResourceDiff, interface{}) bool {
	return func(d *schema.ResourceDiff, meta interface{}) bool {
		return d.HasChange(name)
	}
}
