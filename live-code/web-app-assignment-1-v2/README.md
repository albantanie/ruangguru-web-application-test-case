# Web Application

## Assignment 1 - Task Tracker Plus

### Description

**Task Tracker Plus** adalah aplikasi pengelolaan tugas yang dirancang untuk membantu mahasiswa dalam mengatur jadwal studi mereka. Dibangun menggunakan bahasa pemrograman Go (Golang).

Dalam assignment ini, kita akan mengimplementasikan API menggunakan _Golang web framework Gin_ untuk mengelola data _task_ dan _category_. API harus mengizinkan client untuk:

- Menambahkan tugas baru dengan mengirimkan permintaan **POST** ke endpoint `/task/add`.
- Mengambil informasi tugas berdasarkan ID dengan mengirimkan permintaan **GET** ke endpoint `/task/get/:id`.
- Memperbarui informasi tugas dengan mengirimkan permintaan **PUT** ke endpoint `/task/update/:id`.
- Menghapus tugas dengan mengirimkan permintaan **DELETE** ke endpoint `/task/delete/:id`.
- Mendapatkan daftar tugas dengan mengirimkan permintaan **GET** ke endpoint `/task/list`.
- Mendapatkan daftar tugas berdasarkan ID kategori dengan mengirimkan permintaan **GET** ke endpoint `/task/category/:id`.
- Menambahkan kategori baru dengan mengirimkan permintaan **POST** ke endpoint `/category/add`.
- Mengambil informasi kategori berdasarkan ID dengan mengirimkan permintaan **GET** ke endpoint `/category/get/:id`.
- Memperbarui informasi kategori dengan mengirimkan permintaan **PUT** ke endpoint `/category/update/:id`.
- Menghapus kategori dengan mengirimkan permintaan **DELETE** ke endpoint `/category/delete/:id`.
- Mendapatkan daftar kategori dengan mengirimkan permintaan **GET** ke endpoint `/category/list`.

Gunakan fungsi pada subpackage di `db/filebased` untuk berhubungan dengan database, seluruh fungsinya dapat dipelajari di `db/filebased/README.md` dan juga kamu bisa membaca sendiri kode yang ada di dalamnya.

### Instruction

Pada Assignment ini, kamu harus melengkapi fungsi dari repository `tasks` dan `course` ini dengan mengimplementasikan fungsi-fungsi berikut:

📁 **repository**

Ini adalah fungsi yang berinteraksi dengan database:

> Warning : abaikan code yang tidak berhubungan dengan instruksi di bawah ini pada folder `repository`

- `repository/task.go`

  - `Update`: fungsi akan mengeksekusi sebuah query `UPDATE` untuk memperbarui data tugas tersebut di dalam tabel `tasks`. Query tersebut akan menggunakan nilai dari variabel `model.Task` yang diberikan sebagai argumen.
    - Jika proses tersebut berhasil, fungsi akan mengembalikan `nil` sebagai `error`.
    - Namun jika terjadi error pada proses tersebut, fungsi akan mengembalikan `error` yang terjadi.
  - `Delete`: fungsi untuk menghapus data tugas berdasarkan `id` yang diberikan sebagai argumen. Fungsi ini akan mengeksekusi sebuah query `DELETE` untuk menghapus data tugas dengan `id` tersebut dari tabel `tasks`.
    - Jika proses tersebut berhasil, fungsi akan mengembalikan `nil` sebagai `error`.
    - Namun jika terjadi error pada proses tersebut, fungsi akan mengembalikan `error` yang terjadi.
  - `GetList`: fungsi untuk mengambil daftar data tugas dari tabel `tasks`. Fungsi akan mengeksekusi sebuah query `SELECT` untuk mengambil semua data tugas yang ada dalam tabel.
    - Jika proses tersebut berhasil, fungsi akan mengembalikan slice `[]*model.Task` yang berisi data tugas yang ditemukan, serta `nil` sebagai `error`.
    - Namun jika terjadi error pada proses tersebut, fungsi akan mengembalikan `nil` sebagai slice `[]*model.Task`, serta `error` yang terjadi.
  - `GetTaskCategory`: fungsi untuk mengambil kategori tugas berdasarkan `id` yang diberikan sebagai argumen. Fungsi akan mengeksekusi sebuah query `SELECT` dengan menggunakan join tabel antara `tasks` dan `categories` untuk mengambil informasi kategori tugas yang sesuai dengan `id` yang diberikan.
    - Jika proses tersebut berhasil, fungsi akan mengembalikan slice `[]*model.TaskCategory` yang berisi kategori tugas yang ditemukan, serta `nil` sebagai `error`.
    - Namun jika terjadi error pada proses tersebut, fungsi akan mengembalikan `nil` sebagai slice `[]*model.TaskCategory`, serta `error` yang terjadi.

