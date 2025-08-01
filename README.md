# 🔥 Agni - A Modular Go Code Analyzer

Agni is a blazing-fast, modular static code analysis tool built in Go. Designed to enhance code hygiene and readability, Agni scans your Go projects to flag patterns that affect quality, maintainability, and scalability.

With a lightweight CLI interface, Agni integrates seamlessly into your developer workflow — offering clarity without complexity.

Whether you're working on large enterprise codebases or tight microservices, **Agni keeps your code sharp.**

---

## 🚀 Key Features

- ✅ Detect unused function parameters  
- 💬 Identify unused constants and internal log messages  
- 📁 Detect dead code via automatic [`deadcode`](https://pkg.go.dev/golang.org/x/tools/cmd/deadcode) integration  
- 🔍 Spot unused keys in `Messages`, `FailMessages`, etc.  
- 🧼 Modular design – plug in more detectors easily  
- 🚀 Detect the undefined keys used in messageMap in through out the project. 
- 🧼 Detects capital variable names, function parameters and returning parameters.
- 📁 Detects the packages that are used in the code base but actually are deprecated by golang or organization standards. 

> ⚙️ More powerful static checks are coming in future versions!

---

## 📦 Installation

## Install Agni using `go install`:

go install github.com/Aadi-IRON/agni/cmd/agni@latest

 RUN -> agni check

## If already installed and want latest version
-> go clean -modcache
-> go install github.com/Aadi-IRON/agni/cmd/agni@latest

And then, agni check 
