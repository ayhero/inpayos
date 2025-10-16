#!/bin/bash

# InPayOS Deployment Script
# 支付网关部署脚本
# Usage: ./scripts/deploy.sh [dev|prod] [options]

set -e

# 脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# 默认配置
DEFAULT_ENV="dev"
DEFAULT_ACTION="deploy"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    cat << EOF
InPayOS Deployment Script

Usage: $0 [ENVIRONMENT] [ACTION] [OPTIONS]

ENVIRONMENT:
    dev         Deploy to development environment (default)
    prod        Deploy to production environment

ACTION:
    deploy      Deploy services (default)
    stop        Stop services
    restart     Restart services
    logs        Show service logs
    status      Show service status
    clean       Clean up containers and images

OPTIONS:
    --build     Force rebuild Docker images
    --pull      Pull latest images before deploy
    --backup    Create backup before deployment (prod only)
    --help      Show this help message

Examples:
    $0 dev                          # Deploy to dev environment
    $0 prod deploy --build          # Deploy to prod with rebuild
    $0 dev restart                  # Restart dev services
    $0 prod logs                    # Show prod logs
    $0 dev stop                     # Stop dev services
    $0 prod clean                   # Clean prod environment

EOF
}

# 检查必要条件
check_prerequisites() {
    log_info "检查部署环境..."
    
    # 检查 Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装或不在 PATH 中"
        exit 1
    fi
    
    # 检查 Docker Compose
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        log_error "Docker Compose 未安装"
        exit 1
    fi
    
    # 配置文件存在性检查
    if [[ ! -f "$PROJECT_DIR/${ENVIRONMENT}.yaml" ]]; then
        log_warning "配置文件 ${ENVIRONMENT}.yaml 不存在，请确保配置正确"
    fi
    
    log_success "环境检查通过"
}

# 创建必要目录
    create_directories() {
    log_info "创建必要目录..."
    
    mkdir -p "$PROJECT_DIR/logs"
    
    log_success "目录创建完成"
}

# 备份数据 (仅生产环境)
backup_data() {
    if [[ "$ENVIRONMENT" == "prod" && "$BACKUP" == "true" ]]; then
        log_info "创建生产环境备份..."
        
        BACKUP_DIR="$PROJECT_DIR/backups/$(date +%Y%m%d_%H%M%S)"
        mkdir -p "$BACKUP_DIR"
        
        # 备份配置文件和日志
        cp "$PROJECT_DIR/${ENVIRONMENT}.yaml" "$BACKUP_DIR/" 2>/dev/null || true
        cp -r "$PROJECT_DIR/logs" "$BACKUP_DIR/" 2>/dev/null || true
        
        log_success "备份完成: $BACKUP_DIR"
    fi
}

# 部署服务
deploy_services() {
    log_info "部署 InPayOS 服务 (环境: $ENVIRONMENT)..."
    
    cd "$PROJECT_DIR"
    
    # 设置环境变量
    export ENV="$ENVIRONMENT"
    
    # 环境变量已设置为 ENV
    # 配置验证由应用启动时处理
    
    # Docker Compose 命令参数
    COMPOSE_ARGS=""
    
    if [[ "$PULL_IMAGES" == "true" ]]; then
        log_info "拉取最新镜像..."
        docker-compose pull
    fi
    
    if [[ "$FORCE_BUILD" == "true" ]]; then
        log_info "强制重新构建镜像..."
        docker-compose build --no-cache
    fi
    
    # 启动服务
    docker-compose up -d $COMPOSE_ARGS
    
    log_success "服务部署完成"
    
    # 服务启动完成
    log_info "服务启动完成，请检查应用日志确认配置正确性"
}

# 检查服务健康状态
check_service_health() {
    log_info "检查服务状态..."
    docker-compose ps
}

# 停止服务
stop_services() {
    log_info "停止 InPayOS 服务 (环境: $ENVIRONMENT)..."
    
    cd "$PROJECT_DIR"
    export ENV="$ENVIRONMENT"
    
    docker-compose down
    
    log_success "服务已停止"
}

# 重启服务
restart_services() {
    log_info "重启 InPayOS 服务 (环境: $ENVIRONMENT)..."
    
    stop_services
    sleep 5
    deploy_services
}

# 查看日志
show_logs() {
    log_info "显示 InPayOS 服务日志 (环境: $ENVIRONMENT)..."
    
    cd "$PROJECT_DIR"
    export ENV="$ENVIRONMENT"
    
    docker-compose logs -f --tail=100
}

# 显示服务状态
show_status() {
    log_info "InPayOS 服务状态 (环境: $ENVIRONMENT):"
    
    cd "$PROJECT_DIR"
    export ENV="$ENVIRONMENT"
    
    docker-compose ps
    
    echo
    log_info "端口状态:"
    netstat -tuln | grep -E ":(8080|8081|8082|8083|8084) " || echo "无相关端口监听"
}

# 清理环境
clean_environment() {
    log_warning "这将删除所有容器和镜像，确定继续吗? (y/N)"
    read -r confirmation
    
    if [[ "$confirmation" != "y" && "$confirmation" != "Y" ]]; then
        log_info "取消清理操作"
        return 0
    fi
    
    log_info "清理 InPayOS 环境 (环境: $ENVIRONMENT)..."
    
    cd "$PROJECT_DIR"
    export ENV="$ENVIRONMENT"
    
    # 停止并删除容器
    docker-compose down --remove-orphans
    
    # 删除相关镜像
    docker images | grep inpayos | awk '{print $3}' | xargs docker rmi -f 2>/dev/null || true
    
    log_success "环境清理完成"
}

# 解析命令行参数
parse_arguments() {
    ENVIRONMENT="$DEFAULT_ENV"
    ACTION="$DEFAULT_ACTION"
    FORCE_BUILD="false"
    PULL_IMAGES="false"
    BACKUP="false"
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            dev|prod)
                ENVIRONMENT="$1"
                shift
                ;;
            deploy|stop|restart|logs|status|clean)
                ACTION="$1"
                shift
                ;;
            --build)
                FORCE_BUILD="true"
                shift
                ;;
            --pull)
                PULL_IMAGES="true"
                shift
                ;;
            --backup)
                BACKUP="true"
                shift
                ;;

            --help|-h)
                show_help
                exit 0
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

# 主函数
main() {
    log_info "InPayOS 部署脚本启动"
    
    parse_arguments "$@"
    
    log_info "环境: $ENVIRONMENT"
    log_info "操作: $ACTION"
    
    check_prerequisites
    
    case "$ACTION" in
        deploy)
            create_directories
            backup_data
            deploy_services
            ;;
        stop)
            stop_services
            ;;
        restart)
            restart_services
            ;;
        logs)
            show_logs
            ;;
        status)
            show_status
            ;;
        clean)
            clean_environment
            ;;
        *)
            log_error "未知操作: $ACTION"
            show_help
            exit 1
            ;;
    esac
    
    log_success "操作完成"
}

# 执行主函数
main "$@"