- `repository/category.go`
  - `Update`: fungsi untuk memperbarui kategori berdasarkan `id` yang diberikan sebagai argumen. Fungsi akan mengeksekusi sebuah query `UPDATE` untuk memperbarui data kategori yang sesuai dengan `id` yang diberikan dalam tabel `categories`. Query tersebut akan menggunakan nilai dari variabel `model.Category` yang diberikan sebagai argumen.
    - Jika proses tersebut berhasil, fungsi akan mengembalikan `nil` sebagai `error`.
    - Namun jika terjadi error pada proses tersebut, fungsi akan mengembalikan `error` yang terjadi.
  - `Delete`: fungsi untuk menghapus kategori berdasarkan `id` yang diberikan sebagai argumen. Fungsi akan mengeksekusi sebuah query `DELETE` untuk menghapus data kategori dengan `id` tersebut dari tabel `categories`.
    - Jika proses tersebut berhasil, fungsi akan mengembalikan `nil` sebagai `error`.
    - Namun jika terjadi error pada proses tersebut, fungsi akan mengembalikan `error` yang terjadi.
  - `GetList`: fungsi untuk mengambil daftar kategori dari tabel `categories`. Fungsi akan mengeksekusi sebuah query `SELECT` untuk mengambil semua data kategori yang ada dalam tabel.
    - Jika proses tersebut berhasil, fungsi akan mengembalikan slice `[]*model.Category` yang berisi data kategori yang ditemukan, serta `nil` sebagai `error`.
    - Namun jika terjadi error pada proses tersebut, fungsi akan mengembalikan `nil` sebagai slice `[]*model.Category`, serta `error` yang terjadi.

📁 **service**

Layer service digunakan untuk memproses data sesuai dengan aturan bisnis yang telah ditentukan.

- `service/task.go`
  - `Update`: fungsi ini menggunakan repository `taskRepository` untuk memperbarui data tugas.
    - Jika pembaruan berhasil, fungsi mengembalikan `nil` sebagai `error`
    - Jika terjadi error, fungsi mengembalikan `error` yang terjadi.
  - `Delete`: fungsi ini menggunakan repository `taskRepository` untuk menghapus data tugas berdasarkan `id`.
    - Jika penghapusan berhasil, fungsi mengembalikan `nil` sebagai `error`.
    - Jika terjadi error, fungsi mengembalikan `error` yang terjadi.
  - `GetList`: fungsi ini menggunakan repository `taskRepository` untuk mengambil daftar tugas.
    - Jika pengambilan daftar tugas berhasil, fungsi mengembalikan daftar tugas (`[]*model.Task`) dan `nil` sebagai `error`.
    - Jika terjadi error saat pengambilan daftar tugas, fungsi mengembalikan nil sebagai daftar tugas (`[]*model.Task`) dan `error` yang terjadi.
  - `GetTaskCategory`: fungsi ini menggunakan repository `taskRepository` untuk mengambil daftar tugas beserta dengan kategorinya.
    - Jika pengambilan kategori tugas berhasil, fungsi mengembalikan kategori tugas (`[]*model.TaskCategory`) dan `nil` sebagai `error`.
    - Jika terjadi error saat pengambilan kategori tugas, fungsi mengembalikan `nil` sebagai kategori tugas (`[]*model.TaskCategory`) dan `error` yang terjadi.
- `service/category.go`
  - `Update`: fungsi ini menggunakan repository `categoryRepository` untuk memperbarui data kategori.
    - Jika pembaruan berhasil, fungsi mengembalikan `nil` sebagai `error`
    - Jika terjadi error, fungsi mengembalikan `error` yang terjadi.
  - `Delete`: fungsi ini menggunakan repository `categoryRepository` untuk menghapus data kategori berdasarkan `id`.
    - Jika penghapusan berhasil, fungsi mengembalikan `nil` sebagai `error`.
    - Jika terjadi error, fungsi mengembalikan `error` yang terjadi.
  - `GetList`: fungsi ini menggunakan repository `categoryRepository` untuk mengambil daftar kategori.
    - Jika pengambilan daftar kategori berhasil, fungsi mengembalikan daftar kategori (`[]*model.Category`) dan `nil` sebagai `error`.
    - Jika terjadi error saat pengambilan daftar kategori, fungsi mengembalikan `nil` sebagai daftar kategori (`[]*model.Category`) dan `error` yang terjadi.

📁 **api**

> Warning : abaikan code yang tidak berhubungan dengan instruksi di bawah ini pada folder `api`

