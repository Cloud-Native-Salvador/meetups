package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	paymentLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "payment_processing_duration_seconds",
			Help:    "Tempo de processamento de pagamento",
			Buckets: prometheus.ExponentialBuckets(0.1, 2, 8),
		},
		[]string{"payment_method", "status_code"},
	)
	paymentAttempts = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "payment_attempts_total", Help: "Tentativas de pagamento"},
		[]string{"payment_method"},
	)
	cartAdditions = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "cart_additions_total", Help: "Produtos adicionados ao carrinho"},
		[]string{"product_id", "user_tier"},
	)
	checkoutStarts = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "checkout_start_total", Help: "Início do checkout"},
		[]string{"channel"},
	)
	paymentFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "payment_failures_total", Help: "Falhas no pagamento"},
		[]string{"reason"},
	)
	purchaseCompleted = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "purchase_completed_total", Help: "Compras finalizadas"},
		[]string{"region", "payment_method"},
	)
	userReturnRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "user_return_rate", Help: "Taxa de retorno de usuários"},
		[]string{"segment"},
	)
	paymentStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "payment_status_total", Help: "Pagamentos por status"},
		[]string{"payment_method", "status"},
	)
)

func simulateBusiness(w http.ResponseWriter, r *http.Request) {
	method := r.URL.Query().Get("method")
	if method == "" {
		method = "credit_card"
	}
	product := "abc123"
	region := "br"
	segment := "premium"
	statusCode := 200
	duration := time.Duration(rand.Intn(4000)+100) * time.Millisecond
	time.Sleep(duration)

	paymentAttempts.WithLabelValues(method).Inc()

	if method == "boleto" {
		paymentStatus.WithLabelValues("boleto", "pending").Set(1)
		paymentLatency.WithLabelValues("boleto", "202").Observe(duration.Seconds())

		go func() {
			time.Sleep(5 * time.Second)
			outcome := rand.Float64()
			switch {
			case outcome < 0.5:
				paymentStatus.WithLabelValues("boleto", "pending").Set(0)
				paymentStatus.WithLabelValues("boleto", "confirmed").Set(1)
				purchaseCompleted.WithLabelValues(region, method).Inc()
			case outcome < 0.8:
				paymentStatus.WithLabelValues("boleto", "pending").Set(0)
				paymentStatus.WithLabelValues("boleto", "confirmed").Set(1)
				paymentFailures.WithLabelValues("expired").Inc()
			default:
				paymentStatus.WithLabelValues("boleto", "pending").Set(0)
				paymentStatus.WithLabelValues("boleto", "cancelled").Set(1)
				paymentFailures.WithLabelValues("cancelled").Inc()
			}
		}()
		w.WriteHeader(202)
		w.Write([]byte("Pagamento em boleto gerado - pendente"))
		return
	}

	if rand.Float64() < 0.2 {
		statusCode = 500
		reasons := []string{"timeout", "gateway_error", "insufficient_funds"}
		reason := reasons[rand.Intn(len(reasons))]
		paymentFailures.WithLabelValues(reason).Inc()
	} else {
		purchaseCompleted.WithLabelValues(region, method).Inc()
	}

	paymentLatency.WithLabelValues(method, fmt.Sprintf("%d", statusCode)).Observe(duration.Seconds())
	cartAdditions.WithLabelValues(product, segment).Inc()
	checkoutStarts.WithLabelValues("mobile").Inc()
	userReturnRate.WithLabelValues(segment).Set(rand.Float64())

	w.WriteHeader(statusCode)
	w.Write([]byte("Simulação de fluxo de compra completa"))
}

func main() {
	prometheus.MustRegister(
		paymentLatency,
		paymentAttempts,
		cartAdditions,
		checkoutStarts,
		paymentFailures,
		purchaseCompleted,
		userReturnRate,
		paymentStatus,
	)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/buy", simulateBusiness)

	log.Println("App com métricas orientadas ao negócio rodando na porta :8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
