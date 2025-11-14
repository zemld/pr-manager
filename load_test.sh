#!/bin/bash

# Скрипт для нагрузочного тестирования PR Manager API

BASE_URL="http://localhost:8080"
RESULTS_DIR="load_test_results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

mkdir -p "${RESULTS_DIR}"

echo "=== Нагрузочное тестирование PR Manager API ==="
echo "Время начала: $(date)"
echo ""

# Функция для запуска теста и сохранения результатов
run_test() {
  local test_name=$1
  local method=$2
  local endpoint=$3
  local data=$4
  local duration=$5
  local rate=$6
  
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "Тест: ${test_name}"
  echo "Эндпоинт: ${method} ${endpoint}"
  echo "Длительность: ${duration}с, Rate: ${rate} req/s"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  
  local output_file="${RESULTS_DIR}/${test_name}_${TIMESTAMP}.txt"
  
  if [ -n "$data" ]; then
    echo "$data" > /tmp/load_test_data.json
    hey -n 0 -c ${rate} -z ${duration}s -m ${method} \
      -H "Content-Type: application/json" \
      -D /tmp/load_test_data.json \
      "${BASE_URL}${endpoint}" > "${output_file}" 2>&1
  else
    hey -n 0 -c ${rate} -z ${duration}s -m ${method} \
      "${BASE_URL}${endpoint}" > "${output_file}" 2>&1
  fi
  
  echo ""
  echo "Результаты:"
  tail -20 "${output_file}"
  echo ""
  echo "Полные результаты сохранены в: ${output_file}"
  echo ""
}

# Требования к производительности
# RPS: 5 запросов в секунду
# SLI времени ответа: 300 мс
# SLI успешности: 99.9%

# Тест 1: GET /stats/get (самый тяжелый эндпоинт)
echo ">>> ТЕСТ 1: GET /stats/get (тяжелый запрос)"
run_test "stats_get" "GET" "/stats/get" "" "30" "5"

# Тест 2: GET /team/get
echo ">>> ТЕСТ 2: GET /team/get"
run_test "team_get" "GET" "/team/get?name=team_1" "" "30" "5"

# Тест 3: POST /pullRequest/create
echo ">>> ТЕСТ 3: POST /pullRequest/create"
PR_CREATE_DATA='{"pull_request_id":"pr_load_test_1","pull_request_name":"Load Test PR","author_id":"u_1_1"}'
run_test "pr_create" "POST" "/pullRequest/create" "${PR_CREATE_DATA}" "30" "5"

# Тест 4: POST /pullRequest/merge
echo ">>> ТЕСТ 4: POST /pullRequest/merge"
PR_MERGE_DATA='{"pull_request_id":"pr_1_1_1"}'
run_test "pr_merge" "POST" "/pullRequest/merge" "${PR_MERGE_DATA}" "30" "5"

# Тест 5: GET /users/getReview
echo ">>> ТЕСТ 5: GET /users/getReview"
run_test "user_review" "GET" "/users/getReview?user_id=u_1_1" "" "30" "5"

# Тест 6: Комбинированная нагрузка (все эндпоинты одновременно)
echo ">>> ТЕСТ 6: Комбинированная нагрузка (все эндпоинты)"
echo "Запускаем параллельно несколько эндпоинтов..."

# Запускаем в фоне несколько тестов одновременно
hey -n 0 -c 2 -z 30s -m GET "${BASE_URL}/stats/get" > "${RESULTS_DIR}/combined_stats_${TIMESTAMP}.txt" 2>&1 &
hey -n 0 -c 1 -z 30s -m GET "${BASE_URL}/team/get?name=team_1" > "${RESULTS_DIR}/combined_team_${TIMESTAMP}.txt" 2>&1 &
echo '{"pull_request_id":"pr_load_test_2","pull_request_name":"Load Test PR 2","author_id":"u_1_2"}' > /tmp/load_test_pr.json
hey -n 0 -c 1 -z 30s -m POST -H "Content-Type: application/json" -D /tmp/load_test_pr.json "${BASE_URL}/pullRequest/create" > "${RESULTS_DIR}/combined_pr_create_${TIMESTAMP}.txt" 2>&1 &
hey -n 0 -c 1 -z 30s -m GET "${BASE_URL}/users/getReview?user_id=u_1_2" > "${RESULTS_DIR}/combined_user_review_${TIMESTAMP}.txt" 2>&1 &

wait

echo "Комбинированная нагрузка завершена"
echo ""

# Тест 7: Нагрузка выше требований (стресс-тест)
echo ">>> ТЕСТ 7: Стресс-тест (10 RPS вместо 5)"
run_test "stress_test" "GET" "/stats/get" "" "30" "10"

# Генерируем сводный отчет
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "=== СВОДНЫЙ ОТЧЕТ ==="
echo "Время окончания: $(date)"
echo "Результаты сохранены в директории: ${RESULTS_DIR}/"
echo ""

# Анализируем результаты
for result_file in "${RESULTS_DIR}"/*_${TIMESTAMP}.txt; do
  if [ -f "$result_file" ]; then
    echo "Файл: $(basename $result_file)"
    echo "---"
    grep -E "(Total:|Requests/sec|Average:|50%|95%|99%|Slowest|Fastest|Total data)" "$result_file" | head -10
    echo ""
  fi
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "=== ТРЕБОВАНИЯ К ПРОИЗВОДИТЕЛЬНОСТИ ==="
echo "✓ RPS: ≥ 5 запросов в секунду"
echo "✓ SLI времени ответа: ≤ 300 мс"
echo "✓ SLI успешности: ≥ 99.9%"
echo ""
echo "Проверьте результаты выше на соответствие требованиям."


