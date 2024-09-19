# Devdocs

```sh
#!/bin/bash

set -euo pipefail

if [ $# -eq 1 ]; then
  curl "https://devdocs.io/docs/$SLUG/index.json" | jq --arg slug "$1" '.entries | map({
    title: .name,
    subtitle: .type,
    actions: [
      {title: "Open in Browser", type: "open", url: "https://devdocs.io/\($slug)/\(.path)", exit: true},
      {title: "Copy URL", key: "c", type: "copy", text: "https://devdocs.io/\($slug)/\(.path)", exit: true}
    ]
  }) | { type: "list", items: . }'

  exit 0
fi


curl https://devdocs.io/docs/docs.json | jq 'map({
    title: .name,
    subtitle: (.release // "latest"),
    accessories: [ .slug ],
    actions: [
      {
        title: "Browse entries",
        type: "push",
        args: [ .slug ]
      }
    ]
  }) | { type: "list", items: . }'
```
