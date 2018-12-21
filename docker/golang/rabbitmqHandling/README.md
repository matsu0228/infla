usage

```
$ cd docker/rabbitmq
$ docker-compose up -d

$ docker ps

CONTAINER ID        IMAGE                       COMMAND                  CREATED             STATUS              PORTS                                                                                        NAMES
b38cd616f869        rabbitmq:3.7.4-management   "docker-entrypoint.s…"   3 days ago          Up 2 minutes        4369/tcp, 5671/tcp, 0.0.0.0:5672->5672/tcp, 15671/tcp, 25672/tcp, 0.0.0.0:15672->15672/tcp   rabbitmq

```

open `http://localhost:15672/` to use rabbitMQ manager with your brawser.