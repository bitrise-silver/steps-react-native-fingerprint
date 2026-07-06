# React Native Build Fingerprint

A [Bitrise Step](https://devcenter.bitrise.io/) that computes a deterministic SHA-256 fingerprint of your dependency files and exports it so consequent `restore-cache` / `save-cache` steps can skip the native/bundle build when only JavaScript changed.

This is a lightweight, do-it-yourself pattern. For a fully managed, compilation-level remote cache across Gradle, Xcode (LLVM CAS) and C++, see [Bitrise Build Cache for React Native](https://bitrise.io/platform/build-cache/react-native).

## Inputs

| Key | Default | Description |
| --- | --- | --- |
| `file_paths` | `package.json`<br>`package-lock.json` | Newline-separated files whose contents determine the fingerprint. Blank lines and `#` comments are ignored. |
| `key_prefix` | _(empty)_ | Namespace prepended to the fingerprint (joined with `-`) to form `BUNDLE_HASH_STRING`. |
| `verbose` | `false` | Log the list of files being fingerprinted. |

## Outputs

| Env var | Description |
| --- | --- |
| `BUNDLE_HASH_STRING` | The dependency fingerprint, prefixed with `key_prefix` if set. Always set. Use as the key for `restore-cache` / `save-cache` steps. |
| `BUNDLE_CACHE_FOUND` | `"true"` if a Bitrise key-value cache entry exists for `BUNDLE_HASH_STRING`, `"false"` otherwise. Only checks existence — the cache entry is not downloaded. |


Runnable end-to-end examples:

- Bare React Native — <https://github.com/bitrise-silver/react-native-build-skip-demo>
- Expo — <https://github.com/bitrise-silver/expo-build-skip-demo>

## Development

```bash
go vet ./...
go build ./...
go test ./... -v
```

The hashing logic is pure functions with no Bitrise dependencies, so the unit tests run anywhere. Output export uses `envman`, which is present on Bitrise stacks. The `BUNDLE_CACHE_FOUND` check calls Bitrise's key-value cache service directly, so it needs the `BITRISEIO_ABCS_API_URL` / `BITRISEIO_BITRISE_SERVICES_ACCESS_TOKEN` secrets that are only available in an actual Bitrise CI build.

## License

[MIT](LICENSE)
