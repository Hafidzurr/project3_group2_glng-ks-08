# PROJECT 3 - KELOMPOK 2 - HACKTIV8 - MBKM - GOLANG FOR BACK-END

## Team 2-KS08 Contributors
* Hafidzurrohman Saifullah (GLNG-KS-08-02) - [GitHub@Hafidzurr](https://github.com/Hafidzurr) - Golang For Back-End - Universitas Gunadarma
* Sherly Fauziyah Syaharani (GLNG-KS-08-018) - [GitHub@Sherlyfauz](https://github.com/Sherlyfauz) - Golang For Back-End - Universitas Merdeka Malang 
* Timotius Winsen Bastian (GLNG-KS-08-016) - [GitHub@Kozzen890](https://github.com/Kozzen890) - Golang For Back-End - Universitas Dian Nuswantoro 
##
##
## API URL 
#### https://project3group2glng-ks-08-production.up.railway.app/
##
## Postman Documentation
#### https://documenter.getpostman.com/view/24258835/2s9YeA8tiu
##
## Sytem Requirement
* Golang.
* Postgres SQL.
## Installation Local
#### 1. Open terminal or command prompt
```
git clone https://github.com/Hafidzurr/project1_group2_glng-ks-08.git
cd project1_group2_glng-ks-08
go mod tidy
```
#### 2. Setting Database 

##### Create database in postgres SQL with name `kanban_board` or you can change whats name you like, but coution here you must change database name in `db.go` too.

##### Go to db.go, comment line code from `dns = fmt.Sprintf` - `dbname, dbPort)` and uncomment line code `dsn = "host=host...`.

##### change your `db credential` in `db.go`.


#### 3. Run 
```
go run main.go
```
## Installation and Deploying to Railway
#### 1. Open terminal or command prompt
```
git clone https://github.com/Hafidzurr/project1_group2_glng-ks-08.git
cd project1_group2_glng-ks-08
go mod tidy
```
#### 2. Push into Your New Repo
##### Create a New Repository in Your Github Account
##### Change the Remote URL
```
git remote set-url origin https://github.com/new_user/new_repo.git
```
##### Push to the New Repository 
```
git push -u origin master or your name for repo banch
```
#### 3. Create Account Railway using your github Account and Login
##### Create `New Project` -> Choose `Deploy from github Repo -> Choose `Your Repo Name` -> Wait Deploying Untill Getting Error

#### 4. Adding Postgres SQL into Your Project




##### Create database in postgres SQL with name `kanban_board` or you can change whats name you like, but coution here you must change database name in `db.go` too.

##### Go to db.go, comment line code from `dns = fmt.Sprintf` - `dbname, dbPort)` and uncomment line code `dsn = "host=host...`

##### change your `db credential` in `db.go`

#### 3. Run 
```
go run main.go
```
