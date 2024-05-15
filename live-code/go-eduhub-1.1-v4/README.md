# Web Application

## Live Coding - Go EduHub 1

### Implementation technique

Siswa akan melaksanakan sesi live code di 15 menit terakhir dari sesi mentoring dan di awasi secara langsung oleh Mentor. Dengan penjelasan sebagai berikut:

- **Durasi**: 15 menit pengerjaan
- **Submit**: Maximum 10 menit setelah sesi mentoring menggunakan `grader-cli submit`
- **Obligation**: Wajib melakukan _share screen_ di breakout room yang akan dibuatkan oleh Mentor pada saat mengerjakan Live Coding.

### Description

**Go Eduhub** adalah sebuah aplikasi yang dirancang untuk membantu pengelolaan dan manajemen data siswa dan kursus menggunakan bahasa pemrograman Go. Aplikasi ini memungkinkan pengguna untuk melakukan berbagai operasi seperti menambah dan menampilkan data siswa juga menambah kursus yang terkait dengan siswa tersebut.

Dalam live-code ini, kita akan mengimplementasikan API menggunakan _Golang web framework Gin_ untuk mengelola data _student_ dan _course_. API harus mengizinkan client untuk:

- Menampilkan daftar semua siswa
- Menampilkan satu data siswa berdasarkan  ID
- Menambahkan siswa baru
- Menambahkan kursus ke daftar kursus siswa

Disini sudah ditentukan endpoint untuk setiap operasi untuk mengimplementasikan logika yang diperlukan dari setiap operasi menggunakan repository student dan course di file `main.go` dengan endpoint group sebagai berikut:

```go
student := gin.Group("/student")
{
    student.POST("/add", apiHandler.StudentAPIHandler.AddStudent)
    student.GET("/gets", apiHandler.StudentAPIHandler.GetStudents)
    student.GET("/get/:id", apiHandler.StudentAPIHandler.GetStudentByID)
}

course := gin.Group("/course")
{
    course.POST("/add", apiHandler.CourseAPIHandler.AddCourse)
}
```

### Constraints

Pada live code ini, kamu harus melengkapi fungsi dari repository dan handler api `student` dan `course` sebagai berikut:

📁 **repository**

Ini adalah fungsi yang berinteraksi dengan inmemory database:

- `repository/student.go`
  - `Store`: Function ini menyimpan data mahasiswa yang diberikan sebagai argumen ke dalam database. 
    - Jika proses tersebut berhasil, function akan mengembalikan `nil` sebagai `error`.
    - Namun jika terjadi `error` pada proses tersebut, function akan mengembalikan `error` yang terjadi.

- `repository/course.go`
  - `Store`: Function ini akan menyimpan data kursus yang diberikan sebagai argumen ke dalam database.
    - Jika proses tersebut berhasil, function akan mengembalikan `nil` sebagai `error`.
    - Namun jika terjadi `error` pada proses tersebut, function akan mengembalikan `error` yang terjadi.

📁 **api**

