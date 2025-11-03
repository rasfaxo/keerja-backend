-- Seed data for industries table
-- This includes common industries in Indonesia and globally

INSERT INTO industries (name, slug, description, display_order) VALUES
('Technology', 'technology', 'Information technology, software development, and IT services', 1),
('Healthcare', 'healthcare', 'Healthcare services, medical facilities, and pharmaceuticals', 2),
('Education', 'education', 'Educational institutions, training centers, and e-learning', 3),
('Finance', 'finance', 'Banking, insurance, investment, and financial services', 4),
('Retail', 'retail', 'Retail stores, e-commerce, and consumer goods', 5),
('Manufacturing', 'manufacturing', 'Manufacturing and industrial production', 6),
('Construction', 'construction', 'Construction, real estate development, and infrastructure', 7),
('Transportation', 'transportation', 'Logistics, shipping, and transportation services', 8),
('Hospitality', 'hospitality', 'Hotels, restaurants, tourism, and food & beverage', 9),
('Agriculture', 'agriculture', 'Agriculture, farming, and agribusiness', 10),
('Telecommunications', 'telecommunications', 'Telecommunications and internet service providers', 11),
('Media & Entertainment', 'media-entertainment', 'Media production, entertainment, and creative industries', 12),
('Energy', 'energy', 'Energy production, oil & gas, and renewable energy', 13),
('Automotive', 'automotive', 'Automotive manufacturing and services', 14),
('Fashion', 'fashion', 'Fashion, apparel, and textile industry', 15),
('Consulting', 'consulting', 'Business consulting and professional services', 16),
('Marketing & Advertising', 'marketing-advertising', 'Marketing, advertising, and public relations', 17),
('Legal Services', 'legal-services', 'Legal firms and legal services', 18),
('Real Estate', 'real-estate', 'Real estate agencies and property management', 19),
('Non-Profit', 'non-profit', 'Non-profit organizations and NGOs', 20),
('Government', 'government', 'Government agencies and public sector', 21),
('Arts & Crafts', 'arts-crafts', 'Arts, crafts, and creative industries', 22),
('Sports & Recreation', 'sports-recreation', 'Sports, fitness, and recreational services', 23),
('Beauty & Wellness', 'beauty-wellness', 'Beauty salons, spas, and wellness centers', 24),
('Other', 'other', 'Other industries not listed above', 99)
ON CONFLICT (name) DO NOTHING;

-- Update the sequence to ensure next insert starts from correct value
SELECT setval('industries_id_seq', (SELECT MAX(id) FROM industries));
