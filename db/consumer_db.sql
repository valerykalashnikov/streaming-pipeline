CREATE TABLE statistics (
    consumer_id integer NOT NULL,
    consumed_data bigint NOT NULL,
    UNIQUE (consumer_id)
);