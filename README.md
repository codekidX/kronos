## Kronos - Cron as a service (CronAAS)

> Kronos is a minute based cron runtime packed as as service running over HTTP
> for other services to schedule and run cron tasks. This service can be used to
> abstract all your crons and keep in observation under a single service
> umbrella. This abstraction has a lot of uses which are listed below.

You can use this runtime as a package as well, if you think that your service
will need to spawn too many cron tasks dynamically and so much async threads are
not what you need to keep in your service memory. Check out
[Kronos as a package]() post here.

## Use cases:

- **Deferring execution** - you can defer a HTTP request altogether incase any
  failure occurs at the time of first execution
- **No cron code** - you don't have to maintain cron codes on your service
- **Better control** - you can start/stop any cron whenever it is needed from a
  single place
- **Low overhead** - you can have your service's resources/threads put for good
  use

## APIs

There are 2 APIs provided by Kronos to run your tasks efficiently:

- **Defer** - Deferred tasks are one-shot tasks which runs in a minute
  precision.

> Planned for v0.2 release

- **Repeat** - Repeating tasks are tasks which repeats at a given interval.

> **NOTE:** Kronos does not support time based tasks. Example: Run tasks for
> 3minute at every 3rd hour.
