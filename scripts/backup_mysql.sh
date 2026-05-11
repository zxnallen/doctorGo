#!/usr/bin/env bash
set -euo pipefail

BACKUP_DIR=${BACKUP_DIR:-/data/backups/doctor-go}
DATABASE=${DATABASE:-doctor_go}
MYSQL_USER=${MYSQL_USER:-root}

mkdir -p "$BACKUP_DIR"
mysqldump -u "$MYSQL_USER" -p "$DATABASE" | gzip > "$BACKUP_DIR/${DATABASE}_$(date +%Y%m%d_%H%M%S).sql.gz"
