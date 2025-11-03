-- Seed data for provinces table
-- This includes all 34 provinces of Indonesia based on BPS data

INSERT INTO provinces (name, code) VALUES
('Aceh', '11'),
('Sumatera Utara', '12'),
('Sumatera Barat', '13'),
('Riau', '14'),
('Jambi', '15'),
('Sumatera Selatan', '16'),
('Bengkulu', '17'),
('Lampung', '18'),
('Kepulauan Bangka Belitung', '19'),
('Kepulauan Riau', '21'),
('DKI Jakarta', '31'),
('Jawa Barat', '32'),
('Jawa Tengah', '33'),
('DI Yogyakarta', '34'),
('Jawa Timur', '35'),
('Banten', '36'),
('Bali', '51'),
('Nusa Tenggara Barat', '52'),
('Nusa Tenggara Timur', '53'),
('Kalimantan Barat', '61'),
('Kalimantan Tengah', '62'),
('Kalimantan Selatan', '63'),
('Kalimantan Timur', '64'),
('Kalimantan Utara', '65'),
('Sulawesi Utara', '71'),
('Sulawesi Tengah', '72'),
('Sulawesi Selatan', '73'),
('Sulawesi Tenggara', '74'),
('Gorontalo', '75'),
('Sulawesi Barat', '76'),
('Maluku', '81'),
('Maluku Utara', '82'),
('Papua Barat', '91'),
('Papua', '94')
ON CONFLICT (code) DO NOTHING;

-- Update the sequence to ensure next insert starts from correct value
SELECT setval('provinces_id_seq', (SELECT MAX(id) FROM provinces));
