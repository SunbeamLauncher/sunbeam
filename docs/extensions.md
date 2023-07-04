# Extending Sunbeam

The sunbeam extension system is heavily inspired by [gh](https://cli.github.com). Most of the [documentation](https://docs.github.com/en/github-cli/github-cli/creating-github-cli-extensions) from gh can be applied to sunbeam.

## Using Custom Commands

Use the `sunbeam command browse` commands to browse the sunbeam commands available on github.
To install an extension, you can use the `sunbeam command add <name> <url>` command.

The `subeam command manage` command can be used to manage installed extensions.
Alternatively, use the `list`, `remove`, `upgrade` and `rename` commands directly.

## Writing Custom Commands

Any directory containing a `sunbeam-command` executable is a valid sunbeam command.

To test your extension, use the `sunbeam run ./sunbeam-command` command, or just the shorthand `sunbeam run .`.
You can install the current directory as an extension using the `sunbeam command install <alias> .` command.

> **Warning**: Installing local extension is not yet supported on windows.

You can write command using any language. If you want to distribute your command, make sure that you provide instructions on how to install the required dependencies.

Here are some suggestions if you don't know what language to use:

- Bash is already installed on most systems. Sunbeam provides multiple commands to help you write bash commands. You can use the `sunbeam query` command to generate/manipulate JSON objects.
- If you are more confortable with javascript/typescript, take a look at [deno](https://deno.land/). Types are available both on [npm](https://npmjs.com/package/sunbeam-types).

## Distributing Commands

You have multiple alternatives to distribute sunbeam commands:

### Github Gists

The easiest way to write a sunbeam command is to create a [github gist](https://gist.github.com/) containing a `sunbeam-command` executable.

All the examples from the docs are available as gists.

### Git Repositories

Alternatively, you can create a github repository containing the `sunbeam-command` executable in it's root directory, and push it to github.

### Raw Scripts

## Publishing Extensions

Edit the `catalog/catalog.txt` file in this repository to add your extension to the catalog. Then, create a pull request describing your command.
