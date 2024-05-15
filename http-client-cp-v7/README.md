# Exercise HTTP Client - GET & POST

## Task

Kita sangat tertarik dengan banyak quote dari orang terkenal yang bisa memotivasi diri kita. Oleh karena itu kita sebagai seorang backend developer ingin membuat aplikasi yang menampilkan quote tersebut menggunakan 3rd Party API <https://api.quotable.io/quotes/random>.

> Catatan, Quotes API ini memiliki rate limit tertentu, jika kamu mendapatkan error 429 dari API endpoint, maka berhenti sebentar dan cobalah beberapa menit kemudian.

- Jadi buatlah client Golang dengan menggunakan method **GET** ke API ini: <https://api.quotable.io/quotes/random?limit=3> dan lakukan `decode` terhadap response API tersebut ke dalam bentuk Array terhadap struct berikut:

  ```go
  type Quotes struct {
      Tags []string
      Author string
      Quote string
  }
  ```

- Lalu, buatlah juga client Golang yang menggunakan method **POST** ke api ini: <https://postman-echo.com/post> dengan body json berikut:

  ```json
  {
      "name":  "Dion",
      "email": "dionbe2022@gmail.com",
  }
  ```

  Dan lakukan `decode` terhadap response API tersebut terhadap struct berikut:

  ```go
  type Postman struct {
      Data data 
      Url string
  }
  ```

Berikut _template_ pengerjaannya:

```go
// Lakukan `decode` dari fungsi `ClientGet()` terhadap struct berikut ini:
type Quotes struct {
    Tags []string
    Author string
    Quote string
}

func ClientGet() ([]Quotes, error) {
    // Kerjakan di sini

    return []Quotes, nil // TODO: replace
}

// Lakukan `decode` dari fungsi `ClientPost()` terhadap struct berikut ini:
type data struct {
    Email string
    Name string
}

type Postman struct {
    Data data 
    Url string
}

func ClientPost() (Postman, error) {
    postBody, _ := json.Marshal(map[string]string{
        "name":  "Dion",
        "email": "dionbe2022@gmail.com",
    })
    responseBody := bytes.NewBuffer(postBody)
    fmt.Println(responseBody)

    // Kerjakan di sini

    return Postman, nil // TODO: replace
}
```

Encode **API** to **Struct**:

- API with method **GET**:

> Ini hanyalah sebuah contoh, struktur data pasti sama, namun isi data bisa jadi berbeda.

    ```json
    // curl -X GET "https://api.quotable.io/quotes/random?limit=3"
    [
    {
        "_id": "un7fRc7cVivc",
        "content": "Patience and perseverance have a magical effect before which difficulties disappear and obstacles vanish.",
        "author": "John Adams",
        "tags": [
            "Famous Quotes"
        ],
        "authorSlug": "john-adams",
        "length": 105,
        "dateAdded": "2019-06-27",
        "dateModified": "2023-04-14"
    },
    {
        ...
    },
    ...
    ]
    ```

    to

    ```go
    type Quotes struct {
        Tags string `json:"tags"`
        Author string `json:"author"`
        Quote string `json:"content"`
    }
    ```

- API with method **POST**:

    ```json
    // curl -X POST https://postman-echo.com/post -H "Content-Type: application/json" -d '{"name":"Dion","email":"dionbe2022@gmail.com"}
    {
    "args": {},
    "data": "'{name:Dion,email:dionbe2022@gmail.com}'",
    "files": {},
    "form": {},
    "headers": {
        "x-forwarded-proto": "https",
        "x-forwarded-port": "443",
        "host": "postman-echo.com",
        "x-amzn-trace-id": "Root=1-6329467e-3bb7e4bf7ea9d31b0d6a3edf",
        "content-length": "40",
        "user-agent": "curl/7.83.1",
        "accept": "*/*",
        "content-type": "application/json"
    },
    "json": null,
    "url": "https://postman-echo.com/post"
    }
    ```

    to

    ```go
    type data struct {
        Email string `json:"email"`
        Name string `json:"name"`
    }

    type Postman struct {
        Data data 
        Url string `json:"url"`
    }
    ```

Contoh test case :

> Perhatikan struktur data-nya, isi data-nya bisa jadi berbeda.

```bash
# Method GET:
Input: https://api.quotable.io/quotes/random?limit=3
Output: [{[Wisdom] Francis Bacon Wise men make more opportunities than they find.} {...} ...]

# Method POST:
Input: https://postman-echo.com/post -H "Content-Type: application/json" -d '{"name":"Dion","email":"dionbe2022@gmail.com"}'
Output: {{dionbe2022@gmail.com Dion} https://postman-echo.com/post}
```
  