- `api/task.go`

  - `UpdateTask`: fungsi ini digunakan untuk mengupdate data tugas dengan `id` yang sesuai dengan nilai yang diberikan sebagai argumen.
    - Jika nilai `id` tidak valid, fungsi akan mengirimkan respons dengan kode status `400` dan pesan error dalam format JSON `{"error": "Invalid task ID"}`.
    - Jika data yang diberikan tidak valid, fungsi akan mengirimkan respons dengan kode status `400` dan pesan error dalam format JSON.
    - Selanjutnya, fungsi memanggil `Update` dari `taskService` untuk memperbarui tugas dengan `taskID` yang diberikan sebagai argumen. Assign `taskID` ke struct `model.Task.ID` yang akan dikirim ke fungsi tersebut.
      - Jika proses berhasil, fungsi akan mengirimkan respons dengan kode status `200` dan pesan sukses dalam format JSON `{"message": "update task success"}`.
      - Jika terjadi error pada proses tersebut, fungsi akan mengirimkan respons dengan kode status `500` dan pesan error dalam format JSON.
  - `DeleteTask`: fungsi ini digunakan untuk menghapus tugas dengan `id` yang sesuai dengan nilai yang diberikan sebagai argumen.
    - Jika nilai `id` tidak valid, fungsi akan mengirimkan respons dengan kode status `400` dan pesan error dalam format JSON `{"error": "Invalid task ID"}`.
    - Selanjutnya, fungsi memanggil `Delete` dari `taskService` untuk menghapus tugas dengan `taskID` yang diberikan sebagai argumen.
      - Jika terjadi error saat proses penghapusan, fungsi akan mengirimkan tanggapan JSON dengan kode status HTTP `500` Internal Server Error dan pesan error.
      - Jika tidak terjadi error dalam proses penghapusan, fungsi akan mengirimkan tanggapan JSON dengan kode status HTTP `200` OK dan pesan sukses `{"message": "delete task success"}`.
  - `GetTaskList`: fungsi ini digunakan untuk mendapatkan daftar tugas.

    - Pertama, fungsi memanggil `GetList` dari `taskService` untuk mendapatkan daftar tugas.

      - Jika terjadi error saat mengambil daftar tugas, fungsi akan mengirimkan tanggapan JSON dengan kode status HTTP `500` Internal Server Error dan pesan error yang sesuai ke client.
      - Jika tidak terjadi error, fungsi akan mengirimkan tanggapan JSON dengan kode status HTTP `200` OK dan daftar tugas ke client. Contoh:

        ```json
        [
          {
            "id": 1,
            "title": "Task 1",
            "deadline": "2023-05-30",
            "priority": 2,
            "category_id": "1",
            "status": "In Progress"
          },
          {
            "id": 2,
            "title": "Task 2",
            "deadline": "2023-06-10",
            "priority": 1,
            "category_id": "2",
            "status": "Completed"
          },
          ...
        ]
        ```

  - `GetTaskCategory`: fungsi ini digunakan untuk mendapatkan daftar tugas dengan nama kategorinya.

    - Pertama, fungsi memanggil `GetTaskCategory` dari `taskService` untuk mendapatkan daftar tugas.

      - Jika terjadi error saat mengambil daftar tugas, fungsi akan mengirimkan tanggapan JSON dengan kode status HTTP `500` Internal Server Error dan pesan error yang sesuai ke client.
      - Jika tidak terjadi error, fungsi akan mengirimkan tanggapan JSON dengan kode status HTTP `200` OK dan daftar tugas beserta nama kategorinya ke client. Contoh:

        ```json
        [
          {
            "id": 1,
            "title": "Task 1",
            "category": "Category 1"
          },
          {
            "id": 2,
            "title": "Task 2",
            "category": "Category 2"
          },
          ...
        ]
        ```

