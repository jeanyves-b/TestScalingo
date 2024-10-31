#! /bin/bash

sudo docker compose up > out.log 2> err.log &

end=$((SECONDS+20))
counter=0

while [ $SECONDS -lt $end ]; do
    curl "localhost:5000/getFiltered?language=Java" > /dev/null 2> /dev/null &
    counter=$((counter+1))
done

echo counter: $counter