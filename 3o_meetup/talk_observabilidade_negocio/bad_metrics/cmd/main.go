package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cartAdditionsBad = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "cart_additions_total", Help: "Produtos adicionados ao carrinho"},
	)
	checkoutStartsBad = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "checkout_start_total", Help: "Início do checkout"},
	)
	paymentAttemptsBad = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "payment_attempts_total", Help: "Tentativas de pagamento"},
	)
	purchaseCompletedBad = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "purchase_completed_total", Help: "Compras finalizadas"},
	)
	paymentFailuresBad = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "payment_failures_total", Help: "Falhas no pagamento"},
	)
)

func simulateBusiness(w http.ResponseWriter, r *http.Request) {
	cartAdditionsBad.Inc()
	checkoutStartsBad.Inc()
	paymentAttemptsBad.Inc()
	time.Sleep(time.Duration(rand.Intn(4000)+100) * time.Millisecond)

	if rand.Float64() < 0.2 {
		paymentFailuresBad.Inc()
		w.WriteHeader(500)
		w.Write([]byte("Erro no pagamento"))
		return
	}

	purchaseCompletedBad.Inc()
	w.WriteHeader(200)
	w.Write([]byte("Compra realizada com sucesso"))
}

func main() {
	prometheus.MustRegister(
		cartAdditionsBad,
		checkoutStartsBad,
		paymentAttemptsBad,
		purchaseCompletedBad,
		paymentFailuresBad,
	)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/buy", simulateBusiness)

	log.Println("App com métricas genéricas (ruins) rodando na porta :8087")
	log.Fatal(http.ListenAndServe(":8087", nil))
}
