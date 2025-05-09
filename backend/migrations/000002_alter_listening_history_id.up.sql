-- Изменение типа id в таблице listening_history с SERIAL на UUID
ALTER TABLE listening_history ALTER COLUMN id DROP DEFAULT;
ALTER TABLE listening_history ALTER COLUMN id TYPE UUID USING (uuid_generate_v4());
ALTER TABLE listening_history ALTER COLUMN id SET DEFAULT uuid_generate_v4(); 