> # Things i have learned over days

1. centeralize everything
   - Use application
     - application structure to hold data and `app := &application{}` to mutate data
     - app acts as service and handler provider for the web application
   - Use logger
     - using slog to create logger for different level or error (they )
2. dont centeralize layers (use internal)
   - Database layer (snippets)
     - defer connection pool close for query
     - dont assume all rows are read error and end of row.next() both returns false
   - Custom errors (mainly for error comparision between layers keeping layer clean from un-necessary dependency)
3. use templates without pointer to hold data
   - template cache
   - CSRF
   - common data
4. use handlers and helpers
   - handlers
     - they handel purely routes for me
     - {$} to stop from wild character routes or dynamic routes
     - avoid routes overlaping especially dynamic routes.
     - and can make two handlerfunc for same route with different rest api call `GET POST...`
   - helpers
     - the render parse templete to templete structure using glob and cache them to template data
       - use custom buffer to execute the template first and `buff.WriteTo(http.requestWriter)` in-order to handel run time errors
     - custom serverError, clientError and notFound helper for http errors
