# Terraform Filesystem Provider

This Terraform provider allows you to create and manage files and directories on your local filesystem.

## Features

- Create, update, and delete files
- Create and delete directories
- Manage permissions for files and directories

## Usage

### Provider Configuration

```hcl
terraform {
  required_providers {
    filesystem = {
      source = "jedipunkz/filesystem"
      version = "0.1.0"
    }
  }
}

provider "filesystem" {}
```

### Creating a File

```hcl
resource "filesystem_file" "example" {
  path        = "/tmp/example.txt"
  content     = "Hello, Terraform!"
  permissions = "0644"  # Optional, defaults to "0644"
}
```

### Creating a Directory

```hcl
resource "filesystem_directory" "example_dir" {
  path        = "/tmp/terraform-created-dir"
  permissions = "0755"  # Optional, defaults to "0755"
}
```

## Building the Provider

To build the provider:

```bash
go build -o terraform-provider-filesystem
```

## Testing Locally

To test the provider locally, add the following to your `~/.terraformrc` file:

```
provider_installation {
  dev_overrides {
    "jedipunkz/filesystem" = "/path/to/your/terraform-provider-filesystem"
  }
  direct {}
}
```

## Complete Workflow

Here's how to use the provider from start to finish:

1. **Build the provider**:
   ```bash
   go build -o terraform-provider-filesystem
   ```

2. **Configure Terraform to use your local provider**:
   Create or update your `~/.terraformrc` file with:
   ```
   provider_installation {
     dev_overrides {
       "jedipunkz/filesystem" = "/path/to/your/terraform-provider-filesystem"
     }
     direct {}
   }
   ```

3. **Create a Terraform configuration file** (main.tf):
   ```hcl
   terraform {
     required_providers {
       filesystem = {
         source = "jedipunkz/filesystem"
         version = "0.1.0"
       }
     }
   }

   provider "filesystem" {}

   resource "filesystem_file" "example" {
     path    = "/tmp/example.txt"
     content = "Hello, Terraform!"
   }
   ```

4. **Initialize your Terraform workspace**:
   ```bash
   terraform init
   ```

5. **Plan your changes**:
   ```bash
   terraform plan
   ```

6. **Apply your changes**:
   ```bash
   terraform apply
   ```

7. **Verify the results**:
   ```bash
   cat /tmp/example.txt
   ```

8. **Clean up resources when done**:
   ```bash
   terraform destroy
   ```

## Examples

For more detailed examples, check the [examples](./examples) directory.