# 📘 Agent Skillbook

A clear, structured guide for solving problems effectively.

---

## 🎯 Core Principle

> **Think first, then act. Break big problems into small steps.**

---

## 📋 Problem-Solving Process

### 1. Understand the Request
- What is the user asking for?
- What's the expected output?
- Are there any constraints?

### 2. Break Down the Task
Convert complex requests into simple steps:
```
Example: "Centralize API URLs"
→ Create constants file
→ Update main.go
→ Update other files
→ Test changes
```

### 3. Execute Step-by-Step
- Work through one step at a time
- Wait for confirmation before proceeding
- Document as you go

### 4. Validate & Test
- Does it work correctly?
- Are there any gaps?
- Test the changes

---

## 📁 Task Management System

### When to Use
Use for tasks with:
- Multiple steps (>3 subtasks)
- Changes across multiple files
- Need for progress tracking

### Process

**1. Create Task Documentation**
- File: `AGENT/TASK_NAME.md`
- Include:
  - Overview
  - Changes made (before/after)
  - Files modified
  - Test results
  - Rollback instructions

**2. Update Task Registry**
- File: `AGENT/task_registry.json`
- Add task entry with subtasks
- Track status and files

**3. Track Progress**
- Use `update_todo_list` tool
- Mark subtasks as completed
- Update documentation

### File Structure
```
AGENT/
├── README.md              # System guide
├── task_registry.json     # Task tracking
└── TASK_NAME.md          # Task documentation
```

---

## 🔧 Best Practices

### API Integration
- Inspect API calls in browser DevTools first
- Check response headers (Content-Type, Content-Encoding)
- Test with curl before coding
- Handle compression (gzip, brotli, deflate)

### Error Handling
- `invalid character '<'` → HTML response (auth failure)
- `invalid character 'ð'` → Compressed data
- `unexpected EOF` → Incomplete response
- `connection reset` → Rate limiting

### Code Changes
- Make small, testable changes
- Gather feedback after each step
- Document what was changed and why
- Provide rollback instructions

---

## ✅ Golden Rules

1. **Think before acting** - Understand the problem first
2. **Break it down** - Small steps are easier
3. **Document everything** - For future reference
4. **Test thoroughly** - Verify it works
5. **Be clear** - Simple explanations are best

---

## 📝 Example Workflow

```
Task: "Centralize API URLs"

1. Break down:
   ✓ Create constants.go
   ✓ Update main.go
   ✓ Update bse_shareholding.go
   ✓ Update symbol_mapper
   ✓ Test all URLs

2. Document:
   - Create AGENT/DEBUG_URL_CONSTANTS.md
   - Update AGENT/task_registry.json
   - Track with update_todo_list

3. Test:
   - Build project
   - Run tests
   - Verify URLs work

4. Complete:
   - Mark all subtasks done
   - Document results
   - Provide summary
```

---

**Remember:** Clear thinking → Clear code → Clear results