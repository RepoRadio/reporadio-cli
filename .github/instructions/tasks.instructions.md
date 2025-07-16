# generate-tasks.mdc

## 🌐 Purpose

This macro generates a clear, implementation-agnostic task list based on a Markdown-formatted Product Requirements Document (PRD). The goal is to help teams break down a new feature into logically grouped, actionable units of work without assuming any specific programming language, framework, or platform.

The output is intended to:

* Facilitate planning, backlog grooming, and handoff to developers.
* Support a junior engineer in understanding what needs to be done.
* Enable traceability to the original PRD.

This macro assumes the PRD was generated via `create-prd.mdc`.

---

## ✅ What Makes a Good Task

Each task should:

* Start with a verb (e.g., "Define", "Draft", "Review", "Clarify")
* Describe an action with clear intent
* Be understandable without technical context
* Take no more than a day to complete
* Align with a functional requirement or user story
* Avoid bundling multiple actions into one task

---

## 🔄 Workflow

1. **Ingest PRD**

   * Parse the PRD Markdown and analyze all sections, especially:

     * Overview
     * Goals
     * Functional Requirements
     * User Stories
     * Non-Goals

2. **Generate Parent Tasks**

   * Identify major functional components or workflows
   * Output them as top-level tasks (do not yet break them down)
   * Pause here unless the `--no-pause` flag is provided

3. **(Optional) Pause for Confirmation**

   * Wait for the user to say "Go" before proceeding to sub-task generation
   * This supports collaborative refinement of the parent task structure

4. **Generate Sub-Tasks**

   * For each parent task, decompose into specific, atomic actions
   * Annotate sub-tasks with references to PRD sections for traceability (e.g., `[From: Functional Requirement 2]`)

5. **(Optional) Files to Be Created**

   * If the PRD suggests or implies deliverables (e.g., docs, diagrams), list them and describe their purpose
   * Do not reference any specific code files, directories, or formats

6. **Output Structure**
   Format the final output as follows:

```md
# Task List: [Feature Name]

## Task 1
- [ ] A  <!-- From: User Story 1 -->
- [ ] B  <!-- From: Functional Requirement 3 -->

## Task 2
- [ ] C
- [ ] D

## Optional: Files to Create
- `flowchart.md`: captures key user flows from the PRD
- `acceptance-checklist.md`: tracks criteria from user stories
```

7. **Save Output**

   * Write to `/tasks/tasks-[feature-name].md`

8. **Prompt for Review**

   * Close by prompting the user to review the task list:

     > "Before continuing, please review this task list for missing assumptions, undefined terms, or ambiguous phrasing. Let me know if you'd like to revise any part."

---

## 🔒 Constraints

* ❌ Do NOT reference specific technologies, programming languages, or tools
* ❌ Do NOT output file extensions like `.ts`, `.jsx`, `.py`, etc.
* ✅ Use plain English throughout
* ✅ Assume the audience includes junior developers
* ✅ Be unambiguous, neutral, and implementation-agnostic

---

## 🔍 Clarification Phase

Before generating tasks, ask clarifying questions to ensure accuracy and completeness:

1. **Sequential Questions**
   * Ask one question at a time and wait for the user's response
   * Focus on areas where the PRD may be ambiguous or incomplete
   * Target questions that will impact task structure or scope

2. **Question Categories**
   * **Scope boundaries**: What's explicitly included vs. excluded?
   * **User personas**: Who are the primary users and what are their contexts?
   * **Integration points**: How does this feature connect to existing systems?
   * **Success criteria**: What specific outcomes define completion?
   * **Assumptions**: What dependencies or prerequisites exist?

3. **Question Format**
   * Be specific and actionable
   * Reference relevant PRD sections when asking
   * Avoid yes/no questions when possible

4. **Completion Criteria**
   * Continue asking questions until you have sufficient clarity
   * Confirm understanding before proceeding to task generation
   * Ask: "Do you have any other clarifications before I generate the task list?"

---
