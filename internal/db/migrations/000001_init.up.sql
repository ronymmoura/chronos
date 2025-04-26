CREATE TABLE process (
  id int IDENTITY,
  name varchar(100) NOT NULL,
  description varchar(255) NULL,
  path varchar(255) NOT NULL,
  env text NULL,
  execute_every_secs int NOT NULL DEFAULT 0,
  created_at DateTime NOT NULL DEFAULT(getdate()),
  status varchar(50) NOT NULL DEFAULT('active'),
  running bit NOT NULL DEFAULT(0),
  PRIMARY KEY (id),
  CONSTRAINT status_check CHECK (status IN ('active', 'inactive', 'disabled'))
)

CREATE TABLE process_run (
  id int IDENTITY,
  process_id int NOT NULL,
  started_at DateTime NOT NULL DEFAULT(getdate()),
  ended_at DateTime NULL,
  success bit NOT NULL DEFAULT(0),
  PRIMARY KEY (id),
  FOREIGN KEY (process_id) REFERENCES process(id)
)

CREATE TABLE process_run_log (
  id int IDENTITY,
  process_run_id int NOT NULL,
  log_time DateTime NOT NULL DEFAULT(getdate()),
  message varchar(255) NOT NULL,
  type varchar(50) NOT NULL DEFAULT('info'),
  PRIMARY KEY (id),
  FOREIGN KEY (process_run_id) REFERENCES process_run(id),
  CONSTRAINT type_check CHECK (type IN ('info', 'error', 'warning'))
)
