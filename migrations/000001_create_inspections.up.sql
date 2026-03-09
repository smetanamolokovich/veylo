 CREATE TABLE inspections (                                                                                                                                                                                                                                  
      id              VARCHAR(26) PRIMARY KEY,                                                                                                                                                                                                                
      organization_id VARCHAR(26) NOT NULL,                                                                                                                                                                                                                 
      contract_number VARCHAR(255) NOT NULL,
      status          VARCHAR(50) NOT NULL DEFAULT 'new',
      created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
  );

  CREATE INDEX idx_inspections_organization_id ON inspections (organization_id);