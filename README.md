# RumahResep Project
[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example)
[![Go.Dev reference](https://img.shields.io/badge/gorm-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/gorm.io/gorm?tab=doc)
[![Go.Dev reference](https://img.shields.io/badge/echo-reference-blue?logo=go&logoColor=white)](https://github.com/labstack/echo)
[![Go Reference](https://img.shields.io/badge/midtrans-reference-blue?logo=Midtrans&logoColor=white)](https://github.com/Midtrans/midtrans-go)
[![Go Reference](https://img.shields.io/badge/gmaps-reference-blue?logo=GMaps&logoColor=white)](https://github.com/googlemaps/google-maps-services-go)

# Table of Content
- [Description](#description)
- [How to Use](#how-to-use)
- [Database Schema](#database-schema)
- [Testing Coverage]($testing-coverage)
- [Feature](#feature)
- [Endpoints](#endpoints)
- [Credits](#credits)

# Description
RumahResep merupakan sebuah aplikasi penjualan berbagai macam resep berserta bahan yang menjadi satu. Disini Customer dapat membeli resep dan bahan sekaligus sehingga customer tidak perlu lagi mencari bahan dalam membuat suatu resep. Pada aplikasi ini resep dan bahan ditambahkan oleh penjual itu sendiri yang bertindak sebagai Admin. Untuk supply bahan untuk memenuhi suatu resep, Admin berkerja sama dengan beberapa supplier yang disini disebut sebagai Seller. Seller disini dapat melakukan restock bahan sesuai resep yang nantinya akan dikirim kepada Customer yang melakukan pembelian.

# Database Schema
![ERD](https://github.com/ridzqihammam17/RumahResep-Project/blob/main/screenshoot/erd.png)

# Testing Coverage
Implement Unit Testing average above 85%

![TESTING](https://github.com/ridzqihammam17/RumahResep-Project/blob/main/screenshoot/testing_coverage.PNG)

# Feature
List of overall feature in this Project (To get more details see the API Documentation below)
| No.| Feature        | Role                     | Keterangan                                                             |
| :- | :------------- | :----------------------- | :--------------------------------------------------------------------- |
| 1. | Register       | Admin, Seller, Customer  | Authentication process                                                 |
| 2. | Login          | Admin, Seller, Customer  | Authentication process                                                 |
| 3. | Read Recipe    | Admin, Seller, Customer  | Get all recipe and get details of recipe                               |
| 4. | CUD Recipe     | Admin                    | Create, Update, and Delete Recipe in system                            |
| 4. | CRUD Ingredient| Admin                    | Create, Read, Update, and Delete Ingredient of the recipe in system    |
| 5. | CRUD Category  | Admin                    | Create, Read, Update, and Delete Category in system                    |
| 6. | Cart           | Customer                 | Add recipe to Cart, see Recipe list on Cart, and get details of Cart   |
| 7. | Checkout       | Customer                 | Checkout cart (list recipe item) to transaction and make payment       |
| 8. | Payment        | Customer                 | Make payment and get status payment of transaction                     |
| 9. | Restock        | Seller                   | Restock the ingredient of the recipe                                   |

# How to Use
- Install Go and Database MySQL/XAMPP
- Clone this repository in your $PATH:
```
$ git clone https://github.com/ridzqihammam17/RumahResep-Project.git
```
- Create file .env based on this project 
``
sample-env
``
- Don't forget to create database name as you want in your MySQL
- Run program with command
```
go run main.go
```
# Endpoints
Read the API documentation here [API Endpoint Documentation](https://app.swaggerhub.com/apis/ridzqihammam17/rumah-resep/1.0.0) (Swagger)

# Credits
- [Muchammad Abdurrochman](https://github.com/Abdurrochman25) (Author)
- [Fabrian Ivan Prasetya](https://github.com/fabrianivan-id) (Author)
- [Ridzqiawan Hammam](https://github.com/ridzqihammam17) (Author)
