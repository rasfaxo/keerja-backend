-- Seed data for cities table (Sample data for major cities)
-- Full data should be imported from BPS database or external source
-- This includes major cities in Java for testing purposes

-- DKI Jakarta (province_id will be determined by the code '31')
INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Jakarta Pusat', 'Kota', '3171'
FROM provinces p WHERE p.code = '31'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Jakarta Utara', 'Kota', '3172'
FROM provinces p WHERE p.code = '31'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Jakarta Barat', 'Kota', '3173'
FROM provinces p WHERE p.code = '31'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Jakarta Selatan', 'Kota', '3174'
FROM provinces p WHERE p.code = '31'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Jakarta Timur', 'Kota', '3175'
FROM provinces p WHERE p.code = '31'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Kepulauan Seribu', 'Kabupaten', '3101'
FROM provinces p WHERE p.code = '31'
ON CONFLICT DO NOTHING;

-- Jawa Barat (province_id from code '32')
INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Bogor', 'Kabupaten', '3201'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Sukabumi', 'Kabupaten', '3202'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Cianjur', 'Kabupaten', '3203'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Bandung', 'Kabupaten', '3204'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Garut', 'Kabupaten', '3205'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Tasikmalaya', 'Kabupaten', '3206'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Ciamis', 'Kabupaten', '3207'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Kuningan', 'Kabupaten', '3208'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Cirebon', 'Kabupaten', '3209'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Majalengka', 'Kabupaten', '3210'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Sumedang', 'Kabupaten', '3211'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Indramayu', 'Kabupaten', '3212'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Subang', 'Kabupaten', '3213'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Purwakarta', 'Kabupaten', '3214'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Karawang', 'Kabupaten', '3215'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Bekasi', 'Kabupaten', '3216'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Bandung Barat', 'Kabupaten', '3217'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Pangandaran', 'Kabupaten', '3218'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

-- Kota in Jawa Barat
INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Bogor', 'Kota', '3271'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Sukabumi', 'Kota', '3272'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Bandung', 'Kota', '3273'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Cirebon', 'Kota', '3274'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Bekasi', 'Kota', '3275'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Depok', 'Kota', '3276'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Cimahi', 'Kota', '3277'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Tasikmalaya', 'Kota', '3278'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Banjar', 'Kota', '3279'
FROM provinces p WHERE p.code = '32'
ON CONFLICT DO NOTHING;

-- Banten (province_id from code '36')
INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Pandeglang', 'Kabupaten', '3601'
FROM provinces p WHERE p.code = '36'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Lebak', 'Kabupaten', '3602'
FROM provinces p WHERE p.code = '36'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Tangerang', 'Kabupaten', '3603'
FROM provinces p WHERE p.code = '36'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Serang', 'Kabupaten', '3604'
FROM provinces p WHERE p.code = '36'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Tangerang', 'Kota', '3671'
FROM provinces p WHERE p.code = '36'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Cilegon', 'Kota', '3672'
FROM provinces p WHERE p.code = '36'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Serang', 'Kota', '3673'
FROM provinces p WHERE p.code = '36'
ON CONFLICT DO NOTHING;

INSERT INTO cities (province_id, name, type, code) 
SELECT p.id, 'Tangerang Selatan', 'Kota', '3674'
FROM provinces p WHERE p.code = '36'
ON CONFLICT DO NOTHING;

-- Update the sequence
SELECT setval('cities_id_seq', (SELECT MAX(id) FROM cities));

-- NOTE: This is sample data for testing purposes
-- For production, import complete data from:
-- 1. BPS (Badan Pusat Statistik) official database
-- 2. https://github.com/cahyadsn/wilayah
-- 3. Other reliable Indonesian administrative data sources
