package provider

import (
	"context"
	"testing"

	"github.com/NX211/traefik-proxmox-provider/internal"
)

func TestProviderConfig(t *testing.T) {
	config := CreateConfig()
	if config.PollInterval != "30s" {
		t.Errorf("Expected default PollInterval to be '30s', got %s", config.PollInterval)
	}
	if config.ApiValidateSSL != "true" {
		t.Errorf("Expected default ApiValidateSSL to be 'true', got %s", config.ApiValidateSSL)
	}
	if config.ApiLogging != "info" {
		t.Errorf("Expected default ApiLogging to be 'info', got %s", config.ApiLogging)
	}
}

func TestProviderNew(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "Valid config",
			config: &Config{
				PollInterval:   "5s",
				ApiEndpoint:    "https://proxmox.example.com",
				ApiTokenId:     "test@pam!test",
				ApiToken:       "test-token",
				ApiValidateSSL: "true",
				ApiLogging:     "info",
			},
			wantErr: true, // We expect an error because the domain doesn't exist
		},
		{
			name:    "Nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "Missing poll interval",
			config: &Config{
				ApiEndpoint:    "https://proxmox.example.com",
				ApiTokenId:     "test@pam!test",
				ApiToken:       "test-token",
				ApiValidateSSL: "true",
				ApiLogging:     "info",
			},
			wantErr: true,
		},
		{
			name: "Invalid poll interval",
			config: &Config{
				PollInterval:   "invalid",
				ApiEndpoint:    "https://proxmox.example.com",
				ApiTokenId:     "test@pam!test",
				ApiToken:       "test-token",
				ApiValidateSSL: "true",
				ApiLogging:     "info",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := New(context.Background(), tt.config, "test-provider")
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && provider != nil {
				t.Error("Expected provider to be nil when there's an error")
			}
		})
	}
}

func TestProviderValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "Valid config",
			config: &Config{
				PollInterval:   "5s",
				ApiEndpoint:    "https://proxmox.example.com",
				ApiTokenId:     "test@pam!test",
				ApiToken:       "test-token",
				ApiValidateSSL: "true",
				ApiLogging:     "info",
			},
			wantErr: false,
		},
		{
			name:    "Nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "Missing poll interval",
			config: &Config{
				ApiEndpoint:    "https://proxmox.example.com",
				ApiTokenId:     "test@pam!test",
				ApiToken:       "test-token",
				ApiValidateSSL: "true",
				ApiLogging:     "info",
			},
			wantErr: true,
		},
		{
			name: "Missing endpoint",
			config: &Config{
				PollInterval:   "5s",
				ApiTokenId:     "test@pam!test",
				ApiToken:       "test-token",
				ApiValidateSSL: "true",
				ApiLogging:     "info",
			},
			wantErr: true,
		},
		{
			name: "Missing token ID",
			config: &Config{
				PollInterval:   "5s",
				ApiEndpoint:    "https://proxmox.example.com",
				ApiToken:       "test-token",
				ApiValidateSSL: "true",
				ApiLogging:     "info",
			},
			wantErr: true,
		},
		{
			name: "Missing token",
			config: &Config{
				PollInterval:   "5s",
				ApiEndpoint:    "https://proxmox.example.com",
				ApiTokenId:     "test@pam!test",
				ApiValidateSSL: "true",
				ApiLogging:     "info",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProviderParserConfig(t *testing.T) {
	tests := []struct {
		name        string
		apiEndpoint string
		tokenID     string
		token       string
		wantErr     bool
	}{
		{
			name:        "Valid config",
			apiEndpoint: "https://proxmox.example.com",
			tokenID:     "test@pam!test",
			token:       "test-token",
			wantErr:     false,
		},
		{
			name:        "Missing endpoint",
			apiEndpoint: "",
			tokenID:     "test@pam!test",
			token:       "test-token",
			wantErr:     true,
		},
		{
			name:        "Missing token ID",
			apiEndpoint: "https://proxmox.example.com",
			tokenID:     "",
			token:       "test-token",
			wantErr:     true,
		},
		{
			name:        "Missing token",
			apiEndpoint: "https://proxmox.example.com",
			tokenID:     "test@pam!test",
			token:       "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := newParserConfig(tt.apiEndpoint, tt.tokenID, tt.token, "debug", true)
			if (err != nil) != tt.wantErr {
				t.Errorf("newParserConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if config.ApiEndpoint != tt.apiEndpoint {
					t.Errorf("Expected ApiEndpoint to be %s, got %s", tt.apiEndpoint, config.ApiEndpoint)
				}
				if config.TokenId != tt.tokenID {
					t.Errorf("Expected TokenId to be %s, got %s", tt.tokenID, config.TokenId)
				}
				if config.Token != tt.token {
					t.Errorf("Expected Token to be %s, got %s", tt.token, config.Token)
				}
			}
		})
	}
}

func TestProviderService(t *testing.T) {
	config := map[string]string{
		"traefik.enable":                 "true",
		"traefik.http.routers.test.rule": "Host(`test.example.com`)",
	}

	service := internal.NewService(123, "test-service", config)
	if service.ID != 123 {
		t.Errorf("Expected service ID to be 123, got %d", service.ID)
	}
	if service.Name != "test-service" {
		t.Errorf("Expected service name to be 'test-service', got %s", service.Name)
	}
	if len(service.Config) != 2 {
		t.Errorf("Expected service config to have 2 items, got %d", len(service.Config))
	}
	if len(service.IPs) != 0 {
		t.Errorf("Expected service IPs to be empty, got %d items", len(service.IPs))
	}
}

func TestGetServiceURL(t *testing.T) {
	tests := []struct {
		name        string
		service     internal.Service
		serviceName string
		nodeName    string
		expectedUrl string
	}{
		{
			name:        "IP set, default port and scheme",
			serviceName: "service",
			service: internal.Service{
				Config: map[string]string{
					"traefik.http.services.service.loadbalancer.server.ip": "1.2.3.4",
				},
			},
			expectedUrl: "http://1.2.3.4:80",
		},
		{
			name:        "IP and scheme set, default port (http)",
			serviceName: "service",
			service: internal.Service{
				Config: map[string]string{
					"traefik.http.services.service.loadbalancer.server.ip":     "1.2.3.4",
					"traefik.http.services.service.loadbalancer.server.scheme": "http",
				},
			},
			expectedUrl: "http://1.2.3.4:80",
		},
		{
			name:        "IP and scheme set, default port (https)",
			serviceName: "service",
			service: internal.Service{
				Config: map[string]string{
					"traefik.http.services.service.loadbalancer.server.ip":     "1.2.3.4",
					"traefik.http.services.service.loadbalancer.server.scheme": "https",
				},
			},
			expectedUrl: "https://1.2.3.4:443",
		},
		{
			name:        "IP, port and scheme set",
			serviceName: "service",
			service: internal.Service{
				Config: map[string]string{
					"traefik.http.services.service.loadbalancer.server.ip":     "1.2.3.4",
					"traefik.http.services.service.loadbalancer.server.scheme": "https",
					"traefik.http.services.service.loadbalancer.server.port":   "8080",
				},
			},
			expectedUrl: "https://1.2.3.4:8080",
		},
		{
			name:        "URL is set",
			serviceName: "service",
			service: internal.Service{
				Config: map[string]string{
					"traefik.http.services.service.loadbalancer.server.url": "http://test.com:1234",
				},
			},
			expectedUrl: "http://test.com:1234",
		},
		{
			name:        "URL overrides other settings",
			serviceName: "service",
			service: internal.Service{
				Config: map[string]string{
					"traefik.http.services.service.loadbalancer.server.url":    "http://test.com:1234",
					"traefik.http.services.service.loadbalancer.server.ip":     "1.2.3.4",
					"traefik.http.services.service.loadbalancer.server.scheme": "https",
					"traefik.http.services.service.loadbalancer.server.port":   "8080",
				},
			},
			expectedUrl: "http://test.com:1234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := getServiceURL(tt.service, tt.serviceName, tt.nodeName)
			if url != tt.expectedUrl {
				t.Errorf("Expected URL to be %s, got %s", tt.expectedUrl, url)
			}
		})
	}
}

