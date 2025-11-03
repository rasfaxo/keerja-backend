-- Seed data for districts table (Sample data for testing)
-- Full data should be imported from BPS database or external source
-- This includes sample districts in major cities

-- Districts in Jakarta Pusat (city_id from code '3171')
INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Gambir', '3171010', '10110'
FROM cities c WHERE c.code = '3171'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Tanah Abang', '3171020', '10120'
FROM cities c WHERE c.code = '3171'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Menteng', '3171030', '10130'
FROM cities c WHERE c.code = '3171'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Senen', '3171040', '10140'
FROM cities c WHERE c.code = '3171'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Cempaka Putih', '3171050', '10150'
FROM cities c WHERE c.code = '3171'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Johar Baru', '3171060', '10160'
FROM cities c WHERE c.code = '3171'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Kemayoran', '3171070', '10170'
FROM cities c WHERE c.code = '3171'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Sawah Besar', '3171080', '10180'
FROM cities c WHERE c.code = '3171'
ON CONFLICT DO NOTHING;

-- Districts in Jakarta Selatan (city_id from code '3174')
INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Tebet', '3174010', '12810'
FROM cities c WHERE c.code = '3174'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Setiabudi', '3174020', '12910'
FROM cities c WHERE c.code = '3174'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Mampang Prapatan', '3174030', '12790'
FROM cities c WHERE c.code = '3174'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Pasar Minggu', '3174040', '12510'
FROM cities c WHERE c.code = '3174'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Kebayoran Lama', '3174050', '12210'
FROM cities c WHERE c.code = '3174'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Cilandak', '3174060', '12430'
FROM cities c WHERE c.code = '3174'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Kebayoran Baru', '3174070', '12110'
FROM cities c WHERE c.code = '3174'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Pancoran', '3174080', '12780'
FROM cities c WHERE c.code = '3174'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Jagakarsa', '3174090', '12620'
FROM cities c WHERE c.code = '3174'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Pesanggrahan', '3174100', '12270'
FROM cities c WHERE c.code = '3174'
ON CONFLICT DO NOTHING;

-- Districts in Bandung City (city_id from code '3273')
INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Bandung Wetan', '3273010', '40114'
FROM cities c WHERE c.code = '3273'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Sumur Bandung', '3273020', '40111'
FROM cities c WHERE c.code = '3273'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Cibeunying Kaler', '3273030', '40171'
FROM cities c WHERE c.code = '3273'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Cibeunying Kidul', '3273040', '40121'
FROM cities c WHERE c.code = '3273'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Coblong', '3273050', '40132'
FROM cities c WHERE c.code = '3273'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Sukasari', '3273060', '40152'
FROM cities c WHERE c.code = '3273'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Cidadap', '3273070', '40141'
FROM cities c WHERE c.code = '3273'
ON CONFLICT DO NOTHING;

-- Districts in Bandung Barat (city_id from code '3217')
INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Lembang', '3217010', '40391'
FROM cities c WHERE c.code = '3217'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Parongpong', '3217020', '40559'
FROM cities c WHERE c.code = '3217'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Cisarua', '3217030', '40551'
FROM cities c WHERE c.code = '3217'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Cikalong Wetan', '3217040', '40560'
FROM cities c WHERE c.code = '3217'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Cipeundeuy', '3217050', '40558'
FROM cities c WHERE c.code = '3217'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Ngamprah', '3217060', '40552'
FROM cities c WHERE c.code = '3217'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Cipatat', '3217070', '40553'
FROM cities c WHERE c.code = '3217'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Padalarang', '3217080', '40553'
FROM cities c WHERE c.code = '3217'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Batujajar', '3217090', '40561'
FROM cities c WHERE c.code = '3217'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Cihampelas', '3217100', '40562'
FROM cities c WHERE c.code = '3217'
ON CONFLICT DO NOTHING;

-- Districts in Tangerang Selatan (city_id from code '3674')
INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Serpong', '3674010', '15310'
FROM cities c WHERE c.code = '3674'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Serpong Utara', '3674020', '15310'
FROM cities c WHERE c.code = '3674'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Pondok Aren', '3674030', '15224'
FROM cities c WHERE c.code = '3674'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Ciputat', '3674040', '15411'
FROM cities c WHERE c.code = '3674'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Ciputat Timur', '3674050', '15412'
FROM cities c WHERE c.code = '3674'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Pamulang', '3674060', '15417'
FROM cities c WHERE c.code = '3674'
ON CONFLICT DO NOTHING;

INSERT INTO districts (city_id, name, code, postal_code) 
SELECT c.id, 'Setu', '3674070', '15314'
FROM cities c WHERE c.code = '3674'
ON CONFLICT DO NOTHING;

-- Update the sequence
SELECT setval('districts_id_seq', (SELECT MAX(id) FROM districts));

-- NOTE: This is sample data for testing purposes (only ~50 districts)
-- For production, import complete data (~7000+ districts anjay) from:
-- 1. BPS (Badan Pusat Statistik) official database
-- 2. https://github.com/cahyadsn/wilayah
-- 3. Other reliable Indonesian administrative data sources
--
-- The complete dataset should include all kecamatan across Indonesia
