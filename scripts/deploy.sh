#!/bin/bash

# Deployment script for Request Management System

set -e

# Configuration
ENVIRONMENT=${ENVIRONMENT:-production}
PROJECT_NAME="Chumber-Workflow-System"
DOCKER_IMAGE="$PROJECT_NAME:latest"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    log_step "Checking prerequisites..."
    
    # Check if Docker is installed
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    # Check if Docker Compose is installed
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    
    # Check if .env file exists for production
    if [ "$ENVIRONMENT" = "production" ] && [ ! -f ".env" ]; then
        log_warn ".env file not found. Creating from template..."
        cp .env.example .env
        log_warn "Please edit .env file with production values before continuing."
        read -p "Press Enter to continue after editing .env file..."
    fi
    
    log_info "Prerequisites check completed"
}

# Build application
build_application() {
    log_step "Building application..."
    
    # Build Docker image
    log_info "Building Docker image..."
    docker build -t $DOCKER_IMAGE .
    
    log_info "Application build completed"
}

# Deploy application
deploy_application() {
    log_step "Deploying application..."
    
    if [ "$ENVIRONMENT" = "production" ]; then
        log_info "Deploying to production environment..."
        
        # Stop existing containers
        log_info "Stopping existing containers..."
        docker-compose -f docker-compose.yml -f docker-compose.prod.yml down || true
        
        # Start production containers
        log_info "Starting production containers..."
        docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
        
    else
        log_info "Deploying to development environment..."
        
        # Stop existing containers
        log_info "Stopping existing containers..."
        docker-compose down || true
        
        # Start development containers
        log_info "Starting development containers..."
        docker-compose up -d
    fi
    
    log_info "Application deployment completed"
}

# Wait for services to be ready
wait_for_services() {
    log_step "Waiting for services to be ready..."
    
    # Wait for database
    log_info "Waiting for database to be ready..."
    timeout=60
    while [ $timeout -gt 0 ]; do
        if docker-compose exec -T postgres pg_isready -U postgres >/dev/null 2>&1; then
            log_info "Database is ready"
            break
        fi
        sleep 2
        timeout=$((timeout - 2))
    done
    
    if [ $timeout -le 0 ]; then
        log_error "Database failed to start within timeout"
        exit 1
    fi
    
    # Wait for application
    log_info "Waiting for application to be ready..."
    timeout=60
    while [ $timeout -gt 0 ]; do
        if curl -f http://localhost:8080/health >/dev/null 2>&1; then
            log_info "Application is ready"
            break
        fi
        sleep 2
        timeout=$((timeout - 2))
    done
    
    if [ $timeout -le 0 ]; then
        log_error "Application failed to start within timeout"
        exit 1
    fi
    
    log_info "All services are ready"
}

# Run database migrations
run_migrations() {
    log_step "Running database migrations..."
    
    # Set environment variables for migration script
    export DB_HOST=localhost
    export DB_PORT=5432
    export DB_USER=postgres
    export DB_NAME=request_management_system
    
    if [ "$ENVIRONMENT" = "production" ]; then
        # Load production password from environment or .env file
        if [ -f ".env" ]; then
            export DB_PASSWORD=$(grep DB_PASSWORD .env | cut -d '=' -f2)
        fi
    else
        export DB_PASSWORD=postgres123
    fi
    
    # Run migration script
    ./scripts/migrate.sh migrate
    
    log_info "Database migrations completed"
}

# Verify deployment
verify_deployment() {
    log_step "Verifying deployment..."
    
    # Check if containers are running
    log_info "Checking container status..."
    if [ "$ENVIRONMENT" = "production" ]; then
        docker-compose -f docker-compose.yml -f docker-compose.prod.yml ps
    else
        docker-compose ps
    fi
    
    # Check application health
    log_info "Checking application health..."
    response=$(curl -s http://localhost:8080/health)
    if echo "$response" | grep -q '"status":"ok"'; then
        log_info "Application health check passed"
    else
        log_error "Application health check failed"
        exit 1
    fi
    
    # Check database connection
    log_info "Checking database connection..."
    ./scripts/migrate.sh verify
    
    log_info "Deployment verification completed"
}

# Show deployment status
show_status() {
    log_step "Deployment Status"
    
    echo ""
    echo "Environment: $ENVIRONMENT"
    echo "Application URL: http://localhost:8080"
    echo "Health Check: http://localhost:8080/health"
    
    if [ "$ENVIRONMENT" = "development" ]; then
        echo "pgAdmin URL: http://localhost:5050 (admin@example.com / admin123)"
    fi
    
    echo ""
    echo "API Endpoints:"
    echo "  POST /api/v1/auth/register - User registration"
    echo "  POST /api/v1/auth/login - User login"
    echo "  GET  /api/v1/users/profile - Get user profile"
    echo "  GET  /api/v1/requests - List requests"
    echo "  POST /api/v1/requests - Create request"
    echo ""
    
    echo "Default Admin User:"
    echo "  Username: admin"
    echo "  Password: admin123"
    echo ""
    
    echo "To view logs: docker-compose logs -f app"
    echo "To stop: docker-compose down"
}

# Rollback deployment
rollback_deployment() {
    log_step "Rolling back deployment..."
    
    log_warn "This will stop all containers and remove volumes (data will be lost!)"
    read -p "Are you sure you want to rollback? (yes/no): " confirm
    
    if [ "$confirm" = "yes" ]; then
        if [ "$ENVIRONMENT" = "production" ]; then
            docker-compose -f docker-compose.yml -f docker-compose.prod.yml down -v
        else
            docker-compose down -v
        fi
        
        log_info "Rollback completed"
    else
        log_info "Rollback cancelled"
    fi
}

# Show help
show_help() {
    echo "Deployment Script for Request Management System"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  deploy     Deploy the application (default)"
    echo "  build      Build the application only"
    echo "  migrate    Run database migrations only"
    echo "  verify     Verify deployment"
    echo "  status     Show deployment status"
    echo "  rollback   Rollback deployment"
    echo "  help       Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  ENVIRONMENT  Deployment environment (development|production, default: production)"
}

# Main execution
main() {
    case "${1:-deploy}" in
        "deploy")
            check_prerequisites
            build_application
            deploy_application
            wait_for_services
            run_migrations
            verify_deployment
            show_status
            log_info "Deployment completed successfully!"
            ;;
        "build")
            check_prerequisites
            build_application
            log_info "Build completed successfully!"
            ;;
        "migrate")
            run_migrations
            log_info "Migration completed successfully!"
            ;;
        "verify")
            verify_deployment
            log_info "Verification completed successfully!"
            ;;
        "status")
            show_status
            ;;
        "rollback")
            rollback_deployment
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

