# Design Philosophy

This document outlines the key design decisions and philosophy behind the OpenHS project.

## Why Golang?

### 1. Memory Safety and Efficiency

Card game engines typically make extensive use of pointers (such as in Trigger, Event systems, etc.). Golang provides an excellent balance between:
- Memory safety
- Development and debugging efficiency
- Runtime execution efficiency

This balance is crucial for building a robust and maintainable card game engine.

### 2. Superior LLM Code Generation

Golang has significant advantages for AI-assisted development:
- As a popular language in recent years, there is an abundance of Golang code on GitHub for LLMs to learn from
- Compared to C++ or Rust, Golang has simpler syntax and a more comprehensive standard library
- Implementations of similar functionalities in Golang tend to be more uniform
- These factors result in higher accuracy and better quality when generating Golang code with LLMs, which is beneficial for projects that rely heavily on AI-generated code

### 3. Static Compilation Benefits

As a statically compiled language, Golang offers:
- Compile-time checks that can provide feedback to LLMs
- Reduction in "hallucinations" or incorrect code generation
- Early detection of potential issues before runtime