func TestGenerateConfiguration_NodeWithExplicitNames(t *testing.T) {
	servicesMap := map[string][]internal.Service{
		"pve1": {
			{
				ID:   0,
				Name: "pve1",
				IPs:  []internal.IP{{Address: "192.168.1.10", AddressType: "ipv4"}},
				Config: map[string]string{
					"traefik.enable":                                                   "true",
					"traefik.http.routers.proxmox-ui.rule":                             "Host(`proxmox.example.com`)",
					"traefik.http.routers.proxmox-ui.entrypoints":                      "websecure",
					"traefik.http.routers.proxmox-ui.tls.certresolver":                 "letsencrypt",
					"traefik.http.services.proxmox-ui.loadbalancer.server.port":        "8006",
					"traefik.http.services.proxmox-ui.loadbalancer.server.scheme":      "https",
				},
			},
		},
	}

	config := generateConfiguration(servicesMap)

	// Check router was created
	router, exists := config.HTTP.Routers["proxmox-ui"]
	if !exists {
		t.Fatal("Expected router 'proxmox-ui' to exist")
	}
	if router.Rule != "Host(`proxmox.example.com`)" {
		t.Errorf("Expected rule Host(`proxmox.example.com`), got %s", router.Rule)
	}
	if router.Service != "proxmox-ui" {
		t.Errorf("Expected service 'proxmox-ui', got %s", router.Service)
	}

	// Check service was created
	svc, exists := config.HTTP.Services["proxmox-ui"]
	if !exists {
		t.Fatal("Expected service 'proxmox-ui' to exist")
	}
	if len(svc.LoadBalancer.Servers) != 1 {
		t.Fatalf("Expected 1 server, got %d", len(svc.LoadBalancer.Servers))
	}
	if svc.LoadBalancer.Servers[0].URL != "https://192.168.1.10:8006" {
		t.Errorf("Expected URL https://192.168.1.10:8006, got %s", svc.LoadBalancer.Servers[0].URL)
	}
}

func TestGenerateConfiguration_NodeWithDefaultNaming(t *testing.T) {
	servicesMap := map[string][]internal.Service{
		"pve1": {
			{
				ID:   0,
				Name: "pve1",
				IPs:  []internal.IP{{Address: "10.0.0.1", AddressType: "ipv4"}},
				Config: map[string]string{
					"traefik.enable": "true",
				},
			},
		},
	}

	config := generateConfiguration(servicesMap)

	// Default ID should be "pve1-0"
	_, routerExists := config.HTTP.Routers["pve1-0"]
	if !routerExists {
		t.Fatal("Expected router 'pve1-0' to exist")
	}

	_, svcExists := config.HTTP.Services["pve1-0"]
	if !svcExists {
		t.Fatal("Expected service 'pve1-0' to exist")
	}
}

func TestGenerateConfiguration_NodeSkippedWithoutEnable(t *testing.T) {
	servicesMap := map[string][]internal.Service{
		"pve1": {
			{
				ID:   0,
				Name: "pve1",
				IPs:  []internal.IP{{Address: "10.0.0.1", AddressType: "ipv4"}},
				Config: map[string]string{
					"traefik.http.routers.test.rule": "Host(`test.example.com`)",
				},
			},
		},
	}

	config := generateConfiguration(servicesMap)

	if len(config.HTTP.Routers) != 0 {
		t.Errorf("Expected no routers, got %d", len(config.HTTP.Routers))
	}
	if len(config.HTTP.Services) != 0 {
		t.Errorf("Expected no services, got %d", len(config.HTTP.Services))
	}
}

