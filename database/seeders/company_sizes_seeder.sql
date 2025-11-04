-- Seed data for company_sizes table
-- This includes standard company size categories used in Indonesia

INSERT INTO company_sizes (label, min_employees, max_employees, display_order) VALUES
('1 - 10 karyawan', 1, 10, 1),
('11 - 50 karyawan', 11, 50, 2),
('51 - 200 karyawan', 51, 200, 3),
('201 - 500 karyawan', 201, 500, 4),
('501 - 1000 karyawan', 501, 1000, 5),
('1000+ karyawan', 1001, NULL, 6)
ON CONFLICT DO NOTHING;

-- Update the sequence to ensure next insert starts from correct value
SELECT setval('company_sizes_id_seq', (SELECT MAX(id) FROM company_sizes));
