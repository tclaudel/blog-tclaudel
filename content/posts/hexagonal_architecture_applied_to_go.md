---
title: "Hexagonal architecture with Go"
date: 2023-09-20T21:41:18+02:00
draft: true
---

## Hexagonal architecture

Hexagonal architecture was invented by [Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/) 
and published in 2005.
Is a software architecture that aims to create loosely coupled application components with isolation between 
business logic and technical details. 

The hexagonal architecture is also known as the ports and adapters architecture,

### Why ?

If offers a number of advantages, such as :
- **Independence from framework**: your application is no longer directly dependent on external libraries.
- **Testability**: writing tests is greatly facilitated by the decoupling of dependencies.
- **Flexibility ans scalability**: the application is more flexible and scalable because it is not tied to a specific 
framework. It's easier to change the framework or add new features.

### Key Components

To achieve this, we will define 3 components:
- **Domain**: the business logic of the application, your entities, your business rules, etc.
- **Ports**: the interfaces that define how the domain interacts with the outside world.
- **Adapters**: the implementations of the ports.

## Domain

The domain is the core of the application. It contains the business logic, the entities, the business rules, etc.
It must not depend on anything. It must be completely independent of the outside world.  

**⚠️ NO INFRASTRUCTURE CODE IN THE DOMAIN. ⚠️**

![Domain](./images/content/hexagonal_architecture_with_go/domain.png)
