#!/bin/bash


GAMEDATA_PATH=$1
OUTPUT_FILE=$2
JWP_CMD='wine64 /home/dn/projects/bl3/jwp/apo-jwp.exe'
WINDEBUG=fixme-all


echo $GAMEDATA_PATH > $OUTPUT_FILE


for MAPFILE in $(find $GAMEDATA_PATH -name '*.umap'); do
  echo "Checking map file $MAPFILE..."
  if strings $MAPFILE | grep -qE 'OakMissionSpawner|OakSpawner'; then
    echo "Found spawner data in $MAPFILE, attempting to serialize"
    $JWP_CMD serialize $(echo $MAPFILE | sed 's/\.umap//g') 2>&1 > /dev/null
    SERIALIZED_FILE=$(echo $MAPFILE | sed 's/\.umap/.json/g')
    if [ ! -e "$SERIALIZED_FILE" ]; then
        echo "WARNING: jwp could not serialize map file $MAPFILE"
    else
        echo "$SERIALIZED_FILE" | tee -a $OUTPUT_FILE
    fi
  else
    echo "No spawner data in $MAPFILE"
  fi
done
