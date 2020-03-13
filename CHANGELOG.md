# Changelog
Do this to generate your change history

    git log --pretty=format:'  * [%h](https://github.com/pact-foundation/pact-go/commit/%h) - %s (%an, %ad)' vX.Y.Z..HEAD | egrep -v "wip(:|\()" | grep -v "docs(" | grep -v "chore(" | grep -v Merge | grep -v "test("

<a name="0.0.1"></a>

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
