command to build docker image book:
    docker run --net=book_m -d --rm --name book -e DB_USER_P=postgres -e DB_PASSWORD_P=postgres -e DB_NAME_P=BookMarket -e HOST_P=postgres -e HOST_R=redis book:1
