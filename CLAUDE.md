## Project Overview

- `General`:
  - `Name`: CyberDiner
  - `Description`: An AI-powered food delivery system designed for high performance and personalized user experiences.
  - `Scope`: Graduation Thesis Project
  - `Version`: 1.0.0

## Architecture & Directory Structure

- `CLAUDE.md`: Source of truth for AI agents.
- `.claude`: AI agents configuration files.
  - `settings.json`: AI agents general settings.
  - `hooks`: Guardrail that checks before and after agents response.
  - `skills`: AI agents skills.
- `UIUX_design`: Static design template of the project.
- `docs`: Documents of the project.
  - `architecture.md`: Architecture of the project.
- `src`: Source code of the project.
  - `AI`: FastAPI (Python) Service.
  - `server`: The Go/Gin Backend.
    - `cmd`: Contains the entry points
      - `api`:
        - `routes.go`: API endpoint declarations
        - `main.go`: Entry point: initializes and starts the app
    - `internal`: Private code (cannot be imported by other projects)
      - `app`: Application-wide logic (wire-up, middleware)
      - `handler`: HTTP/Transport layer (routing and controllers)
      - `service`: Business logic (domain rules)
      - `store`: Data access layer (SQL, NoSQL, GORM)
      - `middleware`: The "Filters": Auth, logging, and recovery
      - `model`: Domain entities/structs
      - `locales`: I18next locales files and configurationss
    - `pkg`: Public code (safe for other projects to import)
      - `util`: Small, reusable helper functions
      - `validator`: Custom data validation logic
      - `logger`: Custom logging wrappers
    - `api`: API definitions (OpenAPI/Swagger specs)
    - `configs`: Config files (env, yaml, json)
    - `go.mod`: Dependency management
    - `go.sum`: Dependency checksums
    - `.gitignore`: Git ignored files
    - `.env`: Environment variables
    - `.env.example`: Example environment variables
  - `client`: The React Frontend.
    - `src`:
      - `components`: Reusable UI components
      - `pages`: Page-level components
      - `assets`: Static assets
      - `api`: API client
      - `services`: API client logic
      - `store`: State management
      - `utils`: Utility functions
      - `locales`: I18next locales files and configurations
      - `App.jsx`: Main application component
      - `main.jsx`: Entry point
    - `public/`: Public assets
    - `.gitignore`: Git ignored files
    - `package.json`: Project dependencies and scripts
    - `package-lock.json`: Dependency lock file
    - `.env`: Environment variables
    - `.env.example`: Example environment variables
    - `vite.config.js`: Vite configuration
    - `tailwind.config.js`: Tailwind configuration
    - `index.html`: Main HTML entry

## Tech Stack & Patterns

- `Framework/Language Versions`:
  - `Server`: Gin / Go 1.21+
  - `Client`: React 18+ / Node.js
  - `AI`: FastAPI / Python 3.10+
- `Database`: PostgreSQL (use `gorm` for Go).
- `Communication`:
  - `Client-Server`: REST APIs (JSON).
  - `Server-AI`: REST APIs (JSON).
  - `Server-Database`: GORM (ORM).

## Languages

- English/EN : default language
- Vietnamese/VI : secondary language
