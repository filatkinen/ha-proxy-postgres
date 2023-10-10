FROM postgres
ENV  POSTGRES_PASSWORD=pass
COPY  ./deployment/pgsql/postgres_dump.sql /docker-entrypoint-initdb.d/
EXPOSE 5432
