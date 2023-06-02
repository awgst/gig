package dockercompose

var MySqlDBTemplate = `mysql:
        image: 'mysql:latest'
        container_name: 'gig_mysql'
        ports:
            - '${FORWARD_DB_PORT:-3306}:3306'
        environment:
            MYSQL_ROOT_PASSWORD: '${DB_PASSWORD}'
            MYSQL_ROOT_HOST: "%"
            MYSQL_DATABASE: '${DB_DATABASE}'
            MYSQL_ALLOW_EMPTY_PASSWORD: 'yes'
        volumes:
            - 'gig-mysql:/var/lib/mysql'
            - './database/create-database.sh:/docker-entrypoint-initdb.d/10-create-testing-database.sh'
        networks:
            - gig
        healthcheck:
            test: ["CMD", "mysqladmin", "ping", "-p${DB_PASSWORD}"]
            retries: 3
            timeout: 5s`

var PostgreSqlDBTemplate = `postgresql:
        image: postgres:alpine
        environment:
            - POSTGRES_USER=${DB_USERNAME}
            - POSTGRES_PASSWORD=${DB_PASSWORD}
            - POSTGRES_DB=${DB_DATABASE}
        ports:
            - "${FORWARD_DB_PORT:-5432}:5432"
        volumes:
            - gig-postgresql:/var/lib/postgresql/data
        networks:
            - gig
        healthcheck:
            test: ["CMD", "pg_isready", "-U", "${DB_USERNAME}"]
            retries: 3
            timeout: 5s`
