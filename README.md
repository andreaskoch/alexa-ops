# 「 AlexaOPS 」
The 24/7 endpoint for your "Amazon Echo" based DevOPS skills.

**Deployment**

> You: Alexa, ask OPS to deploy `<application name>`

> Alexa: OK. Deploying version `1.2.3` of `<application name>`.

**Status**

> You: Alexa, ask OPS how `<application name>` is doing?

> Alexa: OK. Give me a few seconds while I am checking the status of `<application name>`.

> Alexa: All KPIs are within normal parameters.
> ◾ I can reach the health check page.
> ◾ We received 15 orders in the last 30 minutes.
> ◾ And no errors can be found in the logs. 

**Restart**

> You: Alexa, ask OPS to restart `<application name>`

> Alexa: OK. Restarting `<application name>`.

**Current version**

> You: Alexa, ask OPS which version of `<application name>` is currently running

> Alexa: `<application name>` is currently running in version `1.2.3`.

**Available version**

> You: Alexa, ask OPS what the latest available version of `<application name>` is.

> Alexa: The latest available version of `<application name>` is `1.3.0`.

## Roadmap

The prototype is missing some crucial parts that still need to be implemented:

**Authorization Verification**

Google Authenticator based authorization checking so not everybody in the room can mess up your production environment.

> Alexa: Before I continue please state your name and authorization code.

> You: My name is `<Andy>` and my authorization code is `<123456>`.

> Alex: OK `<Andy>` I will now start ...`

## License

「 AlexaOPS 」is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.
