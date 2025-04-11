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
  path        = "/tmp/example.txt"
  content     = "Hello, Terraform!"
  permissions = "0644"
}

resource "filesystem_directory" "example_dir" {
  path        = "/tmp/terraform-created-dir"
  permissions = "0755"
}

# ファイルの内容を出力
output "file_content" {
  value = filesystem_file.example.content
}