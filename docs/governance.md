# Project Governance

This document describes the governance structure and decision-making processes of the Mobile Automation Platform project.

---

## 1. Roles and Responsibilities

Our community is structured around three primary roles:

### 1.1. Contributors
Contributors are individuals who participate in the community by writing code, files, issue comments, triage, or documentation. Contributions can range from bug reports to large feature implementations.

### 1.2. Maintainers
Maintainers are trusted contributors who have demonstrated consistent commitment and technical alignment with the project. Maintainers have write/commit access to the primary repositories.
*   **Responsibilities:** Triaging issues, reviewing pull requests, managing release pipelines, and maintaining core library modules.
*   **Earning Maintainership:** A contributor can be nominated for maintainership by any existing maintainer after demonstrating high-quality work, a constructive attitude, and deep understanding of the architecture across several pull requests. A unanimous vote of existing maintainers is required for approval.

### 1.3. Steering Committee
The Steering Committee consists of the project founders and key lead architects.
*   **Responsibilities:** Setting long-term project vision, resolving technical stalemates, and managing domain registrations/assets.

---

## 2. Decision-Making Process

We strive for consensus-based decision-making. However, if a consensus cannot be reached, the following voting process is used:

### 2.1. Technical Architecture Changes
All major technical design changes must go through the **Architecture Decision Record (ADR)** process.
1.  **Drafting:** A proposal is created under `docs/adr/`.
2.  **Discussion:** The community discusses the proposal inside PR reviews.
3.  **Resolution:** 
    *   For minor changes: Approval by at least two maintainers.
    *   For major structural changes (e.g., protocol changes, core interfaces): A voting period of 7 days is opened. A minimum of 3 maintainers must vote. Decisions require a majority vote (50%+1). If there is a tie, the Steering Committee holds the deciding vote.

### 2.2. Standard Pull Requests
Standard bug fixes and minor improvements do not require voting or formal ADRs. They simply require approval from at least one maintainer (who is not the author of the PR) before being merged into `develop`.