- `api/student.go`
  - `GetStudentByID`: fungsi ini akan mengambil data mahasiswa berdasarkan `id` yang diberikan sebagai argumen. Pertama-tama, fungsi akan mengambil nilai `id` dari `c.Param("id")`.
    - Jika nilai `id` tidak valid, fungsi akan mengembalikan status code `400` dan pesan error dalam format JSON.
    - Jika nilai `id` valid, fungsi akan memanggil fungsi `FetchAll()` pada `studentRepo` untuk mengambil semua data mahasiswa.
    - Jika proses tersebut berhasil, fungsi akan mencari data mahasiswa yang memiliki `id` sesuai dengan nilai yang diberikan.
      - Jika ditemukan, fungsi akan mengembalikan status code `200` dan data mahasiswa dalam format JSON.
      - Namun jika tidak ditemukan, fungsi akan mengembalikan status code `404` dan pesan error dalam format JSON.
    - Namun jika terjadi error pada proses tersebut, fungsi akan mengembalikan status code `500` dan pesan error dalam format JSON.
  - `AddStudent`: fungsi ini akan menambahkan data mahasiswa baru ke dalam database. Pertama-tama, fungsi akan menggunakan `c.ShouldBindJSON` untuk mengambil data mahasiswa yang diberikan dalam request body.
    - Jika data yang diberikan tidak valid, fungsi akan mengembalikan status code `400` dan pesan error dalam format JSON.
    - Jika data yang diberikan valid, fungsi akan memanggil fungsi `Store()` pada `studentRepo` untuk menyimpan data mahasiswa baru ke dalam database.
    - Jika proses tersebut berhasil, fungsi akan mengembalikan status code `200` dan pesan sukses dalam format JSON `{"message: "add course success"}`.
    - Namun jika terjadi error pada proses tersebut, fungsi akan mengembalikan status code `500` dan pesan error dalam format JSON.

- `api/course.go`

  - `AddCourse`: fungsi ini digunakan untuk menambahkan data course baru ke dalam sistem. Data course baru dikirim dalam format JSON melalui request body.
    - Jika request berhasil di-parse ke dalam bentuk objek `model.Course`, maka objek tersebut akan disimpan ke dalam database melalui `api.courseRepo.Store` dan mengembalikan status code `200` dan sebuah JSON response dengan pesan sukses.
    - Jika terjadi error dalam proses parse JSON atau proses penyimpanan ke database, maka API akan mengembalikan response JSON dengan status HTTP `400` Bad Request atau `500` Internal Server Error masing-masing beserta pesan error yang dihasilkan.


### Test Case Examples

#### Test Case 1

**Input**:

```http
POST /student/add HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
    "id": 1,
    "nama": "John Doe",
    "email": "john.doe@example.com",
    "telepon": "1234567890",
    "id_kursus": 123
}
```

**Expected Output / Behavior**:

- Jika permintaan berhasil dan data siswa valid, server harus mengembalikan kode status HTTP `200 OK` dan respons JSON dengan pesan sukses.

  ```json
  {
    "message": "add student success"
  }
  ```

- Jika permintaan gagal karena format data siswa tidak sesuai yang diharapkan, server harus mengembalikan kode status HTTP `400 Bad Request` dan respons JSON dengan pesan kesalahan.

  ```json
  {
    "error": "[error messages]"
  }
  ```

- Jika terjadi kesalahan saat menyimpan data siswa, server harus mengembalikan kode status HTTP `500 Internal Server` Error dan respons JSON dengan pesan kesalahan.

  ```json
  {
    "error": "[error messages]"
  }
  ```

#### Test Case 2

**Input**:

```http
GET /student/get/{id} HTTP/1.1
Host: localhost:8080
```

**Expected Output / Behavior**:

- Jika permintaan berhasil dan ID siswa valid, server harus mengembalikan kode status HTTP `200 OK` dan respons JSON yang berisi data siswa yang sesuai dengan ID yang diberikan.

  ```json
  {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@example.com",
      "phone": "1234567890",
      "course_id": 123
  }
  ```

- Jika permintaan gagal karena ID siswa tidak valid, server harus mengembalikan kode status HTTP `400 Bad Request` dan respons JSON dengan pesan kesalahan.

  ```json
  {
    "error": "[error messages]"
  }
  ```

- Jika terjadi kesalahan saat memproses permintaan, server harus mengembalikan kode status HTTP `500 Internal Server` Error dan respons JSON dengan pesan kesalahan.

  ```json
  {
    "error": "[error messages]"
  }
  ```

- Jika permintaan berhasil namun siswa dengan ID yang diberikan tidak ditemukan, server harus mengembalikan kode status HTTP `404 Not Found` dan respons JSON dengan pesan kesalahan.

  ```json
  {
      "error": "student not found"
  }
  ```
