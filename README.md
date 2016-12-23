# Amazon Rekognition のGoSDK 実装サンプル

## SDK経由での利用初期設定
1. AWSのコンソールからユーザーとグループを作成
アクセスの種類を「プログラムによるアクセス」と設定して作成し、グループの権限は「AmazonRekognitionReadOnlyAccess」だけで事足ります。

2. アクセスキーとシークレットアクセスキーをメモ

3. [awsCLI](https://aws.amazon.com/jp/cli/)を利用して認証設定を追加する
`$ aws configure`

4. SDKの追加
`$ go get github.com/aws/aws-sdk-go`

