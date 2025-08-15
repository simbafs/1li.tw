# GeoIP Information Processing Plan

This document outlines the plan to integrate GeoIP information for each click by processing IP addresses in the background.

## 1. Overview

The goal is to enrich click analytics data with geographic information derived from the user's IP address. To achieve this, we will use the `ip-api.com` batch API. To handle the API's rate limit (15 requests per minute) and to ensure that the click recording process remains fast, we will implement a background processing service.

This service will periodically fetch unprocessed clicks from the database, query the GeoIP API, and update the records with the retrieved information.

## 2. Implementation Steps

The implementation is divided into three main phases:

### Phase 1: Database and Domain Layer Modifications

This phase focuses on preparing the database schema and the application's domain layer to store the new GeoIP information.

1.  **Update Database Schema (`sql/schema.sql`)**:
    *   The `url_clicks` table will be modified to include the following new nullable fields:
        *   `country` (TEXT)
        *   `region_name` (TEXT)
        *   `city` (TEXT)
        *   `lat` (REAL)
        *   `lon` (REAL)
        *   `isp` (TEXT)
        *   `as_info` (TEXT, for "AS number and organization")
    *   A new boolean field `is_processed` will be added with a default value of `false`. This flag will track whether a click has been processed by the GeoIP service.

2.  **Update SQL Queries (`sql/queries/url_clicks.sql`)**:
    *   A new query, `GetUnprocessedClicks`, will be created to fetch a batch of records (e.g., up to 100) where `is_processed` is `false`.
    *   A new query, `UpdateClickGeoInfo`, will be created to update a single `url_clicks` record with the new GeoIP data and set `is_processed` to `true`.

3.  **Regenerate `sqlc` Code**:
    *   Run `sqlc generate` to update the Go code based on the schema and query changes.

4.  **Update Domain Layer (`domain/click.go`)**:
    *   The `URLClick` struct will be updated to include the new fields: `Country`, `RegionName`, `City`, `Lat`, `Lon`, `ISP`, `ASInfo`, and `IsProcessed`.

### Phase 2: Background Processor Service

This phase involves creating a dedicated background service to handle the GeoIP data fetching and processing.

1.  **Create `GeoIPProcessor` (`infrastructure/processor/geoip_processor.go`)**:
    *   A new service, `GeoIPProcessor`, will be created. It will depend on the `ClickRepository`.
    *   It will have a `Start()` method that launches a new goroutine for the background processing loop.

2.  **Implement the Processing Loop**:
    *   The `Start()` method will initialize a `time.Ticker` set to a 4-second interval (to stay within the 15 requests/minute rate limit).
    *   The loop will wait for the ticker to fire.
    *   On each tick, the processor will:
        1.  Call `clickRepo.GetUnprocessedClicks()` to get a batch of up to 100 unprocessed clicks.
        2.  If there are no clicks to process, the loop continues to the next tick.
        3.  Extract the IP addresses from the clicks and prepare a request body for the `ip-api.com` batch endpoint.
        4.  Make a `POST` request to `http://ip-api.com/batch?fields=60121&lang=en`.
        5.  Parse the JSON response.
        6.  For each successfully processed IP in the response, call `clickRepo.UpdateClickGeoInfo()` with the corresponding `click.id` and the retrieved geo data.

### Phase 3: Integration and Refactoring

This final phase integrates the new background service into the application and refactors existing code to simplify the click recording process.

1.  **Refactor `URLUseCase` (`application/url_usecase.go`)**:
    *   The `GeoIPService` dependency will be removed from the `URLUseCase`, as its role is now fulfilled by the `GeoIPProcessor`.
    *   The `RecordClick` method will be simplified. It will no longer be responsible for fetching GeoIP data. Its only job is to create a `URLClick` entity with the available information (IP address, User-Agent) and save it to the database with `is_processed` as `false`.

2.  **Update `main.go`**:
    *   In the `main` function, an instance of the `GeoIPProcessor` will be created.
    *   The `geoIPProcessor.Start()` method will be called to begin the background processing as the application starts.

3.  **Update Service Initialization (`presentation/gin/router.go` and `main.go`)**:
    *   The `NewURLUseCase` constructor calls will be updated to remove the `GeoIPService` dependency.

This plan ensures a robust, scalable, and efficient system for enriching click data without impacting the performance of the core URL redirection functionality.
