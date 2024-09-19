# Guidelines

## Choosing a language

Sunbeam extensions are just scripts, so you can use any language you want (as long as it can read and write JSON).

Sunbeam is not aware of the language you are using, so you will have to make sure that your script is executable and that it has the right shebang.

Even though you can use any language, here are some recommendations:

### POSIX Shell

Sunbeam provides multiple helpers to make it easier to share sunbeam extensions, without requiring the user to install additional dependencies (other than sunbeam itself).

- `sunbeam open`: open an url or a file using the default application.
- `sunbeam copy/paste`: copy/paste text from/to the clipboard

```sh
#!/bin/sh

set -eu

jq -n '{ type: "detail", text: "Hello, World!" }'
```

A more complex shell extension can be found [here](./examples/devdocs).

### Deno

[Deno](https://deno.land) is a secure runtime for javascript and typescript. It is an [excellent choice](https://matklad.github.io/2023/02/12/a-love-letter-to-deno.html) for writing scripts that require external dependencies.

Deno allows you to use any npm package by just importing it from a url. This makes it easy to use any library without requiring the user to install it first. The only requirement is that the user already has deno installed.

```ts
#!/usr/bin/env -S deno run -A

console.log(JSON.stringify({
    type: "detail",
    text: "Hello, World!"
}));
```

A more complex typescript extension can be found [here](./examples/hackernews.md).

### Python

If you don't want to use deno/typescript ([you should really give it a try](https://matklad.github.io/2023/02/12/a-love-letter-to-deno.html)), you can use python instead.

Python3 comes preinstalled in macOS and on most linux distributions, so it is a good choice if you want to write an extension that can be used without requiring the user to install additional dependencies.

Make sure to use the `#!/usr/bin/env python3` shebang, as it will make your script more portable.

```python
#!/usr/bin/env python3

import json

print(json.dumps({ "text": "Hello, World!" }))
```

Prefer to not use any external dependencies, as the user will have to install them manually. If you really need to use a dependency, you will need to distribute your extension through pip, and instruct the user how to install it.

See the [file-browser extension](./examples/file-browser.md) for an example.

### Any other language

You can use any language you want, as long as it can write/read JSON to/from stdout/stdin.
Just make sure to use the right shebang, or to compile your script to an binary executable.
