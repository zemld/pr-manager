#!/bin/bash

# Скрипт для подготовки тестовых данных для нагрузочного тестирования
# Создает 20 команд с 200 пользователями (по 10 пользователей в команде)

BASE_URL="http://localhost:8080"

echo "=== Подготовка тестовых данных ==="

# Создаем 20 команд
for team_num in {1..20}; do
  team_name="team_${team_num}"
  
  # Формируем JSON для команды с 10 пользователями
  members_json="["
  for user_num in {1..10}; do
    user_id="u_${team_num}_${user_num}"
    username="User_${team_num}_${user_num}"
    is_active="true"
    if [ $user_num -le 8 ]; then
      is_active="true"  # 8 активных пользователей
    else
      is_active="false"  # 2 неактивных пользователя
    fi
    
    if [ $user_num -gt 1 ]; then
      members_json+=","
    fi
    members_json+="{\"user_id\":\"${user_id}\",\"username\":\"${username}\",\"is_active\":${is_active}}"
  done
  members_json+="]"
  
  # Создаем команду
  response=$(curl -s -X POST "${BASE_URL}/team/add" \
    -H "Content-Type: application/json" \
    -d "{\"team_name\":\"${team_name}\",\"members\":${members_json}}")
  
  if echo "$response" | grep -q "team_name"; then
    echo "✓ Команда ${team_name} создана"
  else
    echo "✗ Ошибка создания команды ${team_name}: $response"
  fi
done

echo ""
echo "=== Создание Pull Request'ов ==="

# Создаем PR для каждого активного пользователя (по 2-3 PR на пользователя)
pr_count=0
for team_num in {1..20}; do
  for user_num in {1..8}; do  # Только активные пользователи
    user_id="u_${team_num}_${user_num}"
    
    # Создаем 2-3 PR для каждого пользователя
    for pr_num in {1..3}; do
      pr_id="pr_${team_num}_${user_num}_${pr_num}"
      pr_name="PR ${team_num}-${user_num}-${pr_num}"
      
      response=$(curl -s -X POST "${BASE_URL}/pullRequest/create" \
        -H "Content-Type: application/json" \
        -d "{\"pull_request_id\":\"${pr_id}\",\"pull_request_name\":\"${pr_name}\",\"author_id\":\"${user_id}\"}")
      
      if echo "$response" | grep -q "pull_request_id"; then
        pr_count=$((pr_count + 1))
        if [ $((pr_count % 50)) -eq 0 ]; then
          echo "✓ Создано ${pr_count} PR"
        fi
      fi
    done
  done
done

echo "✓ Всего создано ${pr_count} PR"

# Мержим часть PR (примерно 30%)
echo ""
echo "=== Мержим часть PR ==="
merged_count=0
for team_num in {1..20}; do
  for user_num in {1..8}; do
    for pr_num in {1..3}; do
      # Мержим каждый третий PR
      if [ $((pr_num % 3)) -eq 0 ]; then
        pr_id="pr_${team_num}_${user_num}_${pr_num}"
        
        response=$(curl -s -X POST "${BASE_URL}/pullRequest/merge" \
          -H "Content-Type: application/json" \
          -d "{\"pull_request_id\":\"${pr_id}\"}")
        
        if echo "$response" | grep -q "MERGED"; then
          merged_count=$((merged_count + 1))
        fi
      fi
    done
  done
done

echo "✓ Смержено ${merged_count} PR"

echo ""
echo "=== Подготовка данных завершена ==="
echo "Команды: 20"
echo "Пользователи: 200 (160 активных, 40 неактивных)"
echo "PR: ~${pr_count} (из них ~${merged_count} смержено)"


