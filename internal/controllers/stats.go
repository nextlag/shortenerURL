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
	clientIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if clientIP == "" {
		http.Error(w, "X-Real-IP header is required", http.StatusBadRequest)
		return
	}

	trustedSubnet := c.cfg.TrustedSubnet
	if trustedSubnet == "" {
		http.Error(w, "Access forbidden: trusted subnet is not configured", http.StatusForbidden)
		return
	}

	if !ipInSubnet(clientIP, trustedSubnet) {
		http.Error(w, "Access forbidden: client IP is not in trusted subnet", http.StatusForbidden)
		return
	}

	stats, err := c.uc.DoGetStats(r.Context())
	if err != nil {
		http.Error(w, "Failed to load stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(stats)
	if err != nil {
		c.log.Error("Failed to write response", zap.Error(err))
	}
}

// ipInSubnet function to check if an IP address belongs to a network.
func ipInSubnet(ipStr string, subnetStr string) bool {
	ip := net.ParseIP(ipStr)
	_, subnet, err := net.ParseCIDR(subnetStr)
	if err != nil {
		return false
	}
	return subnet.Contains(ip)
}
