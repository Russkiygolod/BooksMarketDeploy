command to build docker image api-gw:
    docker run --net=book_m --rm -d --name api-gw -e HOST=book -p 80:80 api-gw
