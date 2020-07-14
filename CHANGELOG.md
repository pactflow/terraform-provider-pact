# Changelog
Do this to generate your change history

    git log --pretty=format:'  * [%h](https://github.com/pact-foundation/pact-go/commit/%h) - %s (%an, %ad)' vX.Y.Z..HEAD | egrep -v "wip(:|\()" | grep -v "docs(" | grep -v "chore(" | grep -v Merge | grep -v "test("

<a name="0.0.1"></a>

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
