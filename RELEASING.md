# Releasing

Releases are automated via [Release Please](https://github.com/googleapis/release-please).

## How it works

1. Merge PRs to `master` using [conventional commits](https://www.conventionalcommits.org/):
   - `feat: ...` — triggers a **minor** version bump
   - `fix: ...` — triggers a **patch** version bump
   - `feat!: ...` or `BREAKING CHANGE:` in the footer — triggers a **major** bump
   - `chore:`, `docs:`, `test:` — no release triggered

2. Release Please automatically opens (or updates) a **Release PR** titled `chore: release vX.Y.Z`.
   This PR bumps `version/version.go` and updates `CHANGELOG.md`.

3. When the Release PR is merged, Release Please creates a `vX.Y.Z` tag and a **draft** GitHub
   Release. It stays a draft until the assets are uploaded, because the Terraform Registry
   only ingests a version when its release is published.

4. The Release Please workflow then calls [`release.yml`](.github/workflows/release.yml), and
   GoReleaser publishes signed multi-platform binaries to the GitHub Release, which the
   Terraform Registry picks up. (The tag push alone does not trigger `release.yml`, because
   tags created with `GITHUB_TOKEN` don't start other workflows.)

## Do NOT manually push `v*` tags

The only way a release should be created is by merging the Release PR.
Manual tags bypass the CHANGELOG update and version file bump.

## Re-publishing assets for an existing release

If a release was tagged but its assets were never published, run the `release` workflow
manually against the tag:

```bash
gh workflow run release.yml --ref master -f tag=vX.Y.Z
```

## Checking the current version

```bash
cat version/version.go
```

## Legacy release process (deprecated)

The old `make release` / `scripts/release.sh` process is **deprecated** and
should not be used. It is kept in place only as a historical reference.
