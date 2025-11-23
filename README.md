# Wound IQ API

Production-ready REST API in Go (Gin) for wound care data backed by PostgreSQL `wound_iq`.

## Features
- CRUD for patients, clinicians, assessments
- Full-assessment creation via stored function `add_full_assessment(...)` (mode A)
- Reports using DB stored functions: get_assessment_full, get_patient_wound_history, get_all_patients, get_all_assessments
- Pagination, filters (assessments by patient_id, clinician_id, date range)
- Graceful shutdown, structured logging (zerolog), request validation, unit tests
- OpenAPI spec included

## Quickstart (local)

1. Ensure Go 1.22+ is installed.
2. Ensure PostgreSQL is running and `wound_iq` DB exists with schema, sample data, and functions loaded.
3. Copy `.env.example` to `.env` and update `DB_DSN`.
4. Run:
   ```
   make run
   ```
5. API will be available at `http://localhost:8080/v1`.

## Example: Create full assessment (calls add_full_assessment)
```bash
curl -X POST http://localhost:8080/v1/assessments \
  -H "Content-Type: application/json" \
  -d @full_assessment.json
```

## Tests
```
make test
```

## Notes
This repository implements Mode A from the specification: `POST /v1/assessments` always calls `add_full_assessment(...)`.
