#!/bin/bash

TARGET_URL="http://localhost:8087/buy"
REQUESTS_PER_BATCH=30
TOTAL_BATCHES=300

echo "Iniciando carga de stress no serviço: bad_app (${TARGET_URL})"
echo "→ ${REQUESTS_PER_BATCH} requisições simultâneas por batch."
echo "→ Total de ${TOTAL_BATCHES} batches."

for ((batch=1; batch<=TOTAL_BATCHES; batch++)); do
  echo "Batch $batch:"
  for ((i=1; i<=REQUESTS_PER_BATCH; i++)); do
    curl -s "${TARGET_URL}" > /dev/null &
  done
  wait
  sleep 0.1
done

echo "Carga finalizada para bad_app."