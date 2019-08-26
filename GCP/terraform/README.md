## GCE instance managements


### GCP Manager vs Terraform

- 詳しく理解していないが、Terraformのほうが良いらしい
  - クラウドに縛られずに定義も可能(方言は多々あるが)


### prepare

- create account

  - GCP console > IAM > service account
  - Role: Compute Admin (beta), Storage Admin

- download account's secret key
  - json
  - put on `secret/` dir



### architecture


現時点の仮directory

```
gcp
├── backend.tf          # TODO: GCSにtfstateファイルを保存するための設定
├── firewall.tf         # TODO: GCP上のファイアウォールルールを設定
├── network.tf          # TODO: カスタムVPCネットワークを構築する設定
├── provider.tf         # GCPに接続するための資格情報、プロジェクト名、リージョンを設定
└── variables.tf        # 各設定ファイル内で使用する変数を設定

```


- TODO: 下記を比較検討する. 環境ごとの設定切り替え/tfstateを分離ができることがポイントなよう
  - https://qiita.com/shogomuranushi/items/e2f3ff3cfdcacdd17f99#
  - https://qiita.com/eigo_s/items/02264a5a7ad0ff6c5387
