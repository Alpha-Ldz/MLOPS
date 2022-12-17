from pyspark.sql import SparkSession
from pyspark.sql.functions import *
import os

spark = SparkSession.builder.appName("process").getOrCreate()

df = spark.readStream.format("kafka")\
    .option("kafka.bootstrap.servers", "localhost:9092")\
    .option("subscribe", "rapport-topic2")\
    .load()

schema = (
        StructType()
        .add("Id", StringType())
        .add("User", StringType())
        .add("Vtype", StringType())
        .add("Longitude", StringType())
        .add("Latitude", StringType())
        .add("Status", StringType())
)
df = df.selectExpr("CAST(value AS STRING)") \
    .select(from_json(col("value"), schema).alias("data")).select("data.*") \
    .writeStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers", "localhost:9092") \
    .option("checkpointLocation", "/tmp/checkpoint") \
    .option("topic", "website-topic") \
    .start() \
    .awaitTermination()

