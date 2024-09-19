# Hacker News (Deno)

```ts
#!/usr/bin/env -S deno run -A --ext=ts

import Parser from "npm:rss-parser";
import { formatDistance } from "npm:date-fns";

const feed = await new Parser().parseURL(
    `https://hnrss.org/frontpage?description=0&count=25`
);

const page = {
    type: "list",
    items: feed.items.map((item: any) => ({
        title: item.title || "",
        subtitle: item.categories?.join(", ") || "",
        accessories: item.isoDate
            ? [
                formatDistance(new Date(item.isoDate), new Date(), {
                    addSuffix: true,
                }),
            ]
            : [],
        actions: [
            {
                title: "Open in browser",
                type: "open",
                url: item.link || "",
            },
            {
                title: "Open Comments in Browser",
                type: "open",
                url: item.guid || "",
            },
            {
                title: "Copy Link",
                type: "copy",
                key: "c",
                text: item.link || "",
                exit: true,
            },
        ],
    })),
};

console.log(JSON.stringify(page));
```
