#!/bin/bash

set -e

# Database credentials and connection info
DB_USER="admin"
DB_PASS="password"
DB_NAME="hecate"
DB_HOST="localhost"

# Output files
SQLITE_DB="./data/hecate.db"
mkdir -p ./data

# Tables to migrate
TABLES=("subreddits" "posts" "comments")

echo "Starting migration process..."

# Create SQLite tables
echo "Creating SQLite tables..."
sqlite3 "$SQLITE_DB" << 'EOF'
CREATE TABLE IF NOT EXISTS subreddits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    num_subscribers INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id TEXT UNIQUE NOT NULL,
    subreddit_name TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT,
    discussion_url TEXT,
    comment_count INTEGER,
    upvotes INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER,
    parent_comment_id INTEGER,
    content TEXT NOT NULL,
    comment_id TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (parent_comment_id) REFERENCES comments(id)
);
EOF
echo "SQLite tables created successfully."

# Migrate each table
for table in "${TABLES[@]}"; do
    echo "Starting migration for table: $table"
    
    # Create a temporary directory for this table
    TEMP_DIR=$(mktemp -d)
    DUMP_FILE="$TEMP_DIR/dump_${table}.sql"
    
    echo "Fetching column information for $table..."
    # Get all column names except 'metadata'
    ALL_COLUMNS=$(PGPASSWORD="$DB_PASS" psql \
        -h "$DB_HOST" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        -t \
        -A \
        -c "SELECT string_agg(column_name, ',') 
            FROM (
                SELECT column_name 
                FROM information_schema.columns 
                WHERE table_name = '$table' AND column_name != 'metadata'
                ORDER BY ordinal_position
            ) sub;")

    echo "All columns for $table: $ALL_COLUMNS"

    # Get column names excluding 'id' and 'metadata'
    COLUMNS_WITHOUT_ID=$(PGPASSWORD="$DB_PASS" psql \
        -h "$DB_HOST" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        -t \
        -A \
        -c "SELECT string_agg(column_name, ',') 
            FROM (
                SELECT column_name 
                FROM information_schema.columns 
                WHERE table_name = '$table' AND column_name NOT IN ('id', 'metadata')
                ORDER BY ordinal_position
            ) sub;")

    echo "Columns without id for $table: $COLUMNS_WITHOUT_ID"

    # Count the number of columns
    COLUMN_COUNT=$(echo "$COLUMNS_WITHOUT_ID" | awk -F',' '{print NF}')

    # Get the actual row count
    ACTUAL_ROW_COUNT=$(PGPASSWORD="$DB_PASS" psql \
        -h "$DB_HOST" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        -t \
        -A \
        -c "SELECT COUNT(*) FROM $table;")
    
    echo "Actual row count for $table: $ACTUAL_ROW_COUNT"

    echo "Dumping data from PostgreSQL for $table..."
    # Dump data
    PGPASSWORD="$DB_PASS" psql \
        -h "$DB_HOST" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        -t \
        -A \
        -F $'\t' \
        -c "SELECT $ALL_COLUMNS FROM $table ORDER BY id;" > "$DUMP_FILE"

    # Create SQLite import file
    SQLITE_FILE="$TEMP_DIR/sqlite_${table}.sql"
    echo "BEGIN TRANSACTION;" > "$SQLITE_FILE"

    echo "Processing data for $table..."
    # Process the dump and create INSERT statements
    TOTAL_ROWS=0
    SUCCESSFUL_INSERTS=0
    SKIPPED_ROWS=0

    # Set a timeout (e.g., 5 minutes)
    TIMEOUT=300
    SECONDS=0

    # Initialize variables for multi-line entries
    current_entry=""
    expected_columns=$((COLUMN_COUNT + 1))  # +1 for the id column

    while IFS= read -r line || [ -n "$line" ]; do
        # Append the current line to the current entry
        current_entry+="$line"
        
        # Check if we have a complete entry (all columns present)
        if [ "$(echo "$current_entry" | awk -F'\t' '{print NF}')" -eq $expected_columns ]; then
            ((TOTAL_ROWS++))

            echo "Processing row $TOTAL_ROWS: $current_entry"

            # Split the entry into an array
            IFS=$'\t' read -ra values <<< "$current_entry"

            # Remove the first element (id) from the array
            unset 'values[0]'
            
            # Build the VALUES part of the INSERT statement
            VALUES_STR=""
            for value in "${values[@]}"; do
                if [ -z "$value" ] || [ "$value" = "\\N" ]; then
                    VALUES_STR="${VALUES_STR}NULL,"
                else
                    # Escape single quotes and wrap in quotes
                    escaped_value=$(echo "$value" | sed "s/'/''/g")
                    VALUES_STR="${VALUES_STR}'${escaped_value}',"
                fi
            done
            VALUES_STR=${VALUES_STR%,}  # Remove trailing comma

            # Create the INSERT statement
            INSERT_STMT="INSERT OR IGNORE INTO $table ($COLUMNS_WITHOUT_ID) VALUES ($VALUES_STR);"
            echo "$INSERT_STMT" >> "$SQLITE_FILE"
            ((SUCCESSFUL_INSERTS++))

            # Reset the current entry
            current_entry=""

            # Print progress every 100 rows
            if [ $((TOTAL_ROWS % 100)) -eq 0 ]; then
                echo "Processed $TOTAL_ROWS rows for $table..."
            fi

            # Check for timeout
            if [ $SECONDS -ge $TIMEOUT ]; then
                echo "Error: Script timed out after $TIMEOUT seconds."
                break
            fi
        fi
    done < "$DUMP_FILE"

    echo "COMMIT;" >> "$SQLITE_FILE"

    echo "Data processing complete for $table."
    echo "Total rows: $TOTAL_ROWS"
    echo "Successful inserts: $SUCCESSFUL_INSERTS"
    echo "Skipped rows: $SKIPPED_ROWS"

    # Debug: Show first few lines of the SQLite file
    echo "Sample of generated SQL for $table:"
    head -n 5 "$SQLITE_FILE"

    echo "Importing data into SQLite for $table..."
    # Import into SQLite
    if ! sqlite3 "$SQLITE_DB" < "$SQLITE_FILE"; then
        echo "Error occurred while importing $table. Check the SQLite file at $SQLITE_FILE for issues."
        exit 1
    fi

    echo "Import successful for $table."

    # Verify the import
    SQLITE_ROW_COUNT=$(sqlite3 "$SQLITE_DB" "SELECT COUNT(*) FROM $table;")
    echo "Row count in SQLite for $table: $SQLITE_ROW_COUNT"
    if [ "$SQLITE_ROW_COUNT" -ne "$ACTUAL_ROW_COUNT" ]; then
        echo "Warning: Row count mismatch for $table. PostgreSQL: $ACTUAL_ROW_COUNT, SQLite: $SQLITE_ROW_COUNT"
    fi

    # Clean up
    rm -rf "$TEMP_DIR"
    echo "Temporary files cleaned up for $table."
done

echo "Migration complete! Database created at $SQLITE_DB"

# Verify the migration
echo "Verifying migration..."
for table in "${TABLES[@]}"; do
    echo "Count of records in $table:"
    sqlite3 "$SQLITE_DB" "SELECT COUNT(*) FROM $table;"
    echo "Sample data from $table:"
    sqlite3 "$SQLITE_DB" "SELECT * FROM $table LIMIT 1;"
done

echo "Migration and verification complete."
