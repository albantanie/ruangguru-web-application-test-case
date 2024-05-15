# Exercise Pemrograman Backend Lanjutan

## Build Simple Blog with Gin Framework 2

### Description

Tugas ini adalah untuk membuat REST API blog sederhana dengan menggunakan framework Gin pada bahasa pemrograman Go. REST API blog ini memiliki fitur-fitur seperti menampilkan semua postingan, menampilkan satu postingan berdasarkan id, dan menambahkan postingan baru dengan menggunakan _authenctication middleware_. Aplikasi ini akan memiliki 2 endpoint yaitu:

- `GET /posts` :
  - Endpoint ini akan mengembalikan seluruh data postingan.
  - Jika menggunakan Query Param `?id={id}` ini akan mengembalikan data postingan berdasarkan id
- `POST /posts` : Endpoint ini akan menambahkan data postingan baru.

Untuk setiap data postingan, terdapat atribut berikut:

- `ID` : Tipe data `int`, id dari postingan.
- `Title` : Tipe data `string`, judul dari postingan.
- `Content` : Tipe data string`, isi dari postingan.
- `CreatedAt` : Tipe data `time.Time`, waktu dibuatnya postingan.
- `UpdatedAt` : Tipe data `time.Time`, waktu terakhir postingan diupdate.

Sesuai dengan struct berikut:

```go
// struct
type Post struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// default data
var Posts = []Post{
  {ID: 1, Title: "Judul Postingan Pertama", Content: "Ini adalah postingan pertama di blog ini.", CreatedAt: time.Now(), UpdatedAt: time.Now()},
  {ID: 2, Title: "Judul Postingan Kedua", Content: "Ini adalah postingan kedua di blog ini.", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}
```

Semua endpoint di atas harus menggunakan basic authentication dengan data user yang telah di sediakan yaitu:

```go
// struct
type User struct {
  Username string `json:"username"`
  Password string `json:"password"`
}

// default data
var users = []User{
  {Username: "user1", Password: "pass1"},
  {Username: "user2", Password: "pass2"},
}
```

Sebaiknya data authentication bisa di encode dengan [base64](https://www.base64encode.org/) pada saat melakukan request, contoh:

```http
Header: "Authorization", "Basic dXNlcjE6cGFzczE="
```

### Constraints

Tugas ini memiliki beberapa batasan dan hal yang harus diperhatikan, yaitu:

- Format data yang dikembalikan harus sesuai dengan format **JSON**.
- Pada saat menampilkan satu postingan berdasarkan id, jika id yang dimasukkan tidak ditemukan pada daftar postingan maka program harus memberikan response dengan status HTTP **404**.
- Pada saat menambahkan postingan baru, program harus mengecek apakah request body yang diberikan sudah sesuai dengan format yang diharapkan. Jika tidak sesuai, program harus memberikan response dengan status HTTP **400**.
- Semua endpoint hanya dapat di akses jika menyertakan _basic authentication_ pada http header saat melakukan request menggunakan middleware. Jika tidak, maka akan memberikan response dengan status HTTP **401**

Kerjakan aplikasi sesuai dengan ketentuan berikut:

- Semua endpoint harus menggunakan basic _authentication middleware_.
  - Jika tidak, maka akan mengembalikan response dengan status `401`
- Pada endpoint `"/posts"` dengan method **GET**, API akan mengembalikan seluruh postingan yang tersimpan pada variabel `Posts` dalam bentuk JSON.
  - Jika menggunakan Query Param `?id={id}` API akan mencari postingan berdasarkan ID yang diminta dan mengembalikan data postingan tersebut dalam bentuk JSON
    - Jika tidak ditemukan postingan dengan ID yang diminta, maka API akan mengembalikan response dengan status `404` dan pesan `"Postingan tidak ditemukan"`.
    - Jika ID yang diberikan bukan angka, maka API akan mengembalikan response dengan status `400` dan pesan `"ID harus berupa angka"`.
- Pada endpoint `"/posts"` dengan method `POST`, API akan menambahkan sebuah postingan baru dengan data yang dikirimkan dalam request body.
  - Jika request body tidak valid, API akan mengembalikan response dengan status `400` dan pesan `"Invalid request body"`.
  - Jika semua data request body sudah valid, maka postingan baru akan disimpan pada variabel Posts dan API akan mengembalikan response dengan status `201` dan pesan `"Postingan berhasil ditambahkan"`

### Test Case Examples

#### Test Case 1

**Input**:

```http
GET /posts
```

**Expected Output / Behavior**:

```http
HTTP status code: 401 Unauthorized
```

**Explanation**:

> Ketika melakukan request `GET /posts` tanpa melampirkan `Header: "Authorization", "Basic dXNlcjE6cGFzczE="` maka akan mengembalikan status code `401 Unauthorized`.

#### Test Case 2

**Input**:

```http
GET /posts
Header: "Authorization", "Basic dXNlcjE6cGFzczE="
```

**Expected Output / Behavior**:

```http
HTTP status code: 200 OK
Response body: {
  "posts": [
    {
      "id": 1,
      "title": "Judul Postingan Pertama",
      "content": "Ini adalah postingan pertama di blog ini.",
      "created_at": "2023-04-05T12:00:00Z",
      "updated_at": "2023-04-05T12:00:00Z"
    },
    {
      "id": 2,
      "title": "Judul Postingan Kedua",
      "content": "Ini adalah postingan kedua di blog ini.",
      "created_at": "2023-04-05T12:00:00Z",
      "updated_at": "2023-04-05T12:00:00Z"
    }
  ]
}
```

**Explanation**:

> Ketika melakukan request `GET /posts`, dengan melampirkan `Header: "Authorization", "Basic dXNlcjE6cGFzczE="` program akan menampilkan semua postingan yang ada pada server dengan status code `200 OK`.

#### Test Case 3

**Input**:

```http
GET /posts?id=1
Header: "Authorization", "Basic dXNlcjE6cGFzczE="
```

**Expected Output / Behavior**:

```http
HTTP status code: 200 OK
Response body: {
  "post": {
    "id": 1,
    "title": "Judul Postingan Pertama",
    "content": "Ini adalah postingan pertama di blog ini.",
    "created_at": "2023-04-05T12:00:00Z",
    "updated_at": "2023-04-05T12:00:00Z"
  }
}
```

**Explanation**:

> Ketika melakukan request `GET /posts/1`, program akan mencari postingan dengan `id=1` pada daftar postingan. Jika postingan dengan `id` tersebut ditemukan, program akan menampilkan detail postingan tersebut dengan status code `200 OK`.

#### Test Case 4

**Input**:

```http
GET /posts?id=999
Header: "Authorization", "Basic dXNlcjE6cGFzczE="
```

**Expected Output / Behavior**:

```http
HTTP status code: 404 Not Found
Response body: {
  "error": "Postingan tidak ditemukan"
}
```

**Explanation**:

> Ketika melakukan request `GET /posts/1`, program akan mencari postingan dengan `id=1` pada daftar postingan. Jika tidak ditemukan, program akan memberikan response dengan status HTTP `404 Not Found`.

#### Test Case 5

**Input**:

```http
POST /posts
Header: "Authorization", "Basic dXNlcjE6cGFzczE="
Content-Type: application/json

{
  "title": "Judul Postingan Baru",
  "content": "Ini adalah postingan baru di blog ini."
}
```

**Expected Output / Behavior**:

```http
HTTP status code: 201 Created
Response body: {
  "message": "Postingan berhasil ditambahkan",
  "post": {
    "id": 3,
    "title": "Judul Postingan Baru",
    "content": "Ini adalah postingan baru di blog ini.",
    "created_at": "2023-04-05T12:00:00Z",
    "updated_at": "2023-04-05T12:00:00Z"
  }
}
```

**Explanation**:

> Ketika melakukan request `POST /posts` dengan data JSON yang valid, program akan menambahkan postingan baru ke dalam daftar postingan dengan ID yang baru dan menampilkan response dengan status HTTP `201 Created`.
