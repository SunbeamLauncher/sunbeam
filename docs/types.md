## PageSchema

This file was automatically generated by json-schema-to-typescript.
DO NOT MODIFY IT BY HAND. Instead, modify the source JSONSchema file,
and run json-schema-to-typescript to regenerate this file.

**POSSIBLE VALUES**

- [List](#list)
- [Detail](#detail)

## Action

**POSSIBLE VALUES**

- object
  - `type`: `'fetch-url'` - The type of the action.
  - `title`: string - The title of the action.
  - `key`: string - The key used as a shortcut.
  - `url`: string - The url to fetch.
  - `method`: `'GET'` | `'POST'` | `'PUT'` | `'DELETE'` - The method to use when fetching.
  - `body`: string - The body to send when fetching.
  - `headers`: object - The headers to send when fetching.
    - `__index`: any
  - `inputs`: [Input](#input)[] - The inputs to show when the action is run.
- object
  - `type`: `'copy-text'` - The type of the action.
  - `title`: string - The title of the action.
  - `text`: string - The text to copy.
  - `key`: string - The key used as a shortcut.
- object
  - `type`: `'open-file'` - The type of the action.
  - `title`: string - The title of the action.
  - `key`: string - The key used as a shortcut.
  - `path`: string - The path to open.
- object
  - `type`: `'open-url'` - The type of the action.
  - `title`: string - The title of the action.
  - `key`: string - The key used as a shortcut.
  - `url`: string - The url to open.
- object
  - `type`: `'run-command'` - The type of the action.
  - `title`: string - The title of the action.
  - `key`: string - The key used as a shortcut.
  - `command`: string - The command to run.
  - `dir`: string - The directory where the command should be run.
  - `onSuccess`: `'reload'` | `'exit'` | `'push'` - The action to take when the command succeeds.
  - `inputs`: [Input](#input)[] - The inputs to show when the action is run.
- object
  - `type`: `'read-file'` - The type of the action.
  - `title`: string - The title of the action.
  - `key`: string - The key used as a shortcut.
  - `path`: string - The path to read.

## Input

**POSSIBLE VALUES**

- object
  - `name`: string - The name of the input.
  - `title`: string - The title of the input.
  - `type`: `'textfield'` - The type of the input.
  - `placeholder`: string - The placeholder of the input.
  - `default`: string - The default value of the input.
  - `secure`: boolean - Whether the input should be secure.
- object
  - `name`: string - The name of the input.
  - `title`: string - The title of the input.
  - `type`: `'checkbox'` - The type of the input.
  - `default`: boolean - The default value of the input.
  - `label`: string - The label of the input.
  - `trueSubstitution`: string - The text substitution to use when the input is true.
  - `falseSubstitution`: string - The text substitution to use when the input is false.
- object
  - `name`: string - The name of the input.
  - `title`: string - The title of the input.
  - `type`: `'textarea'` - The type of the input.
  - `placeholder`: string - The placeholder of the input.
  - `default`: string - The default value of the input.
- object
  - `name`: string - The name of the input.
  - `title`: string - The title of the input.
  - `type`: `'dropdown'` - The type of the input.
  - `items`: object[] - The items of the input.
- `title`: string - The title of the item.
- `value`: string - The value of the item.
  - `default`: string - The default value of the input.

## Preview

The preview to show in the detail view.

**POSSIBLE VALUES**

- object
  - `text`: string - The text to show in the preview.
  - `language`: string - The language of the preview text.
- object
  - `command`: string - The command used to generate the preview.
  - `dir`: string - The directory where the command should be run.
  - `language`: string - The language of the preview text.

## List

**PROPERTIES**

- `type`: `'list'` - The type of the response.
- `title`: string - The title of the page.
- `emptyText`: string - The text to show when the list is empty.
- `showPreview`: boolean - Whether to show details on the right side of the list.
- `actions`: [Action](#action)[] - The global actions attached to the list.
- `items`: [Listitem](#listitem)[] - The items in the list.

## Listitem

**PROPERTIES**

- `title`: string - The title of the item.
- `id`: string - The id of the item.
- `subtitle`: string - The subtitle of the item.
- `preview`: [Preview](#preview)
- `accessories`: string[] - The accessories to show on the right side of the item.
- `actions`: [Action](#action)[] - The actions attached to the item.

## Detail

A detail view displayign a preview and actions.

**PROPERTIES**

- `type`: `'detail'` - The type of the response.
- `title`: string - The title of the page.
- `preview`: [Preview](#preview)
- `actions`: [Action](#action)[] - The actions attached to the detail view.