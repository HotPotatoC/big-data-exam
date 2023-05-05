CREATE KEYSPACE integrated_citizen
WITH replication = {
  'class': 'NetworkTopologyStrategy',
  'replication_factor': 3
};

CREATE TABLE IF NOT EXISTS integrated_citizen.citizen (
  citizen_id TEXT,
  message TEXT,
  transaction_from TEXT,
  transaction_to TEXT,
  total_transaction INT,
  timestamp TIMESTAMP,
  PRIMARY KEY (citizen_id, timestamp)
);