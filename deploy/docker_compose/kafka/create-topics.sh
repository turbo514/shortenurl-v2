#!/bin/bash

BOOTSTRAP_SERVER="localhost:9092"

KAFKA_BIN_DIR="/opt/kafka/bin"

TIMEOUT=60            # 等待 Kafka 就绪的超时时间（秒）
SLEEP_INTERVAL=3      # 每次尝试间隔（秒）

# --- 定义要创建的 topic 列表 ---
topics=(
  "link.click:3:1"
  "link.create:1:1"
  "link.click.dlq:1:1" # 例如，一个死信队列
)

# --- 脚本执行开始 ---

# 1. 等待 Kafka Broker 就绪
echo "Waiting for Kafka broker to be ready..."

for i in $(seq 1 $TIMEOUT); do
  # 尝试执行一个简单的、无副作用的命令来检查连接
  "$KAFKA_BIN_DIR/kafka-topics.sh" --bootstrap-server "$BOOTSTRAP_SERVER" --list 2>/dev/null

  if [ $? -eq 0 ]; then
    echo "Kafka broker is ready after $((i * SLEEP_INTERVAL)) seconds."
    break
  fi

  if [ $i -eq $TIMEOUT ]; then
      echo "Error: Kafka broker did not start within the timeout."
      exit 1
  fi

  sleep $SLEEP_INTERVAL
done

# 2. 循环创建 Topic
echo "Starting topic creation..."
for t in "${topics[@]}"; do
  # 使用参数解构: IFS=":" 用于分隔符
  IFS=":" read -r name partitions rf <<< "$t"

  echo "  -> Creating topic: $name (P:$partitions, RF:$rf)"

  "$KAFKA_BIN_DIR/kafka-topics.sh" \
    --create \
    --if-not-exists \
    --bootstrap-server "$BOOTSTRAP_SERVER" \
    --topic "$name" \
    --partitions "$partitions" \
    --replication-factor "$rf"
done

echo "All topics created successfully."