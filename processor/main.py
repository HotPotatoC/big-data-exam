import json
import os

from kafka import KafkaConsumer, TopicPartition
from pyspark.sql import SparkSession
from pyspark.sql.functions import col, from_json
from pyspark.sql.types import (
    IntegerType,
    StringType,
    StructField,
    StructType,
    TimestampType,
)

# consumer = KafkaConsumer(bootstrap_servers="localhost:9092")

# consumer.subscribe(["messages", "transactions"])

# for message in consumer:
#     print(json.loads(message.value.decode()))

os.environ[
    "PYSPARK_SUBMIT_ARGS"
] = "--packages org.apache.spark:spark-streaming-kafka-0-10_2.12:3.2.0,org.apache.spark:spark-sql-kafka-0-10_2.12:3.2.0,com.datastax.spark:spark-cassandra-connector_2.12:3.2.0 pyspark-shell"


spark = SparkSession.builder.appName("citizen_data_integration").getOrCreate()

df = (
    spark.readStream.format("kafka")
    .option("kafka.bootstrap.servers", "localhost:9092")
    .option("subscribe", "messages,transactions")
    .load()
)

messages_schema = StructType(
    [
        StructField("citizen_id", StringType()),
        StructField("message", StringType()),
        StructField("timestamp", TimestampType()),
    ]
)


transactions_schema = StructType(
    [
        StructField("citizen_id", StringType()),
        StructField("from", StringType()),
        StructField("to", StringType()),
        StructField("total", IntegerType()),
        StructField("timestamp", TimestampType()),
    ]
)

# Parse the JSON data for the messages topic
messages_df = (
    df.filter(col("topic") == "messages")
    .select(from_json(col("value").cast("string"), messages_schema).alias("data"))
    .select("data.*")
)

# Parse the JSON data for the transactions topic
transactions_df = (
    df.filter(col("topic") == "transactions")
    .select(from_json(col("value").cast("string"), transactions_schema).alias("data"))
    .select(
        col("data.citizen_id"),
        col("data.from"),
        col("data.to"),
        col("data.total"),
        col("data.timestamp").alias("transaction_timestamp"),
    )
)

integrated_df = messages_df.join(transactions_df, "citizen_id", "inner")
integrated_df = (
    integrated_df.withColumnRenamed("from", "transaction_from")
    .withColumnRenamed("to", "transaction_to")
    .withColumnRenamed("total", "total_transaction")
    .drop("transaction_timestamp")
)

integrated_df.writeStream.foreachBatch(
    lambda batch_df, batch_id: batch_df.write.format("org.apache.spark.sql.cassandra")
    .option("keyspace", "integrated_citizen")
    .option("table", "citizen")
    .option("spark.cassandra.connection.host", "127.0.0.1")
    .option("spark.cassandra.connection.port", "9042")
    .mode("append")
    .save()
).start().awaitTermination()

query = integrated_df.writeStream.outputMode("append").format("console").start()
query.awaitTermination()
