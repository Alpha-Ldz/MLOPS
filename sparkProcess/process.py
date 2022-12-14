from pyspark.sql import SparkSession
from pyspark.sql.functions import *
import os

spark = SparkSession.builder.appName("process").getOrCreate()

# Lire les données de Kafka en utilisant le format de flux de données Apache Kafka
df = spark.readStream.format("kafka")\
    .option("kafka.bootstrap.servers", "localhost:9092")\
    .option("subscribe", "rapport-topic")\
    .load()

# Afficher les données brutes
df.selectExpr("CAST(key AS STRING)", "CAST(value AS STRING)") \
    .writeStream \
    .format("console") \
    .start() \
    .awaitTermination()

