package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/coreos/coreos-cloudinit/config/validate"
	"github.com/hashicorp/hil"
	"github.com/hashicorp/hil/ast"
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudconfig() *schema.Resource {
	return &schema.Resource{
		Create: cloudconfigRender,
		Read:   cloudconfigRender,
		Update: cloudconfigRender,
		Delete: cloudconfigDelete,
		Exists: cloudconfigExists,

		Schema: map[string]*schema.Schema{
			"template": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Template for cloud config.",
			},
			"vars": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Default:     make(map[string]interface{}),
				Description: "Template variables for cloud config.",
			},
			"gzip": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Compress output w/ gzip (then base64).",
			},
			"validate": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Validate template output.",
			},
			"rendered": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func cloudconfigRender(d *schema.ResourceData, meta interface{}) error {
	rendered, err := renderTemplate(d.Get("template").(string), d.Get("vars").(map[string]interface{}))
	if err != nil {
		return err
	}

	if d.Get("validate").(bool) {
		err := cloudconfigValidate(rendered)
		if err != nil {
			return err
		}
	}

	if d.Get("gzip").(bool) {
		rendered, err = cloudconfigCompress(rendered)
		if err != nil {
			return err
		}
	}

	d.Set("rendered", rendered)
	d.SetId(hash(rendered))
	return nil
}

func cloudconfigDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func cloudconfigExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	current, ok := d.GetOk("rendered")
	if !ok || current == "" {
		return false, nil
	}

	rendered, err := renderTemplate(d.Get("template").(string), d.Get("vars").(map[string]interface{}))
	if err != nil {
		return false, err
	}
	return hash(rendered) == d.Id(), nil
}

// gzip+base64 encode content
func cloudconfigCompress(content string) (string, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	if _, err := gz.Write([]byte(content)); err != nil {
		return "", fmt.Errorf("failed to gzip data: %s", err)
	}
	if err := gz.Flush(); err != nil {
		return "", fmt.Errorf("failed to gzip data: %s", err)
	}
	if err := gz.Close(); err != nil {
		return "", fmt.Errorf("failed to gzip data: %s", err)
	}

	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

// validates the content as cloud-config
func cloudconfigValidate(content string) error {
	report, err := validate.Validate([]byte(content))
	if err != nil {
		return err
	}
	if len(report.Entries()) > 0 {
		return fmt.Errorf("validation failed: %s", report.Entries())
	}
	return nil
}

// taken from https://github.com/hashicorp/terraform/blob/master/builtin/providers/template/datasource_template_file.go
func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}

// taken from https://github.com/hashicorp/terraform/blob/master/builtin/providers/template/datasource_template_file.go
// there, the function is called 'execute'
func renderTemplate(s string, vars map[string]interface{}) (string, error) {
	root, err := hil.Parse(s)
	if err != nil {
		return "", err
	}

	varmap := make(map[string]ast.Variable)
	for k, v := range vars {
		// As far as I can tell, v is always a string.
		// If it's not, tell the user gracefully.
		s, ok := v.(string)
		if !ok {
			return "", fmt.Errorf("unexpected type for variable %q: %T", k, v)
		}
		varmap[k] = ast.Variable{
			Value: s,
			Type:  ast.TypeString,
		}
	}

	cfg := hil.EvalConfig{
		GlobalScope: &ast.BasicScope{
			VarMap:  varmap,
			FuncMap: config.Funcs(),
		},
	}

	result, err := hil.Eval(root, &cfg)
	if err != nil {
		return "", err
	}
	if result.Type != hil.TypeString {
		return "", fmt.Errorf("unexpected output hil.Type: %v", result.Type)
	}

	return result.Value.(string), nil
}
