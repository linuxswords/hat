-- Insert HandicapSet
INSERT INTO handicap_sets (name, created_at, updated_at)
VALUES ('3D 2021-2025 3-Arrow Round (KI generated)', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- 2. Save the inserted ID
-- In SQLite, we can use 1 to get the ID of the inserted row

-- Insert HandicapEntries
-- 3. Insert HandicapEntry rows
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at)
SELECT bc.id, 1.000, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AMHB';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.038, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'CMFS-R';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.887, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'CMBHR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.158, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'CFBHR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.842, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'JFBU';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.384, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AFHB';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.056, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AMLB';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.404, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AFLB';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.998, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VMLB';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.207, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VFLB';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.193, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'SMLB';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.664, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AMFU';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.653, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AFFU';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.663, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VMFU';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.698, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'SMFU';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.813, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VMFSR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.943, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AMBBR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.915, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AFBBR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.883, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VFBBR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.234, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'YAMTR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.582, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'YAFTR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.873, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AMTR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.160, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AFTR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.938, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VMTR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.140, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VFTR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.979, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'YAFBHR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.848, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AMBHR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.221, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AFBHR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.897, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VMBHR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.202, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VFBHR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.979, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'SMBHR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.084, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'SFBHR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.676, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AMBU';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.821, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AFBU';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.679, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VMBU';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.808, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VFBU';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.941, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'JFFU';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.308, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'JMFSR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.804, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'CMFU';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.084, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AMFSR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.972, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'YAMBHR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.660, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'JMFU';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 0.911, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VMBBR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.146, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'JMTR';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.876, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'JMLB';
INSERT INTO handicap_entries (bow_class_id, value, handicap_set_id, created_at, updated_at) SELECT bc.id, 1.276, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'JFBHR';
