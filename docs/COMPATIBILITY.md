# Compatibility notes

## Go and dependencies

The module requires Go 1.26 or newer. Use a maintained patch release.
`golang.org/x/crypto` is updated from its 2022 revision to v0.58.0, with
`golang.org/x/sys` v0.49.0. CI checks Go 1.26 and the stable release.

## Struct validation

`structs.ValidateAndModify` now returns errors for malformed tags, unknown
rules, invalid expressions, incompatible field types, invalid bounds, and
integer absolute-value overflow. Register custom validators before use; a nil
callback unregisters a rule. Registration is synchronized, and callbacks run
outside the registry lock.

Rule declaration order and significant parameter whitespace are preserved.
Regex groups, classes, quantifiers, escapes, and closed `\Q...\E` literals are
supported in combined tags. Escape literal commas outside regex groups,
classes, or quantifiers. An unclosed `\Q` literal is rejected as ambiguous.

Tagged private fields return errors; untagged private fields are skipped.
`required` nested structs and pointers validate their child rules, with cycles
visited once per call. `first_upper` handles the first Unicode rune. Lengths,
truncation, and masks retain byte-based indexing.

Modifications remain sequential: an earlier modification is retained if a later
rule fails. Check the returned error before using the result.

## Tables and CSV

CSV destinations and column bindings are validated even for header-only input.
Malformed input, unsupported or private fields, duplicate mapped columns,
nil embedded pointers, and numeric overflow return errors. Parsed rows are
appended only after the complete input succeeds; errors leave the destination
unchanged. Numeric parsing respects the destination's bit width.

Rendering rejects inaccessible fields instead of panicking. `RenderCSV` flushes
its writer before returning, including for small tables. HTML escaping remains
enabled.

## Files

`file.Copy` rejects source/destination aliases, including hard links and
symlinks, before truncation. Copying remains non-atomic.
`ReadLines` returns scanner errors, including lines over the scanner limit.
`WriteLines` concatenates strings without adding separators and reports write,
flush, sync, and close errors.

Delimiter readers retain their rune-based positions and existing behavior:
`ReadBlocks` starts a block at each delimiter and omits an unterminated final
block. `Each` closes directory handles before callbacks and recursion and does
not follow directory symlinks.

## SSH

Use `SecureShellWithHostKey` with a trusted known-hosts file. Callers own both
returned objects and must close the client and session. Connection setup,
authentication, and session creation share a 15-second deadline; active
sessions have no deadline. Unbracketed IPv6 hosts are supported.

The deprecated `SecureShell` helper still skips host-key verification. Its
returned session now remains connected until closed; closing the session also
closes its owning connection.

## Certificates and keys

`ssl.Generate` creates self-signed RSA-2048 server leaf certificates without CA
or certificate-signing rights. IP names become IP subject alternative names.
Applications requiring a CA need a separate certificate-authority flow.

Outputs must be distinct regular-file paths. Symlink and nonregular outputs,
hard-link aliases, and case-only filenames in the same directory are rejected.
Temporary files are prepared before generating the key. Replacements use mode
0644 for certificates and 0600 for private keys, including existing keys with
broader permissions. These are POSIX permission semantics.

Each output is replaced by rename, subject to OS rename semantics. Updating the
pair is not globally atomic: a later failure can leave a mismatched pair. Handle
errors before deploying or reloading the files.

`ssl.Check` verifies regular-file existence and returns filesystem errors. It
does not validate trust, expiration, or whether the certificate and key match.
