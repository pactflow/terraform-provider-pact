# Changelog
Do this to generate your change history

    git log --pretty=format:'  * [%h](https://github.com/pact-foundation/pact-go/commit/%h) - %s (%an, %ad)' vX.Y.Z..HEAD | egrep -v "wip(:|\()" | grep -v "docs(" | grep -v "chore(" | grep -v Merge | grep -v "test("

<a name="0.0.1"></a>

### v0.3.2 (29 June 2021)
  * [b2ceba8](https://github.com/pactflow/terraform/commit/b2ceba8) - chore: fix makefile to not use bashisms (Matt Fellows, Tue Jun 29 11:21:46 2021 +1000)
  * [43e3e10](https://github.com/pactflow/terraform/commit/43e3e10) - fix: missing line join in makefile (Matt Fellows, Mon Jun 28 21:55:51 2021 +1000)

### v0.3.1 (28 June 2021)
  * [78e23d5](https://github.com/pactflow/terraform/commit/78e23d5) - test: run full create, update and delete acceptance cycle (Matt Fellows, Mon Jun 28 21:25:34 2021 +1000)
  * [a75eb07](https://github.com/pactflow/terraform/commit/a75eb07) - fix: team read should use UUID (Matt Fellows, Mon Jun 28 21:23:32 2021 +1000)

### v0.3.0 (28 June 2021)
  * [7d6f501](https://github.com/pactflow/terraform/commit/7d6f501) - chore: update readme on progress (Matt Fellows, Mon Jun 28 16:35:55 2021 +1000)
  * [e8e9bac](https://github.com/pactflow/terraform/commit/e8e9bac) - fix: bring back bender system user (Matt Fellows, Mon Jun 28 16:27:08 2021 +1000)
  * [1c2206d](https://github.com/pactflow/terraform/commit/1c2206d) - chore: update deps (Matt Fellows, Mon Jun 28 16:16:07 2021 +1000)
  * [77bbe26](https://github.com/pactflow/terraform/commit/77bbe26) - chore: update deps (Matt Fellows, Mon Jun 28 16:07:07 2021 +1000)
  * [794e171](https://github.com/pactflow/terraform/commit/794e171) - chore: update deps (Matt Fellows, Mon Jun 28 16:03:21 2021 +1000)
  * [ecb3822](https://github.com/pactflow/terraform/commit/ecb3822) - chore: update error lint in user resource (Matt Fellows, Mon Jun 28 15:57:59 2021 +1000)
  * [6750113](https://github.com/pactflow/terraform/commit/6750113) - feat: add secrets to team and webhooks resources (Matt Fellows, Mon Jun 28 15:55:56 2021 +1000)
  * [2b79880](https://github.com/pactflow/terraform/commit/2b79880) - fix: incorrect type in auth settings resource (Matt Fellows, Tue Jun 22 10:20:57 2021 +1000)
  * [f3941a6](https://github.com/pactflow/terraform/commit/f3941a6) - chore: tidy go mods (Matt Fellows, Tue Jun 22 10:12:34 2021 +1000)
  * [074afb0](https://github.com/pactflow/terraform/commit/074afb0) - chore: tidy go mods (Matt Fellows, Tue Jun 22 10:00:56 2021 +1000)
  * [42948ca](https://github.com/pactflow/terraform/commit/42948ca) - chore: tidy go modules (Matt Fellows, Tue Jun 22 09:59:58 2021 +1000)
  * [85fa65d](https://github.com/pactflow/terraform/commit/85fa65d) - chore: update deps (Matt Fellows, Tue Jun 22 09:50:23 2021 +1000)
  * [05b9af2](https://github.com/pactflow/terraform/commit/05b9af2) - chore: update deps (Matt Fellows, Tue Jun 22 09:41:50 2021 +1000)
  * [db8b7af](https://github.com/pactflow/terraform/commit/db8b7af) - test: add authentication settings pact test (Matt Fellows, Tue Jun 22 09:31:09 2021 +1000)
  * [945a3d1](https://github.com/pactflow/terraform/commit/945a3d1) - chore: cleanup authentication settings resource (Matt Fellows, Tue Jun 22 09:30:54 2021 +1000)
  * [3e4f914](https://github.com/pactflow/terraform/commit/3e4f914) - fix: webhook resource does not return ID (Matt Fellows, Tue Jun 22 09:30:37 2021 +1000)
  * [454d035](https://github.com/pactflow/terraform/commit/454d035) - feat: add pact tests 🤝 (Matt Fellows, Sat Jun 19 15:22:38 2021 +1000)
  * [344b6e4](https://github.com/pactflow/terraform/commit/344b6e4) - docs: update readme with latest docs (Matt Fellows, Mon Feb 22 15:54:33 2021 +1100)

### v0.2.0 (22 February 2021)
  * [fa89fed](https://github.com/pactflow/terraform/commit/fa89fed) - fix: correct accept headers (Matt Fellows, Mon Feb 22 15:43:13 2021 +1100)
  * [9ae38c7](https://github.com/pactflow/terraform/commit/9ae38c7) - docs: add gh action badge (Matt Fellows, Mon Feb 22 14:45:42 2021 +1100)
  * [f840eb2](https://github.com/pactflow/terraform/commit/f840eb2) - feat: add support for authentication settings (Matt Fellows, Mon Feb 22 14:39:22 2021 +1100)
  * [7645bb2](https://github.com/pactflow/terraform/commit/7645bb2) - docs: fix errant _ in webhook docs (Matt Fellows, Sun Jan 31 12:38:03 2021 +1100)
  * [9796d6d](https://github.com/pactflow/terraform/commit/9796d6d) - chore: add suffixes to all resources in terraform config to avoid conflicts (Matt Fellows, Sun Jan 31 12:32:18 2021 +1100)
  * [7f4ebde](https://github.com/pactflow/terraform/commit/7f4ebde) - chore: update release script to fix changelog (Matt Fellows, Sat Jan 30 08:49:25 2021 +1100)

### v0.1.5 (30 January 2021)
  * [099d745](https://github.com/pactflow/terraform/commit/099d745) - fix: support strings and JSON webhook bodies. Fixes 12 (Matt Fellows, Fri Jan 29 23:59:30 2021 +1100)
  * [c6288d6](https://github.com/pactflow/terraform/commit/c6288d6) - fix: release notes generated incorrect version (Matt Fellows, Fri Jan 22 18:28:40 2021 +1100)
  * [21d68f8](https://github.com/pactflow/terraform/commit/21d68f8) - chore: update changelog for teams/roles/permissions (Matt Fellows, Fri Jan 22 18:27:06 2021 +1100)

### v0.1.4 (22 January 2021)

  Add support for Team, Role and Permissions. Deprecates v1 role and token resource. Other commits include:

  * [c17b759](https://github.com/pactflow/terraform/commit/c17b759) - chore: tidy up error exports (Matt Fellows, Fri Jan 22 18:05:54 2021 +1100)
  * [7c6774f](https://github.com/pactflow/terraform/commit/7c6774f) - chore: remove travis (Matt Fellows, Fri Jan 22 17:56:15 2021 +1100)
  * [27a367b](https://github.com/pactflow/terraform/commit/27a367b) - chore: add github secret to build (Matt Fellows, Fri Jan 22 17:51:38 2021 +1100)
  * [41ce4cd](https://github.com/pactflow/terraform/commit/41ce4cd) - chore: make log dir prior to build running (Matt Fellows, Fri Jan 22 15:23:54 2021 +1100)
  * [c14353b](https://github.com/pactflow/terraform/commit/c14353b) - chore: pin version of terraform for build (Matt Fellows, Fri Jan 22 14:26:41 2021 +1100)
  * [2b464b6](https://github.com/pactflow/terraform/commit/2b464b6) - chore: add GH build for testing (Matt Fellows, Fri Jan 22 14:22:08 2021 +1100)
  * [f102a8f](https://github.com/pactflow/terraform/commit/f102a8f) - fix: oss acceptance test (Matt Fellows, Fri Jan 22 14:15:30 2021 +1100)
  * [62ce5e5](https://github.com/pactflow/terraform/commit/62ce5e5) - chore: add docs and cleanup for teams (Matt Fellows, Fri Jan 22 12:56:59 2021 +1100)
  * [34ee2fc](https://github.com/pactflow/terraform/commit/34ee2fc) - fix: import re-creating pacticipant. Fixes #11 (Matt Fellows, Fri Jan 22 12:00:58 2021 +1100)
  * [6728fd4](https://github.com/pactflow/terraform/commit/6728fd4) - chore: fix oss acceptance build provider config (Matt Fellows, Mon Oct 26 22:12:37 2020 +1100)

### v0.1.2 (26 October 2020)
  * [14ce27e](https://github.com/pactflow/terraform/commit/14ce27e) - docs: reformat to hashicorp directory structure requirements (Matt Fellows, Mon Oct 26 22:10:55 2020 +1100)
  * [e4e4a16](https://github.com/pactflow/terraform/commit/e4e4a16) - chore: update pactflow.tf acceptance config for 0.0.13 (Matt Fellows, Mon Oct 26 21:44:14 2020 +1100)

### v0.1.1 (26 October 2020)
  * [0d28e51](https://github.com/pactflow/terraform/commit/0d28e51) - chore: remove travis release step in favour of GH Actions (Matt Fellows, Mon Oct 26 17:23:29 2020 +1100)
  * [da0b265](https://github.com/pactflow/terraform/commit/da0b265) - chore: use GH actions to publish to Terraform registry (Matt Fellows, Mon Oct 26 17:00:15 2020 +1100)

### v0.1.1 (26 October 2020)
  * [da0b265](https://github.com/pactflow/terraform/commit/da0b265) - chore: use GH actions to publish to Terraform registry (Matt Fellows, Mon Oct 26 17:00:15 2020 +1100)

### v0.1.0 (26 October 2020)
  * [9125d49](https://github.com/pactflow/terraform/commit/9125d49) - fix: issue with validateEvents where input array was not sorted (Matt Fellows, Mon Oct 26 08:14:20 2020 +1100)

### v0.0.7 (25 October 2020)
  * [76cac46](https://github.com/pactflow/terraform/commit/76cac46) - feat: add new provider verification webhook events. Fixes #9 (Matt Fellows, Sun Oct 25 15:24:10 2020 +1100)

### v0.0.6 (25 July 2020)
  * [aeccef1](https://github.com/pactflow/terraform/commit/aeccef1) - fix: omitted webhook consumer/provider would always replace resource. Fixes #7 (Matt Fellows, Sat Jul 25 21:39:35 2020 +1000)
  * [db0789f](https://github.com/pactflow/terraform/commit/db0789f) - fix: don't send null Roles on user update (Matt Fellows, Thu Jul 23 22:46:58 2020 +1000)

### v0.0.5 (23 July 2020)
  * [2b548f6](https://github.com/pactflow/terraform/commit/2b548f6) - fix: temporarily disable user smoke tests (Matt Fellows, Thu Jul 23 22:34:16 2020 +1000)

### v0.0.4 (23 July 2020)
  * [c9fffe4](https://github.com/pactflow/terraform/commit/c9fffe4) - fix: allow non-JSON bodies in webhooks. Fixes #6 (Matt Fellows, Thu Jul 23 22:01:02 2020 +1000)

### v0.0.3 (14 July 2020)
  * [700428a](https://github.com/pactflow/terraform/commit/700428a) - docs: update user (Matt Fellows, Tue Jul 14 21:53:08 2020 +1000)
  * [dd5f834](https://github.com/pactflow/terraform/commit/dd5f834) - chore: differentiate 40x errors (Matt Fellows, Tue Jul 14 21:36:53 2020 +1000)
  * [59b4e97](https://github.com/pactflow/terraform/commit/59b4e97) - docs: how to import resources (Matt Fellows, Tue Jun 30 00:09:26 2020 +1000)
  * [6715549](https://github.com/pactflow/terraform/commit/6715549) - chore: update docs (Matt Fellows, Mon Jun 29 23:21:27 2020 +1000)
  * [9ac3fcc](https://github.com/pactflow/terraform/commit/9ac3fcc) - docs: add note about user resource lifecycle (Matt Fellows, Mon Jun 29 18:42:05 2020 +1000)

### v0.0.3 (16 June 2020)
 * [c18adee](https://github.com/pactflow/terraform/commit/c18adee) - docs: make the ALL case clearer for webhooks. #2 (Matt Fellows, Tue Jun 16 20:37:45 2020 +1000)
 * [f72169c](https://github.com/pactflow/terraform/commit/f72169c) - Correct contract_content_changed event name (#4) (Garry Jeromson, Tue Jun 16 12:29:55 2020 +0200)
 * [1dfb448](https://github.com/pactflow/terraform/commit/1dfb448) - chore: update issue template (Matt Fellows, Mon Mar 16 13:41:25 2020 +1100)
 * [8ea5b36](https://github.com/pactflow/terraform/commit/8ea5b36) - docs: add MIT license (Matt Fellows, Mon Mar 16 13:37:08 2020 +1100)

### v0.0.2 (13 March 2020)


### v0.0.2 (13 March 2020)
 * [b77f3cf](https://github.com/pactflow/terraform/commit/b77f3cf) - chore: run docker early on in the build so its available by the time acceptance tests run (Matt Fellows, Fri Mar 13 18:06:55 2020 +1100)
 * [7af0856](https://github.com/pactflow/terraform/commit/7af0856) - test: initialise plugins before acceptance tests (Matt Fellows, Fri Mar 13 17:59:10 2020 +1100)
 * [6b3f871](https://github.com/pactflow/terraform/commit/6b3f871) - test: run acceptance tests at the end of CI (Matt Fellows, Fri Mar 13 17:47:02 2020 +1100)
 * [307fbab](https://github.com/pactflow/terraform/commit/307fbab) - test: add Pactflow acceptance test (Matt Fellows, Fri Mar 13 17:36:53 2020 +1100)
 * [9bbafbc](https://github.com/pactflow/terraform/commit/9bbafbc) - test: add basic OSS acceptance test (Matt Fellows, Fri Mar 13 17:18:05 2020 +1100)
 * [4441eea](https://github.com/pactflow/terraform/commit/4441eea) - feat: add support for API tokens (Matt Fellows, Fri Mar 13 16:20:10 2020 +1100)

### v0.0.1 (10 March 2020)

Initial release of the Terraform Pact Provider. Supports OSS and Pactflow brokers and the following resources:

* Pacticipants
* Webhooks
* Secrets
