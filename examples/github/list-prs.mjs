#!/usr/bin/env zx

let repo = argv.repo;

const res = await $`gh api /repos/${repo}/pulls`;
const prs = JSON.parse(res);
const items = prs.map((pr) => ({
  title: pr.title,
  subtitle: pr.user.login,
  actions: [
    {
      shortcut: "enter",
      type: "open-url",
      title: "Open in Browser",
      url: pr.html_url,
    },
    {
      shortcut: "ctrl+y",
      type: "copy",
      title: "Copy URL",
      content: pr.html_url,
    },
  ],
}));

for (const item of items) {
  console.log(JSON.stringify(item));
}
