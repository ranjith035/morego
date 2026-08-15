# Project Vision & Mission

Our goal is to build an open-source, Playwright-inspired mobile automation platform that solves the complexity and flakiness of modern mobile app testing.

---

## The Vision

The mobile automation ecosystem has been stagnant. The dominant tool, Appium, is slow, complex to set up, and prone to flaky tests due to its legacy HTTP architecture and lack of built-in orchestration. Alternatives like Maestro are fast and reliable but limit developers to static YAML files, blocking programmatic control, dynamic test flow integration, and complex testing patterns.

We are building a modern mobile automation platform from first principles. It is **NOT** another Appium wrapper. 

Key principles of our vision:
* **Developer Experience First:** High-quality SDKs in standard languages, easy setup, automatic waiting, and excellent debugging tools (Inspector, Recorder, Trace Viewer).
* **Performance Second:** Sub-millisecond communication via gRPC and Protocol Buffers, with a fast, lightweight Go-based Core execution engine that uses less than 50MB RAM per session.
* **Maintainability Always:** Clean, modular interfaces, thin SDK clients, and stable APIs that are treated as 10-year compatibility contracts.

---

## Architecture at a Glance

Instead of communicating over legacy HTTP WebDriver protocols, our platform uses a high-performance, centralized Go Core Engine. All SDK client languages (TypeScript, Java, Python, Go, C#, Kotlin) serve as thin serialization layers that communicate with this core via gRPC streaming.

```
       +------------------+
       |   Client SDK     |
       | (TS, Python, Go) |
       +--------+---------+
                |
                | (gRPC + Protocol Buffers)
                v
       +------------------+
       |   Core Engine    |  <--- Wait Engine, Locator Engine, Scheduler
       |     (Go)         |
       +--------+---------+
                |
                | (Direct Drivers / ADB / XCUITest)
                v
       +------------------+
       |  Mobile Device   |
       +------------------+
```

---

## Core Milestones & Roadmap

Our development is structured into clear phases to ensure high engineering quality:

1. **Phase 1: Repository Foundation** - Directories, documentation, coding guidelines, CI/CD pipelines, and build targets.
2. **Phase 2: Protocol Definition** - Designing Protocol Buffers for session, locator, device, driver, and recording.
3. **Phase 3: Core Engine** - Building the execution engine, scheduler, device manager, and event bus in Go.
4. **Phase 4: Driver SDK** - Defining the driver interface and building wrappers for UiAutomator2, Espresso, XCUITest, ADB, and cloud grids.
5. **Phase 5: Locator Engine** - Implementing accessibility-first locators (`GetByRole`, `GetByLabel`) with auto-ranking.
6. **Phase 6: Auto Wait** - Implementing background wait loops to ensure elements are ready before actions are performed.
7. **Phase 7: Assertions** - Providing web-like expectations (`ToBeVisible`, `ToContainText`, `ToBeChecked`).
8. **Phase 8: Recorder** - Capturing gestures and generating test code in multiple SDK languages.
9. **Phase 9: Inspector** - Building a desktop companion app for inspection, debugging, and live locator preview.
10. **Phase 10: HTML/Trace Reporter** - Recording visual timelines, network requests, screenshots, and logs.
11. **Phase 11: Multi-language SDKs** - Autogenerating and publishing the thin clients.
12. **Phase 12: AI Capabilities** - Self-healing tests, visual validation, and natural language query processing.
13. **Phase 13: Plugin SDK** - Enabling custom locators, drivers, and reporters.
14. **Phase 14: CLI Tools** - Diagnostic doctor tools, test runners, and device managers.
15. **Phase 15: Cloud Device Grid** - Scalable orchestration and reservation service.
