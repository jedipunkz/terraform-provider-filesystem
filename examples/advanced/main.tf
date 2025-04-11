terraform {
  required_providers {
    filesystem = {
      source = "jedipunkz/filesystem"
      version = "0.1.0"
    }
  }
}

provider "filesystem" {}

# シンプルなテキストファイルの作成
resource "filesystem_file" "simple_text" {
  path    = "/tmp/simple.txt"
  content = "これはシンプルなテキストファイルです。"
}

# JSONファイルの作成（設定ファイルなど）
resource "filesystem_file" "config_json" {
  path    = "/tmp/config.json"
  content = jsonencode({
    app_name    = "テスト応用"
    environment = "開発"
    database = {
      host     = "localhost"
      port     = 5432
      username = "user"
      password = "password"
    }
  })
}

# スクリプトファイルの作成
resource "filesystem_file" "bash_script" {
  path        = "/tmp/script.sh"
  content     = <<-EOT
    #!/bin/bash
    echo "Terraformで生成されたスクリプト"
    echo "現在時刻: $(date)"
    ls -la /tmp
  EOT
  permissions = "0755"  # 実行権限を付与
}

# メインプロジェクトディレクトリの作成
resource "filesystem_directory" "project_dir" {
  path = "/tmp/my-project"
}

# サブディレクトリの作成（メインディレクトリに依存）
resource "filesystem_directory" "src_dir" {
  path = "${filesystem_directory.project_dir.path}/src"
  depends_on = [filesystem_directory.project_dir]
}

resource "filesystem_directory" "config_dir" {
  path = "${filesystem_directory.project_dir.path}/config"
  depends_on = [filesystem_directory.project_dir]
}

# 設定ファイルをコンフィグディレクトリに配置
resource "filesystem_file" "app_config" {
  path    = "${filesystem_directory.config_dir.path}/app.conf"
  content = <<-EOT
    # アプリケーション設定
    LOG_LEVEL=INFO
    DEBUG=false
    API_PORT=8080
  EOT
  depends_on = [filesystem_directory.config_dir]
}

# 出力
output "script_path" {
  value = filesystem_file.bash_script.path
}

output "project_structure" {
  value = {
    main_dir = filesystem_directory.project_dir.path
    subdirs  = [
      filesystem_directory.src_dir.path,
      filesystem_directory.config_dir.path
    ]
    config_file = filesystem_file.app_config.path
  }
}