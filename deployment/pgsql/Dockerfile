FROM postgres
ENV  POSTGRES_PASSWORD=pass
COPY  postgres_dump.sql /docker-entrypoint-initdb.d/
EXPOSE 5432
