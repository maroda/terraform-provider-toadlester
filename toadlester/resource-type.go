package toadlester

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceType is a collection of Env Var types
// These are used to configure each ToadLester looping series
// Currently the only read type is a direct Algo: /series/TYPE/ALGO
// Using the list below, an example:
// Type: name = EXP_SIZE
// Setting: value = 500
// Algorithm: algo = up
/*
	EXP_SIZE=5 EXP_LIMIT=250 EXP_TAIL=1 EXP_MOD=250
	FLOAT_SIZE=4 FLOAT_LIMIT=100 FLOAT_TAIL=5 FLOAT_MOD=1.123
	INT_SIZE=10 INT_LIMIT=100 INT_TAIL=1 INT_MOD=1
	RAND_SIZE=1 RAND_LIMIT=500 RAND_TAIL=3 RAND_MOD=500
*/

func dataSourceType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTypeRead, // Read is at /current/json
		Schema: map[string]*schema.Schema{
			"config": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The returned configuration",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Toadlester API Endpoint",
			},
		},
	}
}

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
		return diag.FromErr(err)
	}

	// Confirmation should contain all three settings
	if !strings.Contains(done, envvar) || !strings.Contains(done, algo) || !strings.Contains(done, setval) {
		return diag.Errorf("API response missing expected values: got '%q' expecting '%q', '%q', '%q'", done, envvar, algo, setval)
	}

	// Set schema fields
	if err = d.Set("name", envvar); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("value", setval); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("algo", algo); err != nil {
		return diag.FromErr(err)
	}

	// Set the ID
	now := time.Now().Format("20060102T150405")
	tag := envvar + "_" + now
	d.SetId(tag)

	return nil
}

func resourceTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	// Run and validate response
	done, err := config.Client.ReadType(set)
	if err != nil {
		return diag.FromErr(err)
	}

	// Unmarshal full json config response
	var currConfig map[string]string
	if err = json.Unmarshal([]byte(done), &currConfig); err != nil {
		return diag.FromErr(err)
	}

	// Drift Detection
	// Get actual value from the API response
	actualVal, ok := currConfig[envvar]
	if !ok {
		return diag.Errorf("configuration missing expected key: '%s'", envvar)
	}

	// Set schema fields with actual values
	if err = d.Set("name", envvar); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("value", actualVal); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("algo", algo); err != nil {
		return diag.FromErr(err)
	}

	// Set the ID
	now := time.Now().Format("20060102T150405")
	tag := envvar + "_" + now
	d.SetId(tag)

	return nil
}

func resourceTypeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceTypeCreate(ctx, d, m)
}

func resourceTypeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceTypeCreate(ctx, d, m)
}

func dataSourceTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Provider config
	config := m.(*Config)

	// Data resource has no settings (for future use)
	set := &Setting{}

	// Run and validate response
	done, err := config.Client.ReadType(set)
	if err != nil {
		return diag.FromErr(err)
	}

	// Unmarshal full json config response
	var currConfig map[string]string
	if err = json.Unmarshal([]byte(done), &currConfig); err != nil {
		return diag.FromErr(err)
	}

	// Set schema fields
	if err = d.Set("config", currConfig); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("endpoint", config.Client.BaseURL); err != nil {
		return diag.FromErr(err)
	}

	// Set the ID
	now := time.Now().Format("20060102T150405")
	tag := "config" + "_" + now
	d.SetId(tag)

	return nil
}
