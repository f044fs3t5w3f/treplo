CREATE TABLE files (  
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at DATE,
    file_id VARCHAR(255),
    chat_id BIGINT
);
CREATE UNIQUE INDEX file_id_1762799044404_index ON files USING btree (file_id);