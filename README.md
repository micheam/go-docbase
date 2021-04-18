# go-docbase

[DocBase API](https://help.docbase.io/groups/1472) Go言語向け クライアントライブラリ

## Installation

このパッケージは、go getコマンドでインストールできます:

```
go get github.com/micheam/go-docbase
```

`go doc` コマンドからドキュメントを参照してください:

```
go doc github.com/micheam/go-docbase
```

## Usage

記事の一覧抽出:

```go
var (
    ctx    = context.Background()
    domain = os.Getenv("DOCBASE_DOMAIN")
    param  = url.Values{}
)
docbase.SetToken(os.Getenv("DOCBASE_TOKEN"))
posts, meta, _ := docbase.ListPosts(ctx, domain, param)
for i := range posts {
 e   fmt.Println(posts[i].Title)
}
```

記事詳細の取得:

```go
var (
    ctx    = context.Background()
    domain = os.Getenv("DOCBASE_DOMAIN")
    postID = docbase.PostID(1863830) // 記事ID
)
docbase.SetToken(os.Getenv("DOCBASE_TOKEN"))
post, _ := docbase.GetPost(ctx, domain, postID)
fmt.Printf("%d:%s\n%s", post.ID, post.Title, post.Body)
```

その他の例については、 [examples](./examples/) を参照してみてください。

## TODO

|    | API Method                        | Endpoint                              |
|----|-----------------------------------|---------------------------------------|
|    | 所属チーム取得API                 | https://help.docbase.io/posts/92977   |
|    | ユーザ検索API                     | https://help.docbase.io/posts/680809  |
| ✅ | メモの検索API                     | https://help.docbase.io/posts/92984   |
| ✅ | メモの投稿API                     | https://help.docbase.io/posts/92980   |
| ✅ | メモの詳細取得API                 | https://help.docbase.io/posts/97204   |
| ✅ | メモの更新API                     | https://help.docbase.io/posts/92981   |
|    | メモのアーカイブAPI               | https://help.docbase.io/posts/665804  |
|    | メモのアーカイブ解除API           | https://help.docbase.io/posts/665806  |
|    | メモの削除API                     | https://help.docbase.io/posts/92982   |
|    | コメント投稿API                   | https://help.docbase.io/posts/216289  |
|    | コメント削除API                   | https://help.docbase.io/posts/216290  |
|    | ファイルアップロードAPI           | https://help.docbase.io/posts/225804  |
|    | ファイルダウンロードAPI           | https://help.docbase.io/posts/1084833 |
| ✅ | タグの取得API                     | https://help.docbase.io/posts/92979   |
|    | グループ作成API                   | https://help.docbase.io/posts/652985  |
|    | グループ検索API                   | https://help.docbase.io/posts/92978   |
|    | グループ詳細取得API               | https://help.docbase.io/posts/652983  |
|    | グループへのユーザー追加API       | https://help.docbase.io/posts/665797  |
|    | グループからユーザーを削除するAPI | https://help.docbase.io/posts/665799  |

## License
[MIT](./LICENSE)

## Author
micheam <michto.maeda@gmail.com>
