-- Возврат типа id в таблице listening_history к SERIAL
ALTER TABLE listening_history ALTER COLUMN id DROP DEFAULT;
ALTER TABLE listening_history ALTER COLUMN id TYPE INTEGER USING (floor(random() * 1000000)::integer);
ALTER TABLE listening_history ALTER COLUMN id SET DEFAULT nextval('listening_history_id_seq'); 