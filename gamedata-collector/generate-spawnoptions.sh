#!/bin/bash


GAMEDATA_PATH=${GAMEDATA_PATH}
OUTPUT_FILE=${OUTPUT_FILE}
JWP_CMD=${JWP_CMD:-"wine64 /opt/bl3/jwp/apo-jwp.exe"}
WINDEBUG=fixme-all


echo $GAMEDATA_PATH > $OUTPUT_FILE


for SOFILE in $(find $GAMEDATA_PATH -name '*.uasset' | grep 'Enemies.*Spawning'); do
  echo "Serializing $OFILE"
  $JWP_CMD serialize $(echo $SOFILE | sed 's/\.uasset//g') 2>&1 > /dev/null
  SERIALIZED_FILE=$(echo $SOFILE | sed 's/\.uasset/.json/g')
  if [ ! -e "$SERIALIZED_FILE" ]; then
      echo "WARNING: jwp could not serialize asset file $SOFILE"
  else
      echo "$SERIALIZED_FILE" | tee -a $OUTPUT_FILE
  fi
done
