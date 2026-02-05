# Backend Intern Coding Test
## Promotion Code Validation

---

## 📌 Business Context

Our platform operates multiple systems that manage **promotion codes** for marketing campaigns, memberships, and user rewards.

Due to system separation and historical reasons, promotion codes are stored in **two independent data sources**:

- **Campaign System**
- **Membership System**

Each system maintains its own list of valid promotion codes.

Your task is to determine whether a given promotion code is **eligible** for use by a customer.

---

## ✅ Definition of an Eligible Promotion Code

A promotion code is considered **eligible** if and only if:

- The code exists in **both** data sources  
  (i.e. it is recognized by both the Campaign System and the Membership System)

This ensures that the promotion is **officially issued** and **approved for customer usage**.

---

## 📂 Input Data

### 1. Promotion Code Data Sources

You are provided with two large text files:

- `campaign_codes.txt`
- `membership_codes.txt`

Each file:

- Contains **millions of promotion codes**
- Each line represents **one promotion code**
- Codes are **unique within each file**
- A promotion code:
    - Has a **maximum length of 5 characters**
    - Contains only lowercase English letters (`a` to `z`)

---

### 2. Promotion Code to Validate

- A string `code` representing the promotion code entered by a user
- Guaranteed constraints:
    - `1 ≤ length(code) ≤ 5`
    - Characters range from `a` to `z`

---

## 📤 Output

Return:

- `true` if the promotion code exists in **both** data sources
- `false` otherwise

---

## 📘 Example

### campaign_codes.txt
```
abc
xyz
promo
sale
```

### membership_codes.txt
```
abc
gold
promo
```

### Input
```
code = "promo"
```

### Output
```
true
```

---

## 📝 Notes

- Assume the input files may not fit entirely into memory
- Focus on correctness first, then performance
- Clearly explain your design decisions and trade-offs in your solution

---

## 🚀 Submission Requirements

- Source code implementing the validation logic
- A short explanation covering:
    - Your approach
    - Data structures used
    - Performance considerations
- Any assumptions made during implementation

---

Good luck, and we look forward to reviewing your solution!
