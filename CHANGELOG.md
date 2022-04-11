# Changelog

Do this to generate your change history

    git log --pretty=format:'  * [%h](https://github.com/pact-foundation/pact-go/commit/%h) - %s (%an, %ad)' vX.Y.Z..HEAD | egrep -v "wip(:|\()" | grep -v "docs(" | grep -v "chore(" | grep -v Merge | grep -v "test("


### v0.7.0 (11 April 2022)
  * [d492103](https://github.com/pactflow/terraform/commit/d492103) - feat: administrator support for team resource (Matt Fellows, Mon Apr 11 18:41:13 2022 +1000)
  * [e3894e4](https://github.com/pactflow/terraform/commit/e3894e4) - fix: sort user uuids alphanumerically. Fixes #21 (Matt Fellows, Mon Apr 11 11:38:33 2022 +1000)

### v0.6.0 (08 April 2022)
  * [6beb9a0](https://github.com/pactflow/terraform/commit/6beb9a0) - feat: add main_branch property to applications (Matt Fellows, Fri Apr 8 16:15:38 2022 +1000)

### v0.5.1 (21 March 2022)
  * [9cd0aa5](https://github.com/pactflow/terraform/commit/9cd0aa5) - fix: release under go 1.16 (Matt Fellows, Mon Mar 21 22:01:49 2022 +1100)

### v0.5.0 (21 March 2022)
  * [80ccad5](https://github.com/pactflow/terraform/commit/80ccad5) - feat: support 'contract_requiring_verification_published' event in webhook resource (Matt Fellows, Mon Mar 21 20:47:10 2022 +1100)
  * [b8b2032](https://github.com/pactflow/terraform/commit/b8b2032) - docs: update changelog (Matt Fellows, Tue Mar 8 17:08:22 2022 +1100)
### v0.4.3 (08 March 2022)

- [6b5ae7e](https://github.com/pactflow/terraform/commit/6b5ae7e) - Feat/environments (#20) (Matt Fellows, Tue Mar 8 13:52:05 2022 +1100)
  <a name="0.0.1"></a>

### v0.3.3 (17 January 2022)

- [f254291](https://github.com/pactflow/terraform/commit/f254291) - fix: escape URLs in client API calls, fixes #18 (Matt Fellows, Mon Jan 17 21:56:07 2022 +1100)

### v0.3.2 (29 June 2021)

- [43e3e10](https://github.com/pactflow/terraform/commit/43e3e10) - fix: missing line join in makefile (Matt Fellows, Mon Jun 28 21:55:51 2021 +1000)

### v0.3.1 (28 June 2021)

- [78e23d5](https://github.com/pactflow/terraform/commit/78e23d5) - test: run full create, update and delete acceptance cycle (Matt Fellows, Mon Jun 28 21:25:34 2021 +1000)
- [a75eb07](https://github.com/pactflow/terraform/commit/a75eb07) - fix: team read should use UUID (Matt Fellows, Mon Jun 28 21:23:32 2021 +1000)

### v0.3.0 (28 June 2021)

- [e8e9bac](https://github.com/pactflow/terraform/commit/e8e9bac) - fix: bring back bender system user (Matt Fellows, Mon Jun 28 16:27:08 2021 +1000)
- [6750113](https://github.com/pactflow/terraform/commit/6750113) - feat: add secrets to team and webhooks resources (Matt Fellows, Mon Jun 28 15:55:56 2021 +1000)
- [2b79880](https://github.com/pactflow/terraform/commit/2b79880) - fix: incorrect type in auth settings resource (Matt Fellows, Tue Jun 22 10:20:57 2021 +1000)
- [db8b7af](https://github.com/pactflow/terraform/commit/db8b7af) - test: add authentication settings pact test (Matt Fellows, Tue Jun 22 09:31:09 2021 +1000)
- [3e4f914](https://github.com/pactflow/terraform/commit/3e4f914) - fix: webhook resource does not return ID (Matt Fellows, Tue Jun 22 09:30:37 2021 +1000)
- [454d035](https://github.com/pactflow/terraform/commit/454d035) - feat: add pact tests 🤝 (Matt Fellows, Sat Jun 19 15:22:38 2021 +1000)
- [344b6e4](https://github.com/pactflow/terraform/commit/344b6e4) - docs: update readme with latest docs (Matt Fellows, Mon Feb 22 15:54:33 2021 +1100)

### v0.2.0 (22 February 2021)

- [fa89fed](https://github.com/pactflow/terraform/commit/fa89fed) - fix: correct accept headers (Matt Fellows, Mon Feb 22 15:43:13 2021 +1100)
- [9ae38c7](https://github.com/pactflow/terraform/commit/9ae38c7) - docs: add gh action badge (Matt Fellows, Mon Feb 22 14:45:42 2021 +1100)
- [f840eb2](https://github.com/pactflow/terraform/commit/f840eb2) - feat: add support for authentication settings (Matt Fellows, Mon Feb 22 14:39:22 2021 +1100)
- [7645bb2](https://github.com/pactflow/terraform/commit/7645bb2) - docs: fix errant \_ in webhook docs (Matt Fellows, Sun Jan 31 12:38:03 2021 +1100)

### v0.1.5 (30 January 2021)

- [099d745](https://github.com/pactflow/terraform/commit/099d745) - fix: support strings and JSON webhook bodies. Fixes 12 (Matt Fellows, Fri Jan 29 23:59:30 2021 +1100)
- [c6288d6](https://github.com/pactflow/terraform/commit/c6288d6) - fix: release notes generated incorrect version (Matt Fellows, Fri Jan 22 18:28:40 2021 +1100)

### v0.1.4 (22 January 2021)

Add support for Team, Role and Permissions. Deprecates v1 role and token resource. Other commits include:

- [f102a8f](https://github.com/pactflow/terraform/commit/f102a8f) - fix: oss acceptance test (Matt Fellows, Fri Jan 22 14:15:30 2021 +1100)
- [34ee2fc](https://github.com/pactflow/terraform/commit/34ee2fc) - fix: import re-creating pacticipant. Fixes #11 (Matt Fellows, Fri Jan 22 12:00:58 2021 +1100)

### v0.1.2 (26 October 2020)

- [14ce27e](https://github.com/pactflow/terraform/commit/14ce27e) - docs: reformat to hashicorp directory structure requirements (Matt Fellows, Mon Oct 26 22:10:55 2020 +1100)

### v0.1.1 (26 October 2020)

### v0.1.1 (26 October 2020)

### v0.1.0 (26 October 2020)

- [9125d49](https://github.com/pactflow/terraform/commit/9125d49) - fix: issue with validateEvents where input array was not sorted (Matt Fellows, Mon Oct 26 08:14:20 2020 +1100)

### v0.0.7 (25 October 2020)

- [76cac46](https://github.com/pactflow/terraform/commit/76cac46) - feat: add new provider verification webhook events. Fixes #9 (Matt Fellows, Sun Oct 25 15:24:10 2020 +1100)

### v0.0.6 (25 July 2020)

- [aeccef1](https://github.com/pactflow/terraform/commit/aeccef1) - fix: omitted webhook consumer/provider would always replace resource. Fixes #7 (Matt Fellows, Sat Jul 25 21:39:35 2020 +1000)
- [db0789f](https://github.com/pactflow/terraform/commit/db0789f) - fix: don't send null Roles on user update (Matt Fellows, Thu Jul 23 22:46:58 2020 +1000)

### v0.0.5 (23 July 2020)

- [2b548f6](https://github.com/pactflow/terraform/commit/2b548f6) - fix: temporarily disable user smoke tests (Matt Fellows, Thu Jul 23 22:34:16 2020 +1000)

### v0.0.4 (23 July 2020)

- [c9fffe4](https://github.com/pactflow/terraform/commit/c9fffe4) - fix: allow non-JSON bodies in webhooks. Fixes #6 (Matt Fellows, Thu Jul 23 22:01:02 2020 +1000)

### v0.0.3 (14 July 2020)

- [700428a](https://github.com/pactflow/terraform/commit/700428a) - docs: update user (Matt Fellows, Tue Jul 14 21:53:08 2020 +1000)
- [59b4e97](https://github.com/pactflow/terraform/commit/59b4e97) - docs: how to import resources (Matt Fellows, Tue Jun 30 00:09:26 2020 +1000)
- [9ac3fcc](https://github.com/pactflow/terraform/commit/9ac3fcc) - docs: add note about user resource lifecycle (Matt Fellows, Mon Jun 29 18:42:05 2020 +1000)

### v0.0.3 (16 June 2020)

- [c18adee](https://github.com/pactflow/terraform/commit/c18adee) - docs: make the ALL case clearer for webhooks. #2 (Matt Fellows, Tue Jun 16 20:37:45 2020 +1000)
- [f72169c](https://github.com/pactflow/terraform/commit/f72169c) - Correct contract_content_changed event name (#4) (Garry Jeromson, Tue Jun 16 12:29:55 2020 +0200)
- [8ea5b36](https://github.com/pactflow/terraform/commit/8ea5b36) - docs: add MIT license (Matt Fellows, Mon Mar 16 13:37:08 2020 +1100)

### v0.0.2 (13 March 2020)

### v0.0.2 (13 March 2020)

- [7af0856](https://github.com/pactflow/terraform/commit/7af0856) - test: initialise plugins before acceptance tests (Matt Fellows, Fri Mar 13 17:59:10 2020 +1100)
- [6b3f871](https://github.com/pactflow/terraform/commit/6b3f871) - test: run acceptance tests at the end of CI (Matt Fellows, Fri Mar 13 17:47:02 2020 +1100)
- [307fbab](https://github.com/pactflow/terraform/commit/307fbab) - test: add Pactflow acceptance test (Matt Fellows, Fri Mar 13 17:36:53 2020 +1100)
- [9bbafbc](https://github.com/pactflow/terraform/commit/9bbafbc) - test: add basic OSS acceptance test (Matt Fellows, Fri Mar 13 17:18:05 2020 +1100)
- [4441eea](https://github.com/pactflow/terraform/commit/4441eea) - feat: add support for API tokens (Matt Fellows, Fri Mar 13 16:20:10 2020 +1100)

### v0.0.1 (10 March 2020)

Initial release of the Terraform Pact Provider. Supports OSS and Pactflow brokers and the following resources:

- Pacticipants
- Webhooks
- Secrets
