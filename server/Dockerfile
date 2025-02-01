FROM golang:1.23

WORKDIR /app

COPY go.mod ./
COPY src/ ./src/
COPY docs/ ./docs/
# COPY .env ./
RUN go mod tidy
RUN go mod vendor

COPY *.go ./
# ENV POSTGRES_URL=postgresql://hrms_owho_user:AJEqyOETLDO3YEMvFNX131A3sCEDOm5T@dpg-cqj3gkeehbks73c5748g-a.oregon-postgres.render.com/hrms_owho

RUN go build -o /server

EXPOSE 8080

CMD [ "/server" ]