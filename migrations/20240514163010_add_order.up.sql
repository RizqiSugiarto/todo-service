CREATE SEQUENCE tasks_order_position_seq;

ALTER TABLE tasks
ADD COLUMN "order_position" INT NOT NULL DEFAULT nextval('tasks_order_position_seq');

ALTER SEQUENCE tasks_order_position_seq OWNED BY tasks."order_position";
