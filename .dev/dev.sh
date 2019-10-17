#!/bin/bash
PGDATAFOLDER=".docker/postgres-data"

# this step is needed to check if user who running the script has no access to folder created by docker
if [ -d $PGDATAFOLDER ]; then

    PGOWNER=$(stat -c '%u' $PGDATAFOLDER)
    CURRUID=$(id -u)
    CURRGID=$(id -g)

    # set docker folder permission if needed after it created / modified inside docker container
    if [ $PGOWNER -ne $CURRUID ]; then
        echo "If prompted, please provide sudo to set correct folder permission for $PGDATAFOLDER, after it created / modified inside docker container"
        sudo chown -R $CURRUID:$CURRGID $PGDATAFOLDER
    fi

fi

docker-compose -f .dev/docker-compose.dev.yml down
docker-compose -f .dev/docker-compose.dev.yml up
