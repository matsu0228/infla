# docker sanbox

## simple docker

* docker.app実行
```
docker run -d -p 80:80 --name nginx-sandbox nginx

docker container ls -a
docker image ls
```
* open [http://localhost:80](http://localhost:80) at your broser

* clean server
```
docker container rm **
docker image rm **
```

## docker compose

* あとでやる：go + redis
* ref
* https://qiita.com/TsutomuNakamura/items/7e90e5efb36601c5bc8a
* https://qiita.com/pottava/items/452bf80e334bc1fee69a