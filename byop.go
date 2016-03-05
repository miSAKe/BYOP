package main

import (
        "github.com/prometheus/client_golang/prometheus"
        "github.com/traetox/goDS18B20"
        "net/http"
        "fmt"
)

func main() {
        http.Handle("/metrics", prometheus.Handler())
        
        http.ListenAndServe(":8080", nil)
}

func init() {
        goDS18B20.Setup()

        slaves, _ := goDS18B20.Slaves()
        for i := range slaves {
                probe, _ := goDS18B20.NewProbe(slaves[i])

                guage := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
                        Subsystem: "BYOP",
                        Name:      fmt.Sprintf("sensor_%v_temperature_celcius", i),
                        Help:      fmt.Sprintf("Sensor %v probe temperature in degrees celcius.", i),
                }, func() float64 {
                        probe.Update()
                        temperature, _ := probe.Temperature()
                        return float64(temperature.Celsius())
                })

                prometheus.MustRegister(guage)
        }

}