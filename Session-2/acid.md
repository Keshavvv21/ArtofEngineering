# ACID Properties in Databases

ACID stands for **Atomicity, Consistency, Isolation, Durability**.  
These four properties guarantee **reliable database transactions** in relational systems.  
Without ACID, banking, e-commerce, and mission-critical applications would risk **data corruption** and **inconsistencies**.

---

## 🔹 Atomicity
**Definition:**  
All operations in a transaction are treated as a **single indivisible unit**. Either all succeed, or none do.  

**Real-world example:**  
Money transfer between two accounts: debit from one account **must** be paired with credit to another. If one fails, both are rolled back.

**Implementation:**  
- Transaction logs  
- Rollback mechanisms  
- Two-phase commit protocols  

```mermaid
flowchart TD
    A[Start Transaction] --> B[Debit Account A]
    B --> C{Success?}
    C -- Yes --> D[Credit Account B]
    C -- No --> R[Rollback Transaction]
    D --> E[Commit Transaction]
