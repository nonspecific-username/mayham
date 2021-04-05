#!/bin/bash


GAMEDATA_PATH=${GAMEDATA_PATH}
OUTPUT_FILE=${OUTPUT_FILE}
JWP_SERIALIZE=${JWP_SERIALIZE:-no}
JWP_CMD=${JWP_CMD:-"wine64 /home/dn/projects/bl3/jwp/apo-jwp.exe"}
WINDEBUG=fixme-all


echo $GAMEDATA_PATH > $OUTPUT_FILE


for MAPFILE in $(find $GAMEDATA_PATH -name '*.umap'); do
  echo "Checking map file $MAPFILE..."
  if strings $MAPFILE | grep -qE 'Spawner'; then
    echo "Found spawner data in $MAPFILE"
    SERIALIZED_FILE=$(echo $MAPFILE | sed 's/\.umap/.json/g')
    if [[ "$JWP_SERIALIZE" == 'yes' ]]; then
      $JWP_CMD serialize $(echo $MAPFILE | sed 's/\.umap//g') 2>&1 > /dev/null
      if [ ! -e "$SERIALIZED_FILE" ]; then
          echo "WARNING: jwp could not serialize map file $MAPFILE"
      else
          echo "$SERIALIZED_FILE" | tee -a $OUTPUT_FILE
      fi
    else
      if [[ -e $SERIALIZED_FILE ]]; then
          echo "$SERIALIZED_FILE" | tee -a $OUTPUT_FILE
      fi
    fi
  else
    echo "No spawner data in $MAPFILE"
  fi
done
