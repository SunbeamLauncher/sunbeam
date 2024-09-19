# File Browser

```python
#!/usr/bin/env python3

import sys
import json
import pathlib

directory = sys.argv[1] if len(sys.argv) > 1 else "."
if directory.startswith("~"):
    directory = directory.replace("~", str(pathlib.Path.home()))
root = pathlib.Path(directory)

items = []
for file in root.iterdir():
    item = {
        "title": file.name,
        "accessories": [str(file.absolute())],
        "actions": [],
    }

    if file.is_dir():
        item["actions"].append(
            {
                "title": "Browse",
                "type": "push",
                "args": [str(file.absolute())],
            }
        )

    item["actions"].extend(
        [
            {
                "title": "Open",
                "key": "o",
                "type": "open",
                "path": str(file.absolute()),
            }
        ]
    )

    items.append(item)

print(json.dumps({"type": "list", "items": items}))
```
