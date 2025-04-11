package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func New() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"filesystem_file":      resourceFile(),
			"filesystem_directory": resourceDirectory(),
		},
	}
}

func resourceFile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFileCreate,
		ReadContext:   resourceFileRead,
		UpdateContext: resourceFileUpdate,
		DeleteContext: resourceFileDelete,

		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The path to the file",
			},
			"content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The content of the file",
				Default:     "",
			},
			"permissions": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0644",
				Description: "File permissions in octal format (e.g., '0644')",
			},
		},
	}
}

func resourceDirectory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDirectoryCreate,
		ReadContext:   resourceDirectoryRead,
		DeleteContext: resourceDirectoryDelete,

		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The path to the directory",
			},
			"permissions": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0755",
				ForceNew:    true,
				Description: "Directory permissions in octal format (e.g., '0755')",
			},
		},
	}
}

func parsePermissions(perm string) (os.FileMode, error) {
	var mode os.FileMode
	_, err := fmt.Sscanf(perm, "%o", &mode)
	if err != nil {
		return 0, fmt.Errorf("invalid permission format: %s", perm)
	}
	return mode, nil
}

func resourceFileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	path := d.Get("path").(string)
	content := d.Get("content").(string)
	permStr := d.Get("permissions").(string)

	// Parse permissions
	perm, err := parsePermissions(permStr)
	if err != nil {
		return diag.FromErr(err)
	}

	// Make sure the directory exists
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating directory %s: %s", dir, err))
	}

	// Write the file
	err = os.WriteFile(path, []byte(content), perm)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error writing file %s: %s", path, err))
	}

	// Generate an ID based on path
	hash := sha256.Sum256([]byte(path))
	d.SetId(hex.EncodeToString(hash[:]))

	return resourceFileRead(ctx, d, meta)
}

func resourceFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	path := d.Get("path").(string)

	// Check if the file exists
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// File was deleted outside of Terraform
			d.SetId("")
			return diags
		}
		return diag.FromErr(fmt.Errorf("error reading file %s: %s", path, err))
	}

	// Ensure it's a file, not a directory
	if fileInfo.IsDir() {
		return diag.FromErr(fmt.Errorf("path %s is a directory, not a file", path))
	}

	// Read the file content
	content, err := os.ReadFile(path)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading file %s: %s", path, err))
	}

	if err := d.Set("content", string(content)); err != nil {
		return diag.FromErr(err)
	}

	// Set permissions
	perm := fmt.Sprintf("%04o", fileInfo.Mode().Perm())
	if err := d.Set("permissions", perm); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceFileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	path := d.Get("path").(string)

	if d.HasChange("content") || d.HasChange("permissions") {
		content := d.Get("content").(string)
		permStr := d.Get("permissions").(string)

		// Parse permissions
		perm, err := parsePermissions(permStr)
		if err != nil {
			return diag.FromErr(err)
		}

		// Write the file with new content and/or permissions
		err = os.WriteFile(path, []byte(content), perm)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error writing file %s: %s", path, err))
		}
	}

	return resourceFileRead(ctx, d, meta)
}

func resourceFileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	path := d.Get("path").(string)

	// Delete the file
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return diag.FromErr(fmt.Errorf("error deleting file %s: %s", path, err))
	}

	// Remove ID from state
	d.SetId("")

	return diags
}

func resourceDirectoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	path := d.Get("path").(string)
	permStr := d.Get("permissions").(string)

	// Parse permissions
	perm, err := parsePermissions(permStr)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create the directory
	err = os.MkdirAll(path, perm)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating directory %s: %s", path, err))
	}

	// Set permissions explicitly in case MkdirAll didn't set them correctly
	err = os.Chmod(path, perm)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting permissions for directory %s: %s", path, err))
	}

	// Generate an ID based on path
	hash := sha256.Sum256([]byte(path))
	d.SetId(hex.EncodeToString(hash[:]))

	return resourceDirectoryRead(ctx, d, meta)
}

func resourceDirectoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	path := d.Get("path").(string)

	// Check if the directory exists
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Directory was deleted outside of Terraform
			d.SetId("")
			return diags
		}
		return diag.FromErr(fmt.Errorf("error reading directory %s: %s", path, err))
	}

	// Ensure it's a directory, not a file
	if !fileInfo.IsDir() {
		return diag.FromErr(fmt.Errorf("path %s is a file, not a directory", path))
	}

	// Set permissions
	perm := fmt.Sprintf("%04o", fileInfo.Mode().Perm())
	if err := d.Set("permissions", perm); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDirectoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	path := d.Get("path").(string)

	// Delete the directory
	err := os.RemoveAll(path)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting directory %s: %s", path, err))
	}

	// Remove ID from state
	d.SetId("")

	return diags
}