- `api/category.go`

  - `UpdateCategory`: fungsi ini digunakan untuk mengupdate kategori dengan `id` yang sesuai dengan nilai yang diberikan sebagai argumen.
    - Jika nilai `id` tidak valid _(gagal dikonversi menjadi tipe data int)_, fungsi akan mengirimkan respons JSON dengan kode status HTTP `400` Bad Request dan pesan error dalam format JSON `{"error": "invalid Category ID"}`.
    - Jika data yang diberikan tidak valid, fungsi akan mengirimkan respons JSON dengan kode status HTTP `400` Bad Request dan pesan error dalam format JSON.
    - Jika data yang diberikan valid, fungsi akan memanggil fungsi `Update` pada `categoryService` untuk mengupdate kategori dengan ID yang sesuai dengan nilai yang diberikan.
      - Jika proses update berhasil, fungsi akan mengirimkan respons JSON dengan kode status HTTP `200` OK dan pesan sukses dalam format JSON `{"message": "category update success"}`.
      - Namun jika terjadi error pada proses update, fungsi akan mengirimkan respons JSON dengan kode status HTTP `500` Internal Server Error dan pesan error dalam format JSON.
  - `DeleteCategory`: fungsi ini digunakan untuk menghapus kategori dengan `id` yang sesuai dengan nilai yang diberikan sebagai argumen.
    - Jika nilai `id` tidak valid _(gagal dikonversi menjadi tipe data int)_, fungsi akan mengirimkan respons JSON dengan kode status HTTP `400` Bad Request dan pesan error dalam format JSON `{"error": "Invalid category ID"}`.
    - Jika nilai `id` valid, fungsi akan memanggil fungsi `Delete` pada `categoryService` untuk menghapus kategori dengan `ID` yang sesuai dengan nilai yang diberikan.
      - Jika proses penghapusan berhasil, fungsi akan mengirimkan respons JSON dengan kode status HTTP `200` OK dan pesan sukses dalam format JSON `{"message": "category delete success"}`.
      - Namun jika terjadi error pada proses penghapusan, fungsi akan mengirimkan respons JSON dengan kode status HTTP `500` Internal Server Error dan pesan error dalam format JSON.
  - `GetCategoryList`: fungsi ini digunakan untuk mendapatkan daftar kategori.

    - Fungsi memanggil `GetList` dari `categoryService` untuk mendapatkan daftar kategori.

      - Jika terjadi error saat mengambil daftar kategori, fungsi akan mengirimkan respons JSON dengan kode status HTTP `500` Internal Server Error dan pesan error yang sesuai ke client.
      - Jika tidak terjadi error, fungsi akan mengirimkan respons JSON dengan kode status HTTP `200` OK dan daftar kategori ke client. Contoh:

        ```json
        [
          {
            "id": 1,
            "name": "Category 1"
          },
          {
            "id": 2,
            "name": "Category 2"
          },
          ...
        ]
        ```

### Test Case Examples

#### Test Case 1

**Input**:

```http
PUT /task/update/{id} HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
    "title": "Task 1",
    "deadline": "2023-06-30",
    "priority": 2,
    "category_id": "1",
    "status": "In Progress"
}
```

**Expected Output / Behavior**:

- Jika permintaan berhasil dan ID tugas valid, server harus mengembalikan kode status HTTP `200` OK dan respons JSON dengan pesan sukses.

  ```json
  {
    "message": "update task success"
  }
  ```

- Jika permintaan gagal karena ID tugas tidak valid, server harus mengembalikan kode status HTTP `400` Bad Request dan respons JSON dengan pesan kesalahan.

  ```json
  {
    "error": "Invalid task ID"
  }
  ```

- Jika permintaan gagal karena format data tugas tidak sesuai yang diharapkan, server harus mengembalikan kode status HTTP `400` Bad Request dan respons JSON dengan pesan kesalahan.

  ```json
  {
    "error": "[error messages]"
  }
  ```

- Jika terjadi kesalahan saat memperbarui data tugas, server harus mengembalikan kode status HTTP `500` Internal Server Error dan respons JSON dengan pesan kesalahan.

  ```json
  {
    "error": "[error messages]"
  }
  ```

#### Test Case 2

**Input**:

```http
DELETE /task/delete/{id} HTTP/1.1
Host: localhost:8080
```

**Expected Output / Behavior**:

- Jika permintaan berhasil dan ID tugas valid, server harus mengembalikan kode status HTTP `200` OK dan respons JSON dengan pesan sukses.

  ```json
  {
    "message": "delete task success"
  }
  ```

- Jika permintaan gagal karena ID tugas tidak valid, server harus mengembalikan kode status HTTP `400` Bad Request dan respons JSON dengan pesan kesalahan.

  ```json
  {
    "error": "Invalid task ID"
  }
  ```

- Jika terjadi kesalahan saat menghapus tugas, server harus mengembalikan kode status HTTP `500` Internal Server Error dan respons JSON dengan pesan kesalahan.

  ```json
  {
    "error": "[error messages]"
  }
  ```

#### Test Case 3

**Input**:

```http
GET /task/list HTTP/1.1
Host: localhost:8080
```

**Expected Output / Behavior**:

- Jika permintaan berhasil, server harus mengembalikan kode status HTTP `200` OK dan respons JSON dengan daftar tugas.

  ```json
  [
      {
          "id": 1,
          "title": "Task 1",
          "deadline": "2023-05-30",
          "priority": 2,
          "category_id": "1",
          "status": "In Progress"
      },
      {
          "id": 2,
          "title": "Task 2",
          "deadline": "2023-06-10",
          "priority": 1,
          "category_id": "2",
          "status": "Completed"
      },
      ...
  ]
  ```

- Jika terjadi kesalahan saat mendapatkan daftar tugas, server harus mengembalikan kode status HTTP `500` Internal Server Error dan respons JSON dengan pesan kesalahan.

  ```json
  {
    "error": "[error messages]"
  }
  ```
