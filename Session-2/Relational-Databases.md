# 🗄️ Relational Databases and AWS RDS

---

## 📖 What are Relational Databases?

Relational databases are a **foundational technology** for storing and organizing **structured data**.  

- Data is stored in **tables (relations)** with **rows and columns**.  
- Tables are linked via **primary keys** and **foreign keys**.  
- Managed through **SQL (Structured Query Language)**.  
- Ensure **data integrity, consistency, and efficient retrieval**.

---

## ✅ Why are Relational Databases Used?

- **Data Integrity** → Enforced by schemas & constraints.  
- **Consistency (ACID Properties)** → Transactions remain reliable even during failures.  
- **Structured Querying** → SQL enables complex queries and analytics.  
- **Maturity & Community** → Decades of adoption, strong ecosystem.  

📌 **Use cases**:  
- Financial systems 🏦  
- E-commerce platforms 🛒  
- Inventory management 📦  
- CRM (Customer Relationship Management) 🤝  
- ERP (Enterprise Resource Planning) 🏭  

---

## 🛒 Example: E-commerce Relational Database Schema

Below is a simple schema with **Users, Orders, and Products** tables:

### Users Table

| UserID (PK) | Name       | Email              |
|-------------|------------|--------------------|
| 1           | Alice Doe  | alice@email.com    |
| 2           | Bob Smith  | bob@email.com      |

### Products Table

| ProductID (PK) | ProductName  | Price |
|----------------|--------------|-------|
| 101            | Laptop       | 60000 |
| 102            | Phone        | 15000 |

### Orders Table

| OrderID (PK) | UserID (FK → Users.UserID) | ProductID (FK → Products.ProductID) | Quantity | OrderDate  |
|--------------|----------------------------|--------------------------------------|----------|------------|
| 5001         | 1                          | 101                                  | 1        | 2025-01-01 |
| 5002         | 2                          | 102                                  | 2        | 2025-01-02 |

👉 Here, **Users and Products** are linked to **Orders** via **foreign keys**, ensuring consistency.

---

## ☁️ Amazon RDS: Managed Relational Database Service

Managing relational databases at scale can be **operationally intensive**.  
Amazon **RDS (Relational Database Service)** changes this by offering a **fully managed relational database** in the cloud.

### 🔹 Supported Engines
- MySQL  
- PostgreSQL  
- Oracle  
- SQL Server  
- MariaDB  

### 🔹 Automated Features
- Hardware provisioning  
- Database setup & patching  
- Backups & snapshots  
- Multi-AZ replication for **99.95% availability SLA**  
- Point-in-time recovery (up to 35 days)  

### 🔹 Scalability
- **Read replicas** for scaling read-heavy workloads.  
- **Vertical scaling** (increase instance size).  
- **Horizontal scaling** via sharding + replicas.  

---

## 🌟 Key Benefits of AWS RDS

- 🚀 **Reduced administrative burden** (no manual setup/patching).  
- 🔐 **Enhanced security** (encryption at rest + in transit).  
- 📊 **Monitoring & alerting** via AWS CloudWatch.  
- ⚡ **Seamless scaling** without downtime.  
- 🔄 **Automatic failover** with Multi-AZ deployments.  

---

## 💡 Did You Know?

- The **first relational database model** was introduced by **E.F. Codd in 1970** at IBM.  
- Amazon RDS can scale to **millions of requests per second** while keeping latency low.  

