package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// GetStatsHandler returns user and link statistics.
// Information is available only to trusted users.
func (c *Controller) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if clientIP == "" {
		c.log.Error("[X-Real-IP] header is empty")
		http.Error(w, "X-Real-IP header is required", http.StatusBadRequest)
		return
	}

	trustedSubnet := c.cfg.TrustedSubnet
	if trustedSubnet == "" {
		c.log.Error("trusted_subnet configuration value not specified")
		http.Error(w, "Access forbidden: trusted subnet is not configured", http.StatusForbidden)
		return
	}

	if !ipInSubnet(clientIP, trustedSubnet, c.log) {
		c.log.Info("client IP is not in trusted subnet", zap.String("client ip: ", clientIP))
		http.Error(w, "Access forbidden: client IP is not in trusted subnet", http.StatusForbidden)
		return
	}

	stats, err := c.uc.DoGetStats(r.Context())
	if err != nil {
		c.log.Error("failed to load stats", zap.Error(err))
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

// ipToBinary преобразует IP-адрес в двоичную строку
func ipToBinary(ip string) (string, error) {
	octets := strings.Split(ip, ".")
	if len(octets) != 4 {
		return "", fmt.Errorf("incorrect IP address: %s", ip)
	}

	binaryIP := ""
	for _, octet := range octets {
		num, err := strconv.Atoi(octet)
		if err != nil || num < 0 || num > 255 {
			return "", fmt.Errorf("invalid IP address: %s", ip)
		}
		binaryIP += fmt.Sprintf("%08b", num)
	}
	return binaryIP, nil
}

// ipInCIDR проверяет, принадлежит ли IP-адрес к CIDR-блоку
func ipInSubnet(ip, cidr string, log *zap.Logger) bool {
	parts := strings.Split(cidr, "/")
	if len(parts) != 2 {
		log.Error(fmt.Sprintf("incorrect CIDR block: %s", cidr))
		return false
	}

	binaryIP, err := ipToBinary(parts[0])
	if err != nil {
		log.Error("error", zap.Error(err))
		return false
	}

	maskLen, err := strconv.Atoi(parts[1])
	if err != nil || maskLen < 0 || maskLen > 32 {
		log.Error(fmt.Sprintf("incorrect mask: %s", parts[1]))
		return false
	}

	binaryCheckIP, err := ipToBinary(ip)
	if err != nil {
		log.Error("error", zap.Error(err))
		return false
	}

	// Сравниваем первые maskLen битов адресов
	return binaryIP[:maskLen] == binaryCheckIP[:maskLen]
}
