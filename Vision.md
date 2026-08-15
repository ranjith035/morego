Vision

Build an open-source, Playwright-inspired mobile automation platform.

The framework is NOT another Appium wrapper.

It is a modern automation platform built from first principles.

Inspired by Playwright's developer experience.

Implemented independently.

Core language:

Go

Communication:

gRPC + Protocol Buffers

SDKs

TypeScript
Java
Python
Go
C#
Kotlin

License

Apache 2.0

Guiding Principles

Developer Experience first.

Performance second.

Maintainability always.

Never compromise architecture for short-term gains.

Every public API must remain stable.

Think in 10-year horizons.

Repository Structure
mobile-framework/


docs/


architecture/


vision/


design/


adr/


api/


sdk/


drivers/


recorder/


codegen/


plugins/


security/


performance/


examples/


proto/


core/


cmd/


internal/


pkg/


sdk/


drivers/


tools/


scripts/


examples/


.github/
Phase 1

Repository Foundation

Deliverables

README

Vision

Mission

Architecture

Repository structure

Roadmap

Contributing

License

Code of Conduct

Governance

Engineering Guidelines

Coding Standards

Style Guide

Naming Conventions

API Guidelines

Versioning

Release Process

Decision Framework

Architecture Decision Records

GitHub Templates

GitHub Actions

Dockerfile

Makefile

go.work

go.mod

Phase 2

Protocol

Design Protocol Buffers.

Everything communicates through gRPC.

SDKs never communicate directly with drivers.

Deliverables

proto/

session.proto

locator.proto

device.proto

driver.proto

recorder.proto

reporter.proto

plugin.proto

ai.proto

Generate SDK bindings.

Phase 3

Core Engine

Repository

core/

Implement

Execution Engine

Session Manager

Locator Engine

Wait Engine

Assertion Engine

Action Engine

Scheduler

Device Manager

Event Bus

Configuration

Logging

Metrics

Plugin Loader

Dependency Injection

Error Handling

Tracing

Lifecycle

Interfaces only first.

Implement later.

Phase 4

Driver SDK

Repository

drivers/

Design Driver interface.

Driver


Connect()


Disconnect()


Click()


Fill()


Swipe()


Screenshot()


FindElement()


LaunchApp()


TerminateApp()


InstallApp()


UninstallApp()


GetSource()


ExecuteScript()

Drivers

Appium

ADB

UiAutomator2

Espresso

XCUITest

SeeTest

BrowserStack

Sauce

LambdaTest

Phase 5

Locator Engine

Support

GetByText

GetByRole

GetByLabel

GetByPlaceholder

GetByAccessibilityId

GetByTestId

Locator()

Relative locators

Chained locators

Locator ranking.

Accessibility

↓

Test ID

↓

Role

↓

Text

↓

Resource ID

↓

XPath

Never prefer XPath.

Phase 6

Auto Wait

Everything waits automatically.

Never expose explicit waits.

Wait for

Visible

Stable

Enabled

Clickable

Retry intelligently.

No Thread.sleep.

Phase 7

Assertions

Implement

ToBeVisible

ToContainText

ToHaveValue

ToBeChecked

ToHaveCount

ToExist

ToBeEnabled

ToBeDisabled

Phase 8

Recorder

Repositories

recorder/

inspector/

Capture

Accessibility Events

Touch Events

UI Hierarchy

Window Changes

Keyboard

Clipboard

Activity Changes

Generate

Intermediate Representation

↓

Language Generator

↓

SDK Code

Support

TypeScript

Java

Python

Go

C#

Kotlin

Phase 9

Inspector

Desktop application.

Features

Hierarchy

Element Picker

Timeline

Recorder

Locator Preview

Generated Code

Screenshots

Logs

AI Suggestions

Trace Viewer

Phase 10

Reporter

HTML

Trace

Timeline

Screenshots

Video

Performance

Artifacts

JUnit

JSON

Phase 11

SDKs

Repositories

sdk-typescript

sdk-java

sdk-python

sdk-go

sdk-dotnet

sdk-kotlin

SDKs should only contain

Serialization

Deserialization

Connection

Fluent APIs

No automation logic.

Phase 12

AI

Repository

ai/

Features

Self Healing

Locator Suggestions

Visual AI

Natural Language

Failure Analysis

Code Review

Page Object Generator

AI remains optional.

Phase 13

Plugin SDK

Support

Drivers

Locators

Assertions

AI

Cloud Providers

Reporters

Authentication

Configuration

Phase 14

CLI

mobile doctor

mobile test

mobile devices

mobile recorder

mobile inspector

mobile report

mobile trace

mobile plugin

mobile doctor

Phase 15

Cloud

Device Grid

Authentication

Scheduling

Load Balancing

Queue

Reservations

REST API

Dashboard

Analytics

Architecture Rules

Always prefer

Composition

Interfaces

Dependency Injection

Modularity

Small packages

No global state

No circular dependencies

Go Standards

Idiomatic Go.

Return errors.

Use context.Context.

Small packages.

Table-driven tests.

No panic.

Prefer interfaces.

Performance Goals

10,000 parallel sessions.

<50 MB RAM per session.

Sub-millisecond protocol overhead.

Streaming logs.

Zero-copy protobufs.

Documentation

Every package requires

README

Architecture

Sequence Diagram

Examples

Tests

Benchmarks

API Docs

Testing

Unit

Integration

Contract

Performance

Compatibility

End-to-end

Long-Term Goal

Become the Playwright of mobile automation.

Compete with Appium through superior developer experience.

Not by copying Appium.

Final instruction to Codex

Act as the Chief Architect and Founding Engineer of this project. Every implementation must prioritize long-term maintainability, clean architecture, performance, extensibility, and developer experience. Never optimize for short-term convenience. Design APIs before implementations, keep SDKs thin, centralize automation intelligence in the Go core, use gRPC as the native protocol, and treat every public API as a 10-year compatibility contract. Generate production-quality code, tests, documentation, and architecture diagrams with every major feature.

I would use this as the master prompt for Codex, and then ask it to implement one phase at a time. That will produce a much higher-quality result than asking it to generate the entire framework in one shot