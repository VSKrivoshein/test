CREATE TABLE queue (
                       timestamp UInt64,
                       message String
) ENGINE = Kafka('host.docker.internal:9092', 'topic', 'group1', 'JSONEachRow');

SELECT * FROM queue LIMIT 5;