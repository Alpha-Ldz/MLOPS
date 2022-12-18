from pyspark.sql import SparkSession
from pyspark.sql.functions import from_json, col
from pyspark.sql.types import StructType, StructField, LongType, IntegerType, DoubleType, StringType

# Create a Spark session
spark = SparkSession.builder.appName("KafkaStreamingExample").getOrCreate()

# Define the schema of the data
schema = StructType([
    StructField("Id", LongType(), True),
    StructField("User", LongType(), True),
    StructField("VType", IntegerType(), True),
    StructField("Longitude", DoubleType(), True),
    StructField("Latitude", DoubleType(), True),
    StructField("Status", StringType(), True)
])

# Create a DataFrame from the Kafka topic
df = spark \
  .readStream \
  .format("kafka") \
  .option("kafka.bootstrap.servers", "localhost:9092") \
  .option("subscribe", "rapport-topic2") \
  .load() \
  .selectExpr("CAST(Value AS STRING)") \
  .select(from_json(col("Value"), schema).alias("Data")) \

# Write in console
query = df \
    .writeStream \
    .outputMode("append") \
    .format("console") \
    .start()

query.awaitTermination()
