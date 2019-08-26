## setup

```
- see latest version: https://www.terraform.io/downloads.html

curl -sL https://releases.hashicorp.com/terraform/0.12.6/terraform_0.12.6_darwin_amd64.zip > terraform.zip

unzip terraform.zip

# set path for terraform, then check command

terraform --version
```

## basic comand

```
terraform init
=> Terraform作業ディレクトリの初期化

terraform plan
=> 実行計画作成.
=> 既にTerraformでコード化された環境ならば、実際のリソースやパラメータを変更することなく、パラメータレベルで変更点を確認できる

terraform state show <リソース.リソース名>
=> tfstate内の各リソースの状態を確認


terraform apply
=> 設定ファイルに従ってリソース作成

terraform destroy
=>Terraformが管理するリソースを削除

```