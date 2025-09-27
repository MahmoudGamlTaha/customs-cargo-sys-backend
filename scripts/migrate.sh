#!/bin/bash

# Database migration script for Request Management System

set -e

# Configuration
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-request_management_system}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if PostgreSQL is available
check_postgres() {
    log_info "Checking PostgreSQL connection..."
    
    if ! command -v psql &> /dev/null; then
        log_error "psql command not found. Please install PostgreSQL client."
        exit 1
    fi
    
    export PGPASSWORD=$DB_PASSWORD
    
    if ! psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c '\q' 2>/dev/null; then
        log_error "Cannot connect to PostgreSQL server at $DB_HOST:$DB_PORT"
        exit 1
    fi
    
    log_info "PostgreSQL connection successful"
}

# Create database if it doesn't exist
create_database() {
    log_info "Creating database if it doesn't exist..."
    
    export PGPASSWORD=$DB_PASSWORD
    
    # Check if database exists
    if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -lqt | cut -d \| -f 1 | grep -qw $DB_NAME; then
        log_info "Database '$DB_NAME' already exists"
    else
        log_info "Creating database '$DB_NAME'..."
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "CREATE DATABASE $DB_NAME;"
        log_info "Database '$DB_NAME' created successfully"
    fi
}

# Run migrations
run_migrations() {
    log_info "Running database migrations..."
    
    export PGPASSWORD=$DB_PASSWORD
    
    # Get script directory
    SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
    PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
    
    # Run main schema
    if [ -f "$PROJECT_DIR/database_schema.sql" ]; then
        log_info "Applying main schema..."
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$PROJECT_DIR/database_schema.sql"
        log_info "Main schema applied successfully"
    else
        log_error "Main schema file not found: $PROJECT_DIR/database_schema.sql"
        exit 1
    fi
    
    # Run additional migrations if they exist
    MIGRATIONS_DIR="$PROJECT_DIR/migrations"
    if [ -d "$MIGRATIONS_DIR" ]; then
        log_info "Applying additional migrations..."
        for migration in "$MIGRATIONS_DIR"/*.sql; do
            if [ -f "$migration" ]; then
                log_info "Applying migration: $(basename "$migration")"
                psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$migration"
            fi
        done
        log_info "Additional migrations applied successfully"
    fi
}

# Verify installation
verify_installation() {
    log_info "Verifying database installation..."
    
    export PGPASSWORD=$DB_PASSWORD
    
    # Check if main tables exist
    tables=("users" "companies" "branches" "requests" "request_history" "user_sessions")
    
    for table in "${tables[@]}"; do
        if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "\dt $table" | grep -q "$table"; then
            log_info "Table '$table' exists"
        else
            log_error "Table '$table' not found"
            exit 1
        fi
    done
    
    # Check if sample data exists
    user_count=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM users;")
    log_info "Found $user_count users in database"
    
    log_info "Database verification completed successfully"
}

# Reset database (dangerous!)
reset_database() {
    log_warn "This will completely reset the database and all data will be lost!"
    read -p "Are you sure you want to continue? (yes/no): " confirm
    
    if [ "$confirm" = "yes" ]; then
        log_info "Resetting database..."
        
        export PGPASSWORD=$DB_PASSWORD
        
        # Drop and recreate database
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;"
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "CREATE DATABASE $DB_NAME;"
        
        log_info "Database reset completed"
        
        # Run migrations
        run_migrations
    else
        log_info "Database reset cancelled"
    fi
}

# Show help
show_help() {
    echo "Database Migration Script for Request Management System"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  migrate    Run database migrations (default)"
    echo "  reset      Reset database (WARNING: destroys all data)"
    echo "  verify     Verify database installation"
    echo "  help       Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  DB_HOST     Database host (default: localhost)"
    echo "  DB_PORT     Database port (default: 5432)"
    echo "  DB_USER     Database user (default: postgres)"
    echo "  DB_PASSWORD Database password (default: postgres)"
    echo "  DB_NAME     Database name (default: request_management_system)"
}

# Main execution
main() {
    case "${1:-migrate}" in
        "migrate")
            check_postgres
            create_database
            run_migrations
            verify_installation
            log_info "Migration completed successfully!"
            ;;
        "reset")
            check_postgres
            reset_database
            log_info "Reset completed successfully!"
            ;;
        "verify")
            check_postgres
            verify_installation
            log_info "Verification completed successfully!"
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            log_error "Unknown command: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"

