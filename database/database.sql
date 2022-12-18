CREATE DATABASE mydatabase;

USE mydatabase;

CREATE TABLE mytable (
  id INT,
  name STRING,
  age INT
)
STORED AS ORC;
