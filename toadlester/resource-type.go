package toadlester

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceType is the collection of Env Vars used to configure
// Currently the only read type is a direct Algo: /series/TYPE/ALGO
// Using the list below, an example:
// name = EXP_SIZE
// value = 500
// algo = up
/*
	EXP_SIZE=5
	EXP_LIMIT=250
	EXP_TAIL=1
	EXP_MOD=250
	FLOAT_SIZE=4
	FLOAT_LIMIT=100
	FLOAT_TAIL=5
	FLOAT_MOD=1.123
	INT_SIZE=10
	INT_LIMIT=100
	INT_TAIL=1
	INT_MOD=1
	RAND_SIZE=1
	RAND_LIMIT=500
	RAND_TAIL=3
	RAND_MOD=500
*/
func resourceType() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTypeCreate, // Create is at /reset/TYPE_VAR/VALUE
		ReadContext:   resourceTypeRead,   // Read is at /series/type/[up,down]
		UpdateContext: resourceTypeUpdate, // Update is like Create (updates to a new sequence)
		DeleteContext: resourceTypeDelete, // Delete is like Create (deletes the current sequence)

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 20),
				Description:  "The type setting to change",
			},
			"value": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 20),
				Description:  "New value to set the type",
			},
			"algo": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 4),
				Description:  "Algorithm used to read the type",
			},
		},
	}
}

func resourceTypeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diag diag.Diagnostics

	// Provider config
	config := m.(*Config)

	// Resource data
	envvar := d.Get("name").(string)
	setval := d.Get("value").(string)
	algo := d.Get("algo").(string)

	set := &Setting{
		Name:  envvar,
		Value: setval,
		Algo:  algo,
	}

	// Create a new Type sequence
	done, err := config.Client.CreateType(set)
	if err != nil {
		return diag
	}
	if !strings.Contains(done, envvar) && !strings.Contains(done, algo) && !strings.Contains(done, setval) {
		return diag
	}

	// Set the ID
	d.SetId(envvar)

	return diag
}

func resourceTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diag diag.Diagnostics

	// Provider config
	config := m.(*Config)

	// Resource data
	envvar := d.Get("name").(string)
	setval := d.Get("value").(string)
	algo := d.Get("algo").(string)

	set := &Setting{
		Name:  envvar,
		Value: setval,
		Algo:  algo,
	}

	done, err := config.Client.ReadType(set)
	if err != nil {
		return diag
	}
	if !strings.Contains(done, algo) {
		return diag
	}

	return diag
}

func resourceTypeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceTypeCreate(ctx, d, m)
}

func resourceTypeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceTypeCreate(ctx, d, m)
}