func TestGenerateConfiguration_NodeAndVMCoexist(t *testing.T) {
	servicesMap := map[string][]internal.Service{
		"pve1": {
			// VM service
			{
				ID:   100,
				Name: "webserver",
				IPs:  []internal.IP{{Address: "192.168.1.50", AddressType: "ipv4"}},
				Config: map[string]string{
					"traefik.enable":                                              "true",
					"traefik.http.routers.web.rule":                               "Host(`web.example.com`)",
					"traefik.http.services.web.loadbalancer.server.port":          "80",
				},
			},
			// Node service
			{
				ID:   0,
				Name: "pve1",
				IPs:  []internal.IP{{Address: "192.168.1.10", AddressType: "ipv4"}},
				Config: map[string]string{
					"traefik.enable":                                                   "true",
					"traefik.http.routers.proxmox.rule":                                "Host(`proxmox.example.com`)",
					"traefik.http.services.proxmox.loadbalancer.server.port":           "8006",
					"traefik.http.services.proxmox.loadbalancer.server.scheme":         "https",
				},
			},
		},
	}

	config := generateConfiguration(servicesMap)

	// Check VM router
	if _, exists := config.HTTP.Routers["web"]; !exists {
		t.Error("Expected router 'web' to exist")
	}

	// Check node router
	if _, exists := config.HTTP.Routers["proxmox"]; !exists {
		t.Error("Expected router 'proxmox' to exist")
	}

	// Check VM service
	webSvc, exists := config.HTTP.Services["web"]
	if !exists {
		t.Fatal("Expected service 'web' to exist")
	}
	if webSvc.LoadBalancer.Servers[0].URL != "http://192.168.1.50:80" {
		t.Errorf("Expected URL http://192.168.1.50:80, got %s", webSvc.LoadBalancer.Servers[0].URL)
	}

	// Check node service
	proxmoxSvc, exists := config.HTTP.Services["proxmox"]
	if !exists {
		t.Fatal("Expected service 'proxmox' to exist")
	}
	if proxmoxSvc.LoadBalancer.Servers[0].URL != "https://192.168.1.10:8006" {
		t.Errorf("Expected URL https://192.168.1.10:8006, got %s", proxmoxSvc.LoadBalancer.Servers[0].URL)
	}
}

func TestHandleRouterTLS_ArrayDomains(t *testing.T) {
	tests := []struct {
		name           string
		config         map[string]string
		expectedMain   []string
		expectedSANs   [][]string
		expectNil      bool
	}{
		{
			name: "Array syntax with main and sans",
			config: map[string]string{
				"traefik.http.routers.test.tls.domains[0].main": "example.com",
				"traefik.http.routers.test.tls.domains[0].sans": "*.example.com,www.example.com",
				"traefik.http.routers.test.tls.domains[1].main": "another.com",
			},
			expectedMain: []string{"example.com", "another.com"},
			expectedSANs: [][]string{{"*.example.com", "www.example.com"}, nil},
		},
		{
			name: "Simple domains fallback",
			config: map[string]string{
				"traefik.http.routers.test.tls.domains": "example.com,another.com",
			},
			expectedMain: []string{"example.com", "another.com"},
			expectedSANs: [][]string{nil, nil},
		},
		{
			name:      "No TLS config",
			config:    map[string]string{},
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := internal.Service{Config: tt.config}
			tlsConfig := handleRouterTLS(service, "traefik.http.routers.test")

			if tt.expectNil {
				if tlsConfig != nil {
					t.Error("Expected nil TLS config")
				}
				return
			}

			if tlsConfig == nil {
				t.Fatal("Expected non-nil TLS config")
			}

			if len(tlsConfig.Domains) != len(tt.expectedMain) {
				t.Fatalf("Expected %d domains, got %d", len(tt.expectedMain), len(tlsConfig.Domains))
			}

			for i, domain := range tlsConfig.Domains {
				if domain.Main != tt.expectedMain[i] {
					t.Errorf("Domain[%d].Main = %s, want %s", i, domain.Main, tt.expectedMain[i])
				}
				if tt.expectedSANs[i] != nil {
					if len(domain.SANs) != len(tt.expectedSANs[i]) {
						t.Errorf("Domain[%d].SANs length = %d, want %d", i, len(domain.SANs), len(tt.expectedSANs[i]))
					}
				}
			}
		})
	}
}
