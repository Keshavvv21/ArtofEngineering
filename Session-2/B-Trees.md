# 🌳 B-Trees: The Backbone of Database Indexing

B-trees are a fundamental **self-balancing tree data structure** designed to maintain sorted data and allow **searches, sequential access, insertions, and deletions** in logarithmic time. Their efficiency is paramount for database systems that rely on disk storage.

---

## 🚀 Why B-Trees Matter

### 🔹 Optimized Disk I/O
- B-trees minimize the number of disk reads and writes.  
- They are structured to store **large blocks of data**, retrievable in a single I/O operation.  
- Since **disk access is slower** than memory access, this optimization is critical.  

### 🔹 Fast Data Retrieval
- Keep data sorted and balanced.  
- Any record can usually be found within **3–4 disk seeks**, even in huge databases.  

### 🔹 Efficient Range Queries
- Excellent for **range lookups** (e.g., `age BETWEEN 20 AND 30`).  
- Only a small portion of the tree must be traversed.  

### 🔹 Dynamic Structure
- Auto-balances on insertions and deletions.  
- Consistent performance without costly reorganizations.  

---

## 📖 Did You Know?
The **“B”** in B-tree originally stood for **Boeing** (developed at Boeing Scientific Research Labs).  
But people also say it means **Balanced, Broad, or Bushy** 🌱.

---

## 🏗 Understanding B-Tree Structure

A **B-tree of order m**:
- Each node can have up to **m children**.  
- Each internal node (except root) must have at least **ceil(m/2) children**.  
- Leaf nodes all exist at the **same depth** → ensures balance.  

### 📊 Diagram (Mermaid)

```mermaid
graph TD
    A[Root Node] --> B1[Key 10]
    A --> B2[Key 20]
    A --> B3[Key 30]
    B1 --> L1[Leaf Block with data <10]
    B2 --> L2[Leaf Block with 10–20]
    B3 --> L3[Leaf Block with >20]
