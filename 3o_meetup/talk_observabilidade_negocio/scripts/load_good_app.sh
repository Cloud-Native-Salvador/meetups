#!/bin/bash

TARGET_URL="http://localhost:8088/buy"
METHODS=("credit_card" "pix" "boleto")
REQUESTS_PER_BATCH=50
TOTAL_BATCHES=450

echo "Iniciando stress test contra: ${TARGET_URL}"
echo "→ ${REQUESTS_PER_BATCH} requisições simultâneas por batch."
echo "→ Total de ${TOTAL_BATCHES} batches."

for ((batch=1; batch<=TOTAL_BATCHES; batch++)); do
  echo "Batch $batch:"
  for ((i=1; i<=REQUESTS_PER_BATCH; i++)); do
    method=${METHODS[$RANDOM % ${#METHODS[@]}]}
    curl -s "${TARGET_URL}?method=${method}" > /dev/null &
  done
  wait
  sleep 0.1
done

echo "Carga finalizada para good_app."