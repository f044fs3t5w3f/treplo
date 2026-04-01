CREATE TABLE files (  
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    -- created_at DATE,
    file_id VARCHAR(255),
    chat_id BIGINT,
    message_id BIGINT,
    filepath VARCHAR(255) NULL,
    encoding varchar(20),
    salute_id VARCHAR(36) NULL,
    recognize_task_id VARCHAR(255) NULL,
    recognize_status VARCHAR(8) NULL,
    response_file_id VARCHAR(36) NULL,
    dialogue_content TEXT NULL,
    process_notification_sent boolean NOT NULL DEFAULT false,
);
-- CREATE UNIQUE INDEX file_id_1762799044404_index ON files USING btree (file_id);