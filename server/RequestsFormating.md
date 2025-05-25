The possible requests that frontend need to handle:
In header:
- ANY - session-token as "Session" parameter
- GET - content-type as "Content-Type" parameter
In URL:
- GET(json)     - page as "page" parameter
- GET(task)     - task-id as "id" parameter
- DELETE(task)  - task-id as "id" parameter
- ANY           - language as "lang" parameter

The responses gurantee:
- Status code
- Content-type
- Session (if authorized)