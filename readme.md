# How to use
## Run
Modify `env.example` to match your environment, and rename it to `.env`. 
After that import `news.sql` to mysql to generate tables, don't forget to create the database first

to run go use :
```cmd
go run main.go
```

## API
### Create News
`[POST] http://localhost:8000/api/v1/news/`
```json
{
  "title": "judul bagus", // string
  "content": "contentnya biasa ternyata biasa aja nih", // string
  "status": "publish", // draft, publish, deleted
  "tags": ["ecef5cd5-72dc-42cb-a7e1-ae5578317228"], // tag id, from table tag
  "topic": "topic 1" // string
}
```
### Get All News
`[GET] http://localhost:8000/api/v1/news/` (show all news with status publish)

### Get News By Status
`[GET] http://localhost:8000/api/v1/news/status/:status` 
status only accept `draft`, `publish` and `deleted`

### Get News By Topic
`[GET] http://localhost:8000/api/v1/news/topic/:topic` show news with exact topic value on database

### Get News By Slug
`[GET] http://localhost:8000/api/v1/news/:slug` show news with exact slug value on database

### Update News
`[PATCH] http://localhost:8000/api/v1/news/:id`
```json
{
    "title": "judul bagus", // string
    "content": "contentnya sudah terupdate", // string
    "status": "publish", // draft, publish, deleted
    "tags": ["ecef5cd5-72dc-42cb-a7e1-ae5578317228", "de642a08-c553-479d-80d1-311e6dc687f8"],// tag id, from table tag
    "topic": "topic 3" // string
}
```
we use `PATCH`, so we need to send all field

### Delete News
`[DELETE] http://localhost:8000/api/v1/news/:id`

### Create Tag
`[POST] http://localhost:8000/api/v1/tag/`
```json
{
  "name": "bagus sangat"
}
```

### Get All Tag
`[GET] http://localhost:8000/api/v1/tag/`

### Update Tag
`[PUT] http://localhost:8000/api/v1/tag/:id`
```json
{
  "name": "bagus"
}
```

### Delete Tag
`[DELETE] http://localhost:8000/api/v1/tag/:id`