
# start the environment
make env-start
# building all the binaries
make build-all
# running services
./bin/file-emitter
./bin/publisher --force-scan & 
DB_USER=stats DB_NAME=stats POSTGRES_ENV_PASSWORD=pgpwd4stats ./bin/consumer &

echo "Sleeping for 10 seconds to handle everything"
sleep 10
# running psql to query all the data in the database
count=$(docker exec -it streaming-pipeline-postgres-1 psql -d stats -U stats -t -c "SELECT COUNT(*) from statistics;")
echo "----------------------"
echo "Found records: $count"

ps -ef | grep publisher | grep -v grep | awk '{print $2}' | xargs kill
ps -ef | grep consumer | grep -v grep | awk '{print $2}' | xargs kill

# check
if [[ $count == "0" ]]; then
  echo "Fetched count should not be 0"
  exit 1
fi

#stop env
make env-stop

# remove folder with files
rm -rf /tmp/fileemitter
  
