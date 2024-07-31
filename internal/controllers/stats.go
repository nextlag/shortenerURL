package controllers

import (
	"net"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// GetStatsHandler returns user and link statistics.
// Information is available only to trusted users.
func (c *Controller) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	ipStr := r.Header.Get("X-Real-IP")
	if ipStr == "" {
		ips := r.Header.Get("X-Forwarded-For")
		if ips == "" {
			c.log.Error("Missing both X-Real-IP and X-Forwarded-For headers")
			http.Error(w, "Missing both X-Real-IP and X-Forwarded-For headers", http.StatusBadRequest)
			return
		}
		ipStrs := strings.Split(ips, ",")
		ipStr = strings.TrimSpace(ipStrs[0])
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		c.log.Error("Invalid IP address in headers")
		http.Error(w, "Invalid IP address in headers", http.StatusBadRequest)
		return
	}

	if c.cfg.TrustedSubnet == "" {
		c.log.Error("Trusted subnet configuration value not specified")
		http.Error(w, "Access forbidden: trusted subnet is not configured", http.StatusForbidden)
		return
	}

	_, cidr, err := net.ParseCIDR(c.cfg.TrustedSubnet)
	if err != nil {
		c.log.Error("invalid trusted subnet configuration", zap.Error(err))
		http.Error(w, "Access forbidden: invalid trusted subnet configuration", http.StatusForbidden)
		return
	}

	if !cidr.Contains(ip) {
		c.log.Info("Client IP is not in trusted subnet", zap.String("client_ip", ip.String()))
		http.Error(w, "Access forbidden: client IP is not in trusted subnet", http.StatusForbidden)
		return
	}

	stats, err := c.uc.DoGetStats(r.Context())
	if err != nil {
		c.log.Error("Failed to load stats", zap.Error(err))
		http.Error(w, "Failed to load stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(stats); err != nil {
		c.log.Error("Failed to write response", zap.Error(err))
	}
}
