CREATE TABLE IF NOT EXISTS upload_files (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  biz_type VARCHAR(64) NOT NULL DEFAULT '',
  file_key VARCHAR(512) NOT NULL,
  url VARCHAR(1024) NOT NULL,
  mime_type VARCHAR(128) NOT NULL DEFAULT '',
  size BIGINT NOT NULL DEFAULT 0,
  created_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_upload_files_file_key (file_key),
  KEY idx_upload_files_biz_type (biz_type),
  KEY idx_upload_files_created_by (created_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
