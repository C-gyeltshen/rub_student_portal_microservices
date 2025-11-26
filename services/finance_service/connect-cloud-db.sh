#!/bin/bash
# Script to connect to Render PostgreSQL database

DB_URL="postgresql://finance_service_14r5_user:6QJcvvToLvJ5EUGBGlrYWq15JIV4PrOC@dpg-d4jdlj7gi27c739kmt3g-a.singapore-postgres.render.com/finance_service_14r5?sslmode=require"

echo "Connecting to Render PostgreSQL database..."
echo "Database: finance_service_14r5"
echo "Host: dpg-d4jdlj7gi27c739kmt3g-a.singapore-postgres.render.com"
echo ""

PGPASSWORD="6QJcvvToLvJ5EUGBGlrYWq15JIV4PrOC" psql -h dpg-d4jdlj7gi27c739kmt3g-a.singapore-postgres.render.com -U finance_service_14r5_user -d finance_service_14r5 -p 5432
