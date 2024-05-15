# Exercise file 1

## Calculate profit or loss

Terdapat perusahaan A yang menjual barang dagangan berupa makanan dan minuman. Perusahaan tersebut melakukan pencatatan transaksi sederhana secara manual di dalam sebuah file.

Setiap file akan berisi seluruh transaksi dalam satu bulan. Transaksi yang dicacat hanya dikategorikan sebagai pemasukan atau "income" dan pengeluaran atau "expense". Format transaksinya berupa `[dd/mm/yyyy];[type];[amount]`.

Contoh data transaksi:

- `transactions.txt`

```txt
1/1/2021;income;100000
1/1/2021;expense;50000
2/1/2021;income;200000
2/1/2021;expense;100000
3/1/2021;income;300000
3/1/2021;expense;150000
4/1/2021;income;400000
4/1/2021;expense;200000
...
```

Buatlah sebuah fungsi yang dapat digunakan untuk mendapatkan data dari file tersebut, kemudian menghitung total keuntungan ("_profit_") atau kerugian ("_loss_") dari perhitungan transaksi dalam satu bulan. Kembalikan nilai dalam format yang sama dengan tanggal terakhir dari transaksi tersebut. Contoh format jika total _income_ lebih besar dari _expense_, maka seperti ini `"1/1/2021;profit;100000"`, jika total _expense_ lebih besar dari _income_, maka seperti ini `"1/1/2021;loss;100000"`.

Contoh, jika di dalam `transaction.txt` berisi data berikut:

```txt
1/1/2021;income;100000
1/1/2021;expense;50000
2/1/2021;income;200000
2/1/2021;expense;100000
3/1/2021;income;300000
3/1/2021;expense;150000
4/1/2021;income;400000
4/1/2021;expense;200000
```

Kita perlu menghitung total keuntungan atau kerugian dari tanggal 1 Januari 2021 sampai 4 Januari 2021.

100000 + 50000 + 200000 - 100000 + 300000 - 150000 + 400000 - 200000 = 500000 (_income_)

Karena, hasil hasil total _income_ lebih besar dari total _expense_, maka terjadi keuntungan. Kita dapat mengembalikan nilai `"4/1/2021;profit;500000"`.

Berikut _template_ yang diperlukan:

```go

// gunakan fungsi ini untuk mendapatkan data file
// kembalikan data berupa slice of string
func Readfile(path string) ([]string, error) {
    // kerjakan di sini
}

func CalculateProfitLoss(data string) string {
    // kerjakan di sini
}
```

> Hint : kalian dapat memisah setiap baris data di `transactions.txt` dengan pemisah `"\n"` (_new line_)
