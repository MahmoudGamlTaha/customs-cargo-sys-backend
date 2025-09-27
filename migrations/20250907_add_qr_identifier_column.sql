-- Add qr_identifier column to requests table
ALTER TABLE requests ADD COLUMN qr_identifier VARCHAR(255);

-- Create index for qr_identifier for faster lookups
CREATE INDEX idx_requests_qr_identifier ON requests(qr_identifier);

-- Update existing requests to generate qr_identifier (this will be handled by the application)
-- The application will need to populate this field for existing records
