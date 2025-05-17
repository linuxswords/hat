-- Insert HandycapSet
INSERT INTO handycap_sets (name, created_at, updated_at)
VALUES ('3D 2021-2025 3-Arrow Round', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- 2. Save the inserted ID
-- In SQLite, we can use 1 to get the ID of the inserted row

-- Insert HandycapEntries
-- 3. Insert HandycapEntry rows
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at)
SELECT bc.id, 1.000, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AMHB';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.963, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'CMFS-R';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.127, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'CMBHR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.864, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'CFBHR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.188, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'JFBU';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.722, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AFHB';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.947, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AMLB';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.712, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AFLB';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.002, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VMLB';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.829, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VFLB';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.838, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'SMLB';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.505, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AMFU';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.531, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AFFU';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.509, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VMFU';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.432, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'SMFU';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.230, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VMFSR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.060, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AMBBR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.093, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AFBBR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.132, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VFBBR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.810, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'YAMTR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.632, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'YAFTR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.145, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AMTR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.862, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AFTR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.066, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VMTR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.878, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VFTR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.021, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'YAFBHR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.179, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AMBHR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.819, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AFBHR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.115, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VMBHR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.832, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VFBHR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.021, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'SMBHR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.922, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'SFBHR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.480, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AMBU';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.218, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AFBU';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.474, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VMBU';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.237, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VFBU';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.062, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'JFFU';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.765, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'JMFSR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.243, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'CMFU';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.922, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'AMFSR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.030, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'YAMBHR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.516, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'JMFU';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 1.097, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'VMBBR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.873, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'JMTR';
INSERT INTO handycap_entries (bow_class_id, value, handycap_set_id, created_at, updated_at) SELECT bc.id, 0.533, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP FROM bow_classes bc WHERE bc.code = 'JMLB';
