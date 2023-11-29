# Can you handle it Redis?

Simple tool to check redis write and read massively. It will write, read and remove keys continuously.
Each goroutine will write and read a key, and then remove it.

Just run it with the following command:

```
docker run -it --rm ghcr.io/cruizba/can-you-handle-it-redis:latest \
    -cluster=redis-1:6379,redis-2:6379,redis-3:6379 \
    -password=yourpassword \
    -goroutines 10 \
    -sleep 1s
```

- `cluster` is a comma separated list of redis nodes.
- `password` is the redis password.
- `goroutines` is the number of goroutines to use.
- `sleep` is the time to sleep between each write/read.

In case of errors, they will appear in the stout.

I've written this just to check redis faultolerance while starting and removing servers from a redis cluster.

> Note: This can polute your redis cluster, so use it with caution.

