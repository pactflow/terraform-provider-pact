# Changelog
Do this to generate your change history

    git log --pretty=format:'  * [%h](https://github.com/pact-foundation/pact-go/commit/%h) - %s (%an, %ad)' vX.Y.Z..HEAD | egrep -v "wip(:|\()" | grep -v "docs(" | grep -v "chore(" | grep -v Merge | grep -v "test("

<a name="0.0.1"></a>

### v0.0.1 (10 March 2020)


### v0.0.1 (10 March 2020)


### v0.0.1 (10 March 2020)


Initial release of the Terraform Pact Provider. Supports OSS and Pactflow brokers and the following resources:

* Pacticipants
* Webhooks
* Secrets
