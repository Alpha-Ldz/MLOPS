FROM bde2020/hive

COPY database.sql /scripts/database.sql

RUN beeline -u "jdbc:hive2://localhost:10000/default" -f /scripts/database.sql

EXPOSE 10000

CMD hive --service hiveserver2
