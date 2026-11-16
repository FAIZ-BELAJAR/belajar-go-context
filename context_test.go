package belajar_golang_context

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// Pengenalan Context
// Context merupakan sebuah data yang membawa value, sinyal cancel, sinyal timeout, dan sinyal deadline
// Context biadanya dibuat per request(misal setiap ada request masuk ke server web melalui http request)
// Context digunakan untuk mempermudah kita meneruskan value, dan sinyal antar proses

// Kenapa Context Perlu dipelajari
// Context di goang biasa digunakan untuk mengirim data request atau sinyal ke proses lain
// Dengan menggunakan context,ketika kita ingin membatalkan semua proses, kita cukup mengirim sinyal ke context, maka secara otomatis semua proses akan dibatalkan
// Hampir semua bagian di golang memanfaatkan context,seperti database, http server, http client, dan lain lain
// Bahkan di google sendiri, ketika menggunakan golang, context wajib digunakan dan selalu dikirim ke setiap function yang dikirim

// Package Context
// Context direpresentasikan di dalam sebuah interface Context
// interface Context terdapat di dalam package context
// https://golang.org/pkg/context/\

// Context memiiki 4 Method utama yakni
// Deadline() (deadline time.Tim, ok bool)
// Done() <-chan struct{}
// Err error
// Value(key interface{}) interface{}

// Membuat Context
// Karena context adalah sebuah interface, untuk membuat context kita butuh sebuah struct yang sesuai dengan kontrrak interface context
// Namun kit atidak perlu membuatnya secara manual
// Di Golang package context terdapat function yang bisa kita gunakan untuk membuat Context

// Function untuk membuat context
// context.Background(): Membuat context kosong, tidak pernah dibatalkan, tidak pernah timeout, dan tidak memiliki value apapun, biasanya digunakan di main function atau dalam test, atau dalam awal proses request terjadi
// context.TODO: Membuat context kosong seperti context.Background(), namun biasanya menggunakan ini ketika belum jelas context apa yang ingin digunakan
func TestMembuatContext(t *testing.T) {
	background := context.Background()
	todo := context.TODO()

	fmt.Println(background)
	fmt.Println(todo)

}

// Parent dan Child Context
// Context menganut konsep parent child
// Artinya saat kita membuat sebuah context kita bisa membuat child context yang sudah ada
// parent context bisa membuat banyak child, namun child hanya bisa memiliki 1 parent
// konsep ini mirip dengan pewarisan di pemrograman berorientasi object

// Hubungan Antara Parent dan child context
// Parent dan Child conetxt akan selalu terhubung
// Saat kita nanti melakukan misal pembatalan context A, maka semua child dan sub child dari context A, akan ikut dibatalkan
// Namun jika misal kita membatalkan context B, hanya contextB dan semua child dan sub child nya yang dibatalkan, Parent context nya tak akan dibatalkan
// Begitu juga nanti saat kita menyisipkan data ke dalam context A, semua child dan sub child nya bisa mendapatkan data tersebut
// Namun jika kita menyisipkan data di context B, hanya context B dan semua child dan sub child nya yang mendapat data parent context B tidak akan mendapat data

// Immutable
// Context merupakan object yang immutable, artinya setelah context dibuat, dia tidak bisa diubah lagi
// Ketika kita menambahkan value ke dalam context, atau menambahkan pengaturan timeout dan yang lainnya, secara otomatis akan membentuk child context baru, bukan merubah context tersebut

// Context with Value
// Pada saat awal kita membuat context, context tidak memilki value
// Kita bisa menambahkan sebuah value, dengan data Pair(key-value) ke dalam context
// Saat kita menambah value ke context, secara otomatis akan tercipta child context baru, artinya original context nya tidak akan berubah sama sekali
// Untuk membuat menambahkan value ke context, kita bisa menggunakan function context.WithValue(parent, key, value)
// Kita bisa memngambil context value milik sendiri dan parent namun tak bisa mengambil value milik child
// Menggunakan method Value(key)
func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "Faiz", "Hannan")
	contextC := context.WithValue(contextA, "Faiz22", "Hannan22")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)

	fmt.Println(contextD.Value("Faiz"))
}

// Context With Cancel
// Selain menambahkan value ke context, kita juga bisa menambahkan sinyal cancel ke context
// Kapan sinyal cancel diperlukan dalam contextt?
// Biasanya ketika kita butuh menjalankan proses lain, dan kita ingin bisa memberi sinyal cancel ke proses tersebut
// Biasanya proses ini berupa goroutine yang berbeda sehingga dengan mudah jika kita ingin membatalkan eksekusi goroutine, kita bisa mengirim sinyal cancel ke context nya
// Namun ingat, goroutine yang menggunakan context, tetap harus melakukan pengecekan terhadap contextnya, jika tidak, tidak ada gunanya
// Untuk membuat context dengan cancel kita bisa menggunakan function di package context
// context.WithCancel(parent)

func TestWithCancel(t *testing.T) {
	parent := context.Background()
	var wg sync.WaitGroup

	deadline := time.Now().Add(time.Second * 5)
	konteks, cancel := context.WithDeadline(parent, deadline)
	defer cancel()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		i := 0
		for {
			i++
			select {
			case <-ctx.Done():
				fmt.Println("Perulangan selesai:", ctx.Err()) // ctx.Err() akan bernilai context.DeadlineExceeded
				return
			default:
				fmt.Println("Perulangan ke: ", i)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}(konteks)

	wg.Wait()
	fmt.Println("Program selesai")

}
