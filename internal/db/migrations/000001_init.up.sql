CREATE TABLE process (
  id int IDENTITY,
  name varchar(100) NOT NULL,
  description varchar(255) NULL,
  path varchar(255) NOT NULL,
  execute_every_secs int NOT NULL DEFAULT 0,
  created_at DateTime NOT NULL DEFAULT(getdate()),
  last_run DateTime NULL 
)
