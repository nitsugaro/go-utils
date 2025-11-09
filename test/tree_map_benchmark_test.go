package test

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	goutils "github.com/nitsugaro/go-utils"
)

// 🧠 Benchmark: mide el overhead del TreeMap haciendo muchas operaciones
func BenchmarkTreeMapOverhead(b *testing.B) {
	b.ReportAllocs()
	rand.Seed(time.Now().UnixNano())

	for n := 0; n < b.N; n++ {
		m := goutils.NewTreeMap()

		// Inserta 1000 objetos
		for i := 0; i < 1000; i++ {
			key := "user." + strconv.Itoa(i)
			m.Set(key, map[string]any{
				"id":     i,
				"name":   "User_" + strconv.Itoa(i),
				"score":  rand.Float64() * 100,
				"active": i%2 == 0,
			})
		}

		// Lee 1000 valores al azar
		for i := 0; i < 1000; i++ {
			key := "user." + strconv.Itoa(rand.Intn(1000)) + ".name"
			_ = m.Get(key).AsStringOr("")
		}

		// Actualiza 500 valores
		for i := 0; i < 500; i++ {
			key := "user." + strconv.Itoa(rand.Intn(1000)) + ".score"
			m.Set(key, rand.Float64()*200)
		}

		// Clona y accede al clon
		copy := m.Clone()
		for i := 0; i < 200; i++ {
			key := "user." + strconv.Itoa(rand.Intn(1000)) + ".active"
			_ = copy.Get(key).AsBoolOr(false)
		}

		// Elimina 200 claves al azar
		for i := 0; i < 200; i++ {
			key := "user." + strconv.Itoa(rand.Intn(1000))
			m.TryDelete(key)
		}

		// Convierte a JSON
		_ = m.ToJsonString(false)
	}
}

// 💨 Benchmark: compara contra un map[string]any normal
func BenchmarkRawMapOverhead(b *testing.B) {
	b.ReportAllocs()
	rand.Seed(time.Now().UnixNano())

	for n := 0; n < b.N; n++ {
		raw := make(map[string]any)

		// Inserta 1000 objetos
		for i := 0; i < 1000; i++ {
			raw[strconv.Itoa(i)] = map[string]any{
				"id":     i,
				"name":   "User_" + strconv.Itoa(i),
				"score":  rand.Float64() * 100,
				"active": i%2 == 0,
			}
		}

		// Lee 1000 valores al azar
		for i := 0; i < 1000; i++ {
			k := strconv.Itoa(rand.Intn(1000))
			_ = raw[k]
		}

		// Actualiza 500 valores
		for i := 0; i < 500; i++ {
			k := strconv.Itoa(rand.Intn(1000))
			if user, ok := raw[k].(map[string]any); ok {
				user["score"] = rand.Float64() * 200
			}
		}

		// Elimina uno al azar
		delete(raw, strconv.Itoa(rand.Intn(1000)))
	}
}
