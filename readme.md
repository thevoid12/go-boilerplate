# my go+pgsql+ tailwind+ htmx+alphine boilerplate

## tech stack
- pgsql https://www.youtube.com/watch?v=AeIksLEHp8E
- viper
- go
- goose (migration) https://www.youtube.com/watch?v=fA8QK69zwlw
- 
## for pgsql
- first install pgsql in your machine
mac: brew install postgresql
ubuntu: sudo apt update && sudo apt install postgresql postgresql-contrib -y
- check if it is successfully installed
psql --version

- sqlc for postgressql https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html
- on mac: brew services start postgresql@14
- psql -h localhost -U $(whoami) -d postgres (now we have a user postgres which is the default user)
<!-- - createdb ( creates default database so from the next time onwards running psql connects to pgsql) -->
- alter USER postgres WITH password 'postgres';
- \q to quit postgres
- brew services restart postgresql@14
- now we have the default user postgres. next step is to create our db
- createdb void.pgsql (since pgsql is running in the background it gets created)
---
If you're starting from scratch and using **Goose** for database migrations in Golang, here's a step-by-step guide to set up PostgreSQL and integrate it with Goose:

---

### **Step 1: Install PostgreSQL on macOS**
1. Install PostgreSQL using Homebrew:
   ```bash
   brew install postgresql@14
   ```
2. Start the PostgreSQL service:
   ```bash
   brew services start postgresql@14
   ```

---

### **Step 2: Set Up PostgreSQL**
1. Connect to the default `postgres` database:
   ```bash
   psql -h localhost -U $(whoami) -d postgres
   ```
2. Create a new user (if needed) and set a password:
   ```sql
   CREATE USER your_user WITH PASSWORD 'your_password';
   ```
3. Create a new database for your application:
   ```sql
   CREATE DATABASE your_database_name;
   ```
4. Grant privileges to the user:
   ```sql
   GRANT ALL PRIVILEGES ON DATABASE your_database_name TO your_user;
   ```
5. Exit PostgreSQL:
   ```sql
   \q
   ```

---

### **Step 3: Install Goose**
https://github.com/pressly/goose
Goose is a database migration tool for Go. Install it using:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```
Verify the installation:
```bash
goose -version
```

---

### **Step 4: Set Up Your Go Project**
1. Initialize a Go module (if not already done):
   ```bash
   go mod init your_project_name
   ```
2. Install the Goose package:
   ```bash
   go get github.com/pressly/goose/v3
   ```

---

### **Step 5: Create Migration Files**
1. Create a directory for migrations:
   ```bash
   mkdir -p migrations
   ```
2. Generate a new migration file:
   ```bash
   goose create create_users_table sql
   ```
   This will create a file like `migrations/00001_create_users_table.sql` with `up` and `down` sections.

3. Write your SQL schema in the `up` section:
   ```sql
   -- +goose Up
   CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       name TEXT NOT NULL,
       email TEXT UNIQUE NOT NULL
   );
   ```
4. Write the rollback SQL in the `down` section:
   ```sql
   -- +goose Down
   DROP TABLE users;
   ```

---

### **Step 6: Configure Goose**
1. Create a `.env` file to store your database connection details:
   ```bash
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_user
   DB_PASSWORD=your_password
   DB_NAME=your_database_name
   DB_SSLMODE=disable
   ```
2. Load the environment variables in your Go application:
   ```bash
   export $(cat .env | xargs)
   ```

---

### **Step 7: Run Migrations**
1. Apply migrations:
   ```bash
   goose -dir migrations postgres "user=your_user password=your_password dbname=your_database_name sslmode=disable" up
   ```
   eg:
   ```bash
   goose -dir migrations postgres "user=postgres password=postgres dbname=void sslmode=disable" up
   ```
2. Rollback migrations (if needed):
   ```bash
   goose -dir migrations postgres "user=your_user password=your_password dbname=your_database_name sslmode=disable" down
   ```

---

### **Step 8: Integrate Goose with Your Go Application**
1. Use the `goose` package in your Go code to manage migrations programmatically:
   ```go
   package main

   import (
       "database/sql"
       "log"
       "os"

       _ "github.com/lib/pq"
       "github.com/pressly/goose/v3"
   )

   func main() {
       // Load environment variables
       dbUser := os.Getenv("DB_USER")
       dbPassword := os.Getenv("DB_PASSWORD")
       dbName := os.Getenv("DB_NAME")
       dbHost := os.Getenv("DB_HOST")
       dbPort := os.Getenv("DB_PORT")

       // Create connection string
       connStr := "user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " host=" + dbHost + " port=" + dbPort + " sslmode=disable"

       // Open database connection
       db, err := sql.Open("postgres", connStr)
       if err != nil {
           log.Fatalf("Failed to open database: %v", err)
       }
       defer db.Close()

       // Run migrations
       if err := goose.Up(db, "migrations"); err != nil {
           log.Fatalf("Failed to run migrations: %v", err)
       }

       log.Println("Migrations applied successfully!")
   }
   ```

---

### **Step 9: Verify the Database**
1. Connect to your database:
   ```bash
   psql -h localhost -U your_user -d your_database_name
   ```
2. Check the tables:
   ```sql
   \dt
   ```
   You should see the `users` table created by your migration.

---

### **Step 10: Automate Migrations**
You can add a `Makefile` to automate common tasks:
```Makefile
migrate-up:
    goose -dir migrations postgres "user=your_user password=your_password dbname=your_database_name sslmode=disable" up

migrate-down:
    goose -dir migrations postgres "user=your_user password=your_password dbname=your_database_name sslmode=disable" down
```
Run migrations with:
```bash
make migrate-up
```

---

This setup ensures a clean and reproducible way to manage PostgreSQL and migrations using Goose in your Golang project. Let me know if you need further clarification!
