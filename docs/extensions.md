# Extending Subeam

The sunbeam extension system is heavily inspired by [gh](https://cli.github.com). Most of the [documentation](https://docs.github.com/en/github-cli/github-cli/creating-github-cli-extensions) from gh can be applied to sunbeam.

## Using Extensions

Use the `sunbeam extension browse` commands to browse the sunbeam extensions available on github.

To install an extension, you can use the `sunbeam extension install <alias> <url>` command.

For example, to install the `sunbeam-file-browser` extension as `file-browser`:

```bash
# Install the extension
sunbeam extension install file-browser pomdtr/sunbeam-file-browser

# Run the extension
sunbeam file-browser
```

The `subeam extension manage` command can be used to manage installed extensions.

Alternatively, you can use the `list`, `remove` and `upgrade` commands directly.

## Writing Extensions

### Script Extensions

Any directory containing a `sunbeam-extension` executable is a valid sunbeam extension.

To publish an extension, you can create a github repository containing the `sunbeam-extension` executable, and push it to github.

If you want your extension to be listed in the `sunbeam extension browse` command, you can add the `sunbeam-extension` topic to your repository.

### Binary Extensions

Sunbeam also supports binary extensions. In this case, sunbeam will download the binary from a github release instead of cloning the extension repository. The [sunbeam-extension-precompile](https://github.com/pomdtr/sunbeam-extension-precompile) github action can be used to automatically compile and publish your extension as a binary.