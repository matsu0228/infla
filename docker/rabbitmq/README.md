

* setup
```
docker-compose up -d
```

* management plugin (open with your brawser)
```
http://localhost:15672 
```

* cli
```
# 標準コマンド: `rabbitmqctl`
# ------------------------------
docker exec my-queue rabbitmqctl --help

# vhost
docker exec rabbitmq  rabbitmqctl list_vhosts
docker exec rabbitmq  rabbitmqctl add_vhost test

# ユーザー
docker exec rabbitmq rabbitmqctl list_users

# Management Plugin に付属するツール: `rabbitmqadmin`
# ------------------------------

# 新規 Queue の定義
docker exec my-queue rabbitmqadmin -u admin -p admin -V my-vhost declare queue name=my-q auto_delete=false durable=true
```