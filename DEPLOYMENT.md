# Deployment Guide

This guide provides comprehensive instructions for deploying the Request Management System in various environments, from development to production.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Development Deployment](#development-deployment)
- [Production Deployment](#production-deployment)
- [Cloud Deployment](#cloud-deployment)
- [Security Considerations](#security-considerations)
- [Monitoring and Maintenance](#monitoring-and-maintenance)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### System Requirements

#### Minimum Requirements
- **CPU**: 2 cores
- **RAM**: 4GB
- **Storage**: 20GB available space
- **Network**: Internet connection for dependencies

#### Recommended Requirements
- **CPU**: 4+ cores
- **RAM**: 8GB+
- **Storage**: 50GB+ SSD
- **Network**: High-speed internet connection

### Software Dependencies

#### Required Software
- **Docker**: Version 20.10 or higher
- **Docker Compose**: Version 2.0 or higher
- **Git**: For source code management

#### Optional Software
- **Go**: Version 1.21+ (for development)
- **PostgreSQL Client**: For database management
- **Nginx**: For reverse proxy (production)

### Installation Commands

#### Ubuntu/Debian
```bash
# Update package list
sudo apt update

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Install Git
sudo apt install git

# Logout and login to apply Docker group changes
```

#### CentOS/RHEL
```bash
# Install Docker
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install docker-ce docker-ce-cli containerd.io
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Install Git
sudo yum install git
```

#### macOS
```bash
# Install Docker Desktop from https://www.docker.com/products/docker-desktop
# Or using Homebrew
brew install --cask docker

# Install Git (if not already installed)
brew install git
```

## Development Deployment

### Quick Start

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd request-management-system
   ```

2. **Start development environment**:
   ```bash
   # Using Docker Compose
   docker-compose up -d
   
   # Or using Make
   make docker-run
   ```

3. **Verify deployment**:
   ```bash
   # Check health endpoint
   curl http://localhost:8080/health
   
   # Check container status
   docker-compose ps
   ```

4. **Access services**:
   - **API**: http://localhost:8080/api/v1
   - **Health Check**: http://localhost:8080/health
   - **pgAdmin**: http://localhost:5050 (admin@example.com / admin123)

### Development Configuration

The development environment uses default configurations suitable for local development:

```yaml
# docker-compose.yml (relevant sections)
services:
  postgres:
    environment:
      POSTGRES_PASSWORD: postgres123
    ports:
      - "5432:5432"
  
  app:
    environment:
      APP_ENV: development
      LOG_LEVEL: info
      JWT_SECRET: development-secret-key
    ports:
      - "8080:8080"
```

### Development Tools

#### Enable pgAdmin
```bash
# Start with pgAdmin
docker-compose --profile tools up -d

# Access pgAdmin at http://localhost:5050
# Email: admin@example.com
# Password: admin123
```

#### Live Reload Development
```bash
# Install Air for live reload
go install github.com/cosmtrek/air@latest

# Start development with live reload
make dev
```

#### Database Management
```bash
# Reset database
make db-reset

# Run migrations manually
./scripts/migrate.sh migrate

# Verify database
./scripts/migrate.sh verify
```

## Production Deployment

### Pre-deployment Checklist

- [ ] Server meets minimum requirements
- [ ] Docker and Docker Compose installed
- [ ] SSL certificates obtained (if using HTTPS)
- [ ] Environment variables configured
- [ ] Database backup strategy planned
- [ ] Monitoring solution prepared
- [ ] Log aggregation configured

### Step-by-Step Production Deployment

#### 1. Server Preparation

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install required packages
sudo apt install -y curl wget git ufw

# Configure firewall
sudo ufw allow ssh
sudo ufw allow 80
sudo ufw allow 443
sudo ufw --force enable

# Create application user
sudo useradd -m -s /bin/bash appuser
sudo usermod -aG docker appuser
```

#### 2. Application Setup

```bash
# Switch to application user
sudo su - appuser

# Clone repository
git clone <repository-url>
cd request-management-system

# Create production environment file
cp .env.example .env
```

#### 3. Environment Configuration

Edit the `.env` file with production values:

```bash
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database Configuration (CHANGE THESE!)
DB_PASSWORD=your_secure_database_password_here

# JWT Configuration (CHANGE THIS!)
JWT_SECRET=your_super_secure_jwt_secret_key_here_minimum_32_characters

# Application Configuration
APP_ENV=production
LOG_LEVEL=warn
```

#### 4. SSL Certificate Setup (Optional)

```bash
# Create SSL directory
mkdir -p ssl

# Copy your SSL certificates
cp /path/to/your/cert.pem ssl/
cp /path/to/your/key.pem ssl/

# Or use Let's Encrypt
sudo apt install certbot
sudo certbot certonly --standalone -d yourdomain.com
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem ssl/cert.pem
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem ssl/key.pem
sudo chown appuser:appuser ssl/*
```

#### 5. Deploy Application

```bash
# Deploy using script
ENVIRONMENT=production ./scripts/deploy.sh

# Or deploy manually
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

#### 6. Verify Production Deployment

```bash
# Check container status
docker-compose -f docker-compose.yml -f docker-compose.prod.yml ps

# Check application health
curl http://localhost:8080/health

# Check logs
docker-compose -f docker-compose.yml -f docker-compose.prod.yml logs -f app
```

### Production with Nginx Reverse Proxy

#### 1. Enable Nginx Service

```bash
# Start with Nginx
docker-compose -f docker-compose.yml -f docker-compose.prod.yml --profile nginx up -d
```

#### 2. Configure Custom Domain

Update `nginx.conf` with your domain:

```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;
    
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    # ... rest of configuration
}
```

#### 3. DNS Configuration

Point your domain to the server IP:

```
A    yourdomain.com    YOUR_SERVER_IP
A    www.yourdomain.com    YOUR_SERVER_IP
```

### Environment Variables Reference

#### Required Production Variables

```bash
# Database (MUST CHANGE)
DB_PASSWORD=secure_password_here

# JWT (MUST CHANGE)
JWT_SECRET=secure_jwt_secret_minimum_32_characters

# Application
APP_ENV=production
LOG_LEVEL=warn
```

#### Optional Production Variables

```bash
# Server
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_IDLE_TIMEOUT=120s

# Database
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=300s

# Security
CORS_ALLOWED_ORIGINS=https://yourdomain.com
RATE_LIMIT_ENABLED=true
```

## Cloud Deployment

### AWS Deployment

#### Using EC2

1. **Launch EC2 Instance**:
   ```bash
   # Recommended: t3.medium or larger
   # OS: Ubuntu 20.04 LTS
   # Security Group: Allow ports 22, 80, 443
   ```

2. **Connect and Setup**:
   ```bash
   ssh -i your-key.pem ubuntu@your-ec2-ip
   
   # Follow production deployment steps
   ```

3. **Configure Load Balancer** (Optional):
   ```bash
   # Create Application Load Balancer
   # Target Group: Port 8080
   # Health Check: /health
   ```

#### Using ECS (Elastic Container Service)

1. **Build and Push Image**:
   ```bash
   # Build image
   docker build -t request-management-system .
   
   # Tag for ECR
   docker tag request-management-system:latest 123456789012.dkr.ecr.region.amazonaws.com/request-management-system:latest
   
   # Push to ECR
   docker push 123456789012.dkr.ecr.region.amazonaws.com/request-management-system:latest
   ```

2. **Create ECS Task Definition**:
   ```json
   {
     "family": "request-management-system",
     "networkMode": "awsvpc",
     "requiresCompatibilities": ["FARGATE"],
     "cpu": "512",
     "memory": "1024",
     "containerDefinitions": [
       {
         "name": "app",
         "image": "123456789012.dkr.ecr.region.amazonaws.com/request-management-system:latest",
         "portMappings": [
           {
             "containerPort": 8080,
             "protocol": "tcp"
           }
         ],
         "environment": [
           {
             "name": "DB_HOST",
             "value": "your-rds-endpoint"
           }
         ]
       }
     ]
   }
   ```

#### Using RDS for Database

```bash
# Create RDS PostgreSQL instance
aws rds create-db-instance \
  --db-instance-identifier request-management-db \
  --db-instance-class db.t3.micro \
  --engine postgres \
  --master-username postgres \
  --master-user-password your-secure-password \
  --allocated-storage 20 \
  --vpc-security-group-ids sg-xxxxxxxxx
```

### Google Cloud Platform

#### Using Compute Engine

```bash
# Create VM instance
gcloud compute instances create request-management-vm \
  --image-family=ubuntu-2004-lts \
  --image-project=ubuntu-os-cloud \
  --machine-type=e2-medium \
  --tags=http-server,https-server

# SSH and deploy
gcloud compute ssh request-management-vm
```

#### Using Cloud Run

```bash
# Build and deploy
gcloud builds submit --tag gcr.io/PROJECT_ID/request-management-system
gcloud run deploy --image gcr.io/PROJECT_ID/request-management-system --platform managed
```

### Digital Ocean

#### Using Droplets

```bash
# Create droplet via web interface or API
# Choose Ubuntu 20.04, minimum 2GB RAM

# SSH and deploy
ssh root@your-droplet-ip
```

#### Using App Platform

```yaml
# app.yaml
name: request-management-system
services:
- name: api
  source_dir: /
  github:
    repo: your-username/request-management-system
    branch: main
  run_command: ./main
  environment_slug: go
  instance_count: 1
  instance_size_slug: basic-xxs
  envs:
  - key: DB_HOST
    value: ${db.HOSTNAME}
  - key: DB_PASSWORD
    value: ${db.PASSWORD}
databases:
- name: db
  engine: PG
  version: "13"
```

## Security Considerations

### SSL/TLS Configuration

#### Let's Encrypt (Free SSL)

```bash
# Install Certbot
sudo apt install certbot

# Obtain certificate
sudo certbot certonly --standalone -d yourdomain.com

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

#### Custom SSL Certificate

```bash
# Copy certificates to ssl/ directory
cp your-cert.pem ssl/cert.pem
cp your-key.pem ssl/key.pem

# Set proper permissions
chmod 600 ssl/key.pem
chmod 644 ssl/cert.pem
```

### Firewall Configuration

```bash
# Ubuntu UFW
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw deny 8080/tcp  # Block direct access to app
sudo ufw enable

# CentOS/RHEL Firewalld
sudo firewall-cmd --permanent --add-service=ssh
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

### Database Security

```bash
# Change default passwords
DB_PASSWORD=very_secure_password_here

# Restrict database access
# In docker-compose.prod.yml, remove ports mapping for postgres
# Only allow internal container communication
```

### Application Security

```bash
# Use strong JWT secret
JWT_SECRET=your_super_secure_jwt_secret_minimum_32_characters

# Enable rate limiting
RATE_LIMIT_ENABLED=true

# Restrict CORS origins
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

## Monitoring and Maintenance

### Health Monitoring

#### Basic Health Check

```bash
# Create health check script
cat > /home/appuser/health-check.sh << 'EOF'
#!/bin/bash
response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
if [ $response -eq 200 ]; then
    echo "$(date): Application is healthy"
else
    echo "$(date): Application is unhealthy (HTTP $response)"
    # Restart application
    cd /home/appuser/request-management-system
    docker-compose restart app
fi
EOF

chmod +x /home/appuser/health-check.sh

# Add to crontab
crontab -e
# Add: */5 * * * * /home/appuser/health-check.sh >> /var/log/health-check.log 2>&1
```

#### Advanced Monitoring with Prometheus

```yaml
# Add to docker-compose.yml
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    profiles:
      - monitoring

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    profiles:
      - monitoring
```

### Log Management

#### Log Rotation

```bash
# Create logrotate configuration
sudo cat > /etc/logrotate.d/request-management << 'EOF'
/var/log/request-management/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 appuser appuser
    postrotate
        docker-compose -f /home/appuser/request-management-system/docker-compose.yml restart app
    endscript
}
EOF
```

#### Centralized Logging

```yaml
# Add to docker-compose.yml
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:7.15.0
    environment:
      - discovery.type=single-node
    profiles:
      - logging

  logstash:
    image: docker.elastic.co/logstash/logstash:7.15.0
    profiles:
      - logging

  kibana:
    image: docker.elastic.co/kibana/kibana:7.15.0
    ports:
      - "5601:5601"
    profiles:
      - logging
```

### Database Maintenance

#### Automated Backups

```bash
# Create backup script
cat > /home/appuser/backup-db.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/home/appuser/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/request_management_$DATE.sql"

mkdir -p $BACKUP_DIR

# Create backup
docker-compose exec -T postgres pg_dump -U postgres request_management_system > $BACKUP_FILE

# Compress backup
gzip $BACKUP_FILE

# Remove backups older than 30 days
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete

echo "$(date): Database backup completed: $BACKUP_FILE.gz"
EOF

chmod +x /home/appuser/backup-db.sh

# Schedule daily backups
crontab -e
# Add: 0 2 * * * /home/appuser/backup-db.sh >> /var/log/backup.log 2>&1
```

#### Database Restore

```bash
# Restore from backup
gunzip backup_file.sql.gz
docker-compose exec -T postgres psql -U postgres -d request_management_system < backup_file.sql
```

### Application Updates

#### Rolling Updates

```bash
# Create update script
cat > /home/appuser/update-app.sh << 'EOF'
#!/bin/bash
cd /home/appuser/request-management-system

# Pull latest changes
git pull origin main

# Build new image
docker-compose build app

# Rolling update
docker-compose up -d --no-deps app

# Verify health
sleep 10
response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
if [ $response -eq 200 ]; then
    echo "$(date): Update successful"
    # Clean up old images
    docker image prune -f
else
    echo "$(date): Update failed, rolling back"
    docker-compose restart app
fi
EOF

chmod +x /home/appuser/update-app.sh
```

## Troubleshooting

### Common Issues

#### Application Won't Start

```bash
# Check container logs
docker-compose logs app

# Check container status
docker-compose ps

# Check system resources
df -h
free -m
docker system df
```

#### Database Connection Issues

```bash
# Check database container
docker-compose logs postgres

# Test database connection
docker-compose exec postgres psql -U postgres -d request_management_system -c "SELECT 1;"

# Check environment variables
docker-compose exec app env | grep DB_
```

#### Performance Issues

```bash
# Check resource usage
docker stats

# Check database performance
docker-compose exec postgres psql -U postgres -d request_management_system -c "
SELECT query, calls, total_time, mean_time 
FROM pg_stat_statements 
ORDER BY total_time DESC 
LIMIT 10;"

# Check application metrics
curl http://localhost:8080/health
```

### Debug Mode

```bash
# Enable debug logging
# In .env file:
LOG_LEVEL=debug

# Restart application
docker-compose restart app

# View debug logs
docker-compose logs -f app
```

### Recovery Procedures

#### Complete System Recovery

```bash
# Stop all services
docker-compose down

# Remove volumes (WARNING: This will delete all data)
docker-compose down -v

# Restore from backup
./scripts/migrate.sh reset
# Restore database from backup file

# Restart services
docker-compose up -d
```

#### Partial Recovery

```bash
# Restart specific service
docker-compose restart app
docker-compose restart postgres

# Rebuild and restart
docker-compose up -d --build app
```

### Performance Tuning

#### Database Optimization

```sql
-- Add indexes for better performance
CREATE INDEX CONCURRENTLY idx_requests_client_id ON requests(client_id);
CREATE INDEX CONCURRENTLY idx_requests_status ON requests(status);
CREATE INDEX CONCURRENTLY idx_requests_created_at ON requests(created_at);
CREATE INDEX CONCURRENTLY idx_users_role ON users(role);
```

#### Application Optimization

```bash
# Increase connection pool size
# In .env file:
DB_MAX_OPEN_CONNS=50
DB_MAX_IDLE_CONNS=10

# Tune server timeouts
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_IDLE_TIMEOUT=120s
```

### Support and Maintenance

#### Regular Maintenance Tasks

1. **Weekly**:
   - Check application logs
   - Verify backup integrity
   - Monitor disk space
   - Review security logs

2. **Monthly**:
   - Update system packages
   - Review database performance
   - Clean up old Docker images
   - Test disaster recovery procedures

3. **Quarterly**:
   - Security audit
   - Performance review
   - Capacity planning
   - Update documentation

#### Emergency Contacts

- **System Administrator**: admin@yourcompany.com
- **Database Administrator**: dba@yourcompany.com
- **Security Team**: security@yourcompany.com
- **On-call Support**: +1-555-SUPPORT

---

This deployment guide provides comprehensive instructions for deploying the Request Management System in various environments. For additional support, refer to the [README](README.md) and [API Documentation](API.md).

