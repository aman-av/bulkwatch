# 📘 Skillbook: Problem-Solving Playbook

A structured guide to approach, analyze, and solve any user request effectively and consistently.

---

## 🎯 Purpose

This Skillbook defines a **standard thinking framework** to:
- Understand user intent
- Break down problems
- Produce clear, useful outputs
- Avoid confusion and inefficiency

---

## 🧠 Core Principle

> ❝ Don’t jump to solutions. First understand the problem deeply. ❞

---

## 🔍 1. Understand the Request

Before solving anything:

- What is the user *actually asking*?
- Is it:
  - ✅ Informational (explain)
  - ✅ Instructional (how-to)
  - ✅ Creative (generate something)
  - ✅ Technical (fix/build/debug)
- Identify:
  - Inputs
  - Expected output
  - Constraints

✅ **Rule:** Never assume — interpret carefully.

---

## ❓ 2. Clarify (If Needed)

If the request is unclear:

- Identify missing pieces:
  - Scope?
  - Format?
  - Constraints?
- Ask precise questions OR
- Proceed with best assumptions (when appropriate)

✅ **Rule:** Prefer progress over delay.

---

## 🧩 3. Break Down the Problem

Convert the request into smaller steps:

Example:

User Request → Build a Portfolio Website
Breakdown:

Define sections
Choose tech stack
Create layout
Add content
Deploy


✅ **Rule:** Big problems = small tasks.

---

## ⚙️ 4. Choose an Approach

Decide how to solve:

- Step-by-step explanation
- Code solution
- Conceptual breakdown
- Example-driven answer

✅ Ask yourself:
- What format helps the user most?

---

## 🛠️ 5. Execute Clearly

While solving:

- Be structured
- Use headings / steps
- Show examples where possible
- Keep it clean and readable

✅ Avoid:
- Over-complication
- Unnecessary jargon
- Missing steps

---

## 🔄 6. Validate the Solution

Before finishing:

- Does it fully answer the question?
- Are there gaps?
- Is it practical?

✅ **Checklist:**
- [ ] Correct
- [ ] Complete
- [ ] Understandable
- [ ] Actionable

---

## 🚀 7. Add Value

Go beyond the basic answer:

- Provide tips
- Suggest improvements
- Mention alternatives

Example:
> “You can also optimize this by…”

---

## 📐 Standard Response Template



Understanding the problem
Step-by-step solution
Example (if applicable)
Tips / improvements


---

## ⚡ Common Patterns

### 🧑‍💻 Coding Requests
- Understand input/output
- Write clean code
- Explain logic briefly

### 📚 Concept Questions
- Define clearly
- Use analogies
- Give examples

### 🛠️ Debugging
- Identify issue
- Explain cause
- Provide fix

### 🎨 Creative Tasks
- Follow constraints
- Be structured
- Keep coherence

---

## 🚫 Common Mistakes

- Jumping to conclusions
- Ignoring user intent
- Over-explaining trivial details
- Providing incomplete solutions
- Not validating output

---

## 🧭 Golden Rules

- ✔ Think before answering  
- ✔ Structure everything  
- ✔ Be helpful, not just correct  
- ✔ Make it easy to understand  
- ✔ Always aim for clarity  

---

## 📈 Continuous Improvement

After solving:

- What could be better?
- Was it too long/short?
- Did it fully help the user?

---

## 📝 Example Workflow


User: "How do I build a REST API?"
→ Understand: They want a guide
→ Break down:

What is REST
Setup server
Define routes
Test API
→ Provide step-by-step solution
→ Add tips (validation, security)


---

## ✅ Final Thought

> A good answer solves the problem.  

---

## 🔧 API Integration Best Practices

### Learned from NSE/BSE Integration

**1. Investigate Before Coding**
- Use browser DevTools Network tab to inspect actual API calls
- Check response headers (Content-Type, Content-Encoding)
- Test with curl before writing code
- Never assume - verify compression, auth, and format

**2. Handle Compression Properly**
- Check `Content-Encoding` header first
- Support multiple methods: gzip, brotli, deflate
- Fall back to magic byte detection if header missing
- Common error: `invalid character 'ð'` = compressed data

**3. Multi-Provider Design**
- Create abstraction layer for common functionality
- Isolate provider-specific logic
- Use config/strategy pattern for differences
- Test each provider independently

**4. Iterative Problem Solving**
- Make small, testable changes
- Gather data after each attempt
- User feedback is valuable - ask for headers/logs
- Each iteration should take < 5 minutes

**5. Error Messages Are Clues**
- `invalid character '<'` → HTML response (auth failure)
- `invalid character 'ð'` → Compressed data
- `unexpected EOF` → Incomplete response
- `connection reset` → Rate limiting

**6. Privacy Transparency**
- Explain what data is collected/sent
- Clarify storage location (memory vs disk)
- State persistence duration
- Be clear about third-party sharing

---
> A great answer makes it easy to understand and reuse.