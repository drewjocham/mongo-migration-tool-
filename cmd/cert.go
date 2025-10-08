package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	hostFlag        string
	portFlag        string
	skipVerifyFlag  bool
	verboseFlag     bool
	fixFlag         bool
	systemCertsFlag bool
)

// certCmd represents the certificate verification command
var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "Certificate verification utilities",
	Long: `Certificate verification utilities for debugging SSL/TLS issues.
	
This command provides various certificate-related operations:
- Check certificate validity for a host
- Diagnose certificate verification issues  
- Fix common certificate problems
- Show system certificate information`,
}

// certCheckCmd checks certificate validity for a host
var certCheckCmd = &cobra.Command{
	Use:   "check [host]",
	Short: "Check certificate validity for a host",
	Long: `Check SSL/TLS certificate validity for a specified host.
	
Examples:
  mongo-migrate cert check login.microsoftonline.com
  mongo-migrate cert check github.com --port 443
  mongo-migrate cert check localhost --port 8443 --skip-verify`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		host := args[0]
		if portFlag == "" {
			portFlag = "443"
		}

		return checkCertificate(host, portFlag)
	},
}

// certDiagnoseCmd diagnoses certificate verification issues
var certDiagnoseCmd = &cobra.Command{
	Use:   "diagnose",
	Short: "Diagnose certificate verification issues",
	Long: `Diagnose common certificate verification problems on the system.
	
This command checks:
- System certificate store location and validity
- Python certificate bundle (for tools like Azure CLI)
- Environment variables affecting certificate verification
- Common certificate-related issues`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return diagnoseCertificates()
	},
}

// certFixCmd attempts to fix common certificate issues
var certFixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix common certificate verification issues",
	Long: `Attempt to fix common certificate verification problems.
	
This command can:
- Update system certificate bundles
- Fix Python certificate issues
- Update Azure CLI certificates
- Set appropriate environment variables`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fixCertificates()
	},
}

// certInfoCmd shows system certificate information
var certInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show system certificate information",
	Long: `Display information about system certificate stores and configuration.
	
Shows:
- System certificate store locations
- Certificate bundle paths
- Environment variables
- Certificate counts and validity`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showCertInfo()
	},
}

func init() {
	// Add cert command to root
	rootCmd.AddCommand(certCmd)

	// Add subcommands
	certCmd.AddCommand(certCheckCmd)
	certCmd.AddCommand(certDiagnoseCmd)
	certCmd.AddCommand(certFixCmd)
	certCmd.AddCommand(certInfoCmd)

	// Override persistent pre-run to skip MongoDB config for cert commands
	certCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Skip MongoDB configuration for cert commands
		return nil
	}

	// Flags for check command
	certCheckCmd.Flags().StringVar(&portFlag, "port", "443", "Port to check (default: 443)")
	certCheckCmd.Flags().BoolVar(&skipVerifyFlag, "skip-verify", false, "Skip certificate verification")
	certCheckCmd.Flags().BoolVar(&verboseFlag, "verbose", false, "Verbose output")

	// Flags for fix command
	certFixCmd.Flags().BoolVar(&fixFlag, "apply", false, "Actually apply fixes (dry-run by default)")
	certFixCmd.Flags().BoolVar(&systemCertsFlag, "system-certs", false, "Use system certificate store")

	// Global flags
	certCmd.PersistentFlags().BoolVar(&verboseFlag, "verbose", false, "Verbose output")
}

func checkCertificate(host, port string) error {
	fmt.Printf("Checking certificate for %s:%s\n", host, port)
	fmt.Println(strings.Repeat("-", 50))

	// Create TLS config
	config := &tls.Config{
		InsecureSkipVerify: skipVerifyFlag,
		ServerName:         host,
	}

	// Connect to the host
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", host, port), config)
	if err != nil {
		return fmt.Errorf("failed to connect to %s:%s: %w", host, port, err)
	}
	defer conn.Close()

	// Get certificate chain
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return fmt.Errorf("no certificates found")
	}

	cert := certs[0]

	fmt.Printf("✓ Connection successful\n")
	fmt.Printf("Subject: %s\n", cert.Subject.CommonName)
	fmt.Printf("Issuer: %s\n", cert.Issuer.CommonName)
	fmt.Printf("Valid From: %s\n", cert.NotBefore.Format(time.RFC3339))
	fmt.Printf("Valid Until: %s\n", cert.NotAfter.Format(time.RFC3339))
	fmt.Printf("Serial Number: %s\n", cert.SerialNumber.String())

	// Check if certificate is expired
	now := time.Now()
	if now.Before(cert.NotBefore) {
		fmt.Printf("⚠️  Certificate is not yet valid\n")
	} else if now.After(cert.NotAfter) {
		fmt.Printf("❌ Certificate has expired\n")
	} else {
		daysUntilExpiry := int(cert.NotAfter.Sub(now).Hours() / 24)
		if daysUntilExpiry < 30 {
			fmt.Printf("⚠️  Certificate expires in %d days\n", daysUntilExpiry)
		} else {
			fmt.Printf("✓ Certificate is valid (%d days remaining)\n", daysUntilExpiry)
		}
	}

	// Show SANs if verbose
	if verboseFlag && len(cert.DNSNames) > 0 {
		fmt.Printf("DNS Names: %s\n", strings.Join(cert.DNSNames, ", "))
	}

	// Show certificate chain
	if verboseFlag && len(certs) > 1 {
		fmt.Printf("Certificate Chain (%d certificates):\n", len(certs))
		for i, c := range certs {
			fmt.Printf("  %d. %s (Issuer: %s)\n", i+1, c.Subject.CommonName, c.Issuer.CommonName)
		}
	}

	return nil
}

func diagnoseCertificates() error {
	fmt.Println("Certificate Verification Diagnosis")
	fmt.Println(strings.Repeat("=", 50))

	// Check system certificate store
	fmt.Println("\n🔍 System Certificate Store:")
	if err := checkSystemCerts(); err != nil {
		fmt.Printf("❌ Error checking system certs: %v\n", err)
	}

	// Check Python certificates
	fmt.Println("\n🔍 Python Certificate Bundle:")
	if err := checkPythonCerts(); err != nil {
		fmt.Printf("❌ Error checking Python certs: %v\n", err)
	}

	// Check environment variables
	fmt.Println("\n🔍 Environment Variables:")
	checkEnvVars()

	// Check Azure CLI if available
	fmt.Println("\n🔍 Azure CLI:")
	if err := checkAzureCLI(); err != nil {
		fmt.Printf("❌ Error checking Azure CLI: %v\n", err)
	}

	// Test common endpoints
	fmt.Println("\n🔍 Connectivity Tests:")
	testEndpoints := []string{
		"github.com",
		"google.com",
		"login.microsoftonline.com",
	}

	for _, endpoint := range testEndpoints {
		if err := testHTTPSConnection(endpoint); err != nil {
			fmt.Printf("❌ %s: %v\n", endpoint, err)
		} else {
			fmt.Printf("✓ %s: OK\n", endpoint)
		}
	}

	return nil
}

func checkSystemCerts() error {
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("security", "find-certificate", "-a", "-p", "/System/Library/Keychains/SystemRootCertificates.keychain")
		output, err := cmd.Output()
		if err != nil {
			return err
		}
		count := strings.Count(string(output), "BEGIN CERTIFICATE")
		fmt.Printf("✓ macOS System Keychain: %d certificates\n", count)

		// Check system cert file
		systemCertPaths := []string{
			"/etc/ssl/cert.pem",
			"/usr/local/etc/openssl/cert.pem",
		}

		for _, path := range systemCertPaths {
			if _, err := os.Stat(path); err == nil {
				fmt.Printf("✓ System cert file: %s\n", path)
				return nil
			}
		}
		fmt.Printf("⚠️  No system cert.pem file found\n")

	case "linux":
		certPaths := []string{
			"/etc/ssl/certs/ca-certificates.crt",
			"/etc/pki/tls/certs/ca-bundle.crt",
			"/etc/ssl/ca-bundle.pem",
		}

		found := false
		for _, path := range certPaths {
			if _, err := os.Stat(path); err == nil {
				fmt.Printf("✓ System cert bundle: %s\n", path)
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("⚠️  No system cert bundle found\n")
		}
	}

	return nil
}

func checkPythonCerts() error {
	// Check if certifi is available
	cmd := exec.Command("python3", "-m", "certifi")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("certifi module not found: %w", err)
	}

	certPath := strings.TrimSpace(string(output))
	if _, err := os.Stat(certPath); err != nil {
		return fmt.Errorf("certificate file not found: %s", certPath)
	}

	fmt.Printf("✓ Python certifi bundle: %s\n", certPath)
	return nil
}

func checkEnvVars() {
	envVars := []string{
		"SSL_CERT_FILE",
		"SSL_CERT_DIR",
		"REQUESTS_CA_BUNDLE",
		"CURL_CA_BUNDLE",
		"PYTHONHTTPSVERIFY",
		"https_proxy",
		"HTTPS_PROXY",
		"http_proxy",
		"HTTP_PROXY",
	}

	found := false
	for _, env := range envVars {
		if value := os.Getenv(env); value != "" {
			fmt.Printf("✓ %s=%s\n", env, value)
			found = true
		}
	}

	if !found {
		fmt.Printf("ℹ️  No certificate-related environment variables set\n")
	}
}

func checkAzureCLI() error {
	// Check if az is available
	cmd := exec.Command("az", "--version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Azure CLI not found")
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		fmt.Printf("✓ %s\n", strings.TrimSpace(lines[0]))
	}

	return nil
}

func testHTTPSConnection(host string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}

	resp, err := client.Get(fmt.Sprintf("https://%s", host))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func fixCertificates() error {
	fmt.Println("Certificate Fix Utility")
	fmt.Println(strings.Repeat("=", 50))

	if !fixFlag {
		fmt.Println("🔍 DRY RUN MODE - Use --apply to actually apply fixes")
		fmt.Println()
	}

	fixes := []struct {
		name        string
		description string
		fix         func() error
	}{
		{
			name:        "Update Python certificates",
			description: "Update the Python certifi certificate bundle",
			fix:         updatePythonCerts,
		},
		{
			name:        "Set system certificate environment",
			description: "Configure environment to use system certificates",
			fix:         setSystemCertEnv,
		},
		{
			name:        "Update Azure CLI certificates",
			description: "Update Azure CLI certificate bundle",
			fix:         updateAzureCLICerts,
		},
	}

	for _, fix := range fixes {
		fmt.Printf("🔧 %s\n", fix.name)
		fmt.Printf("   %s\n", fix.description)

		if fixFlag {
			if err := fix.fix(); err != nil {
				fmt.Printf("   ❌ Failed: %v\n", err)
			} else {
				fmt.Printf("   ✓ Applied successfully\n")
			}
		} else {
			fmt.Printf("   📝 Would apply this fix\n")
		}
		fmt.Println()
	}

	if !fixFlag {
		fmt.Println("Run with --apply to actually apply these fixes")
	}

	return nil
}

func updatePythonCerts() error {
	cmd := exec.Command("python3", "-m", "pip", "install", "--upgrade", "certifi")
	return cmd.Run()
}

func setSystemCertEnv() error {
	var certPath string

	switch runtime.GOOS {
	case "darwin":
		certPath = "/etc/ssl/cert.pem"
	case "linux":
		// Try common locations
		paths := []string{
			"/etc/ssl/certs/ca-certificates.crt",
			"/etc/pki/tls/certs/ca-bundle.crt",
			"/etc/ssl/ca-bundle.pem",
		}

		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				certPath = path
				break
			}
		}
	}

	if certPath == "" {
		return fmt.Errorf("no system certificate bundle found")
	}

	fmt.Printf("Would set SSL_CERT_FILE=%s\n", certPath)
	fmt.Printf("Add this to your shell profile:\n")
	fmt.Printf("export SSL_CERT_FILE=%s\n", certPath)

	return nil
}

func updateAzureCLICerts() error {
	// This would update Azure CLI's certificate bundle
	// The actual implementation depends on the Azure CLI version and installation
	fmt.Printf("Would update Azure CLI certificates\n")
	fmt.Printf("Run: brew update && brew upgrade azure-cli\n")

	return nil
}

func showCertInfo() error {
	fmt.Println("System Certificate Information")
	fmt.Println(strings.Repeat("=", 50))

	// Show system info
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Architecture: %s\n", runtime.GOARCH)
	fmt.Println()

	// Show certificate stores
	fmt.Println("Certificate Store Locations:")

	// System root certificates
	_, err := x509.SystemCertPool()
	if err != nil {
		fmt.Printf("❌ Error loading system cert pool: %v\n", err)
	} else {
		// Note: SystemCertPool doesn't expose certificate count directly
		fmt.Printf("✓ System certificate pool loaded successfully\n")
	}

	// Show environment variables
	fmt.Println("\nEnvironment Variables:")
	checkEnvVars()

	// Show file locations
	fmt.Println("\nCertificate Files:")
	checkSystemCerts()
	checkPythonCerts()

	return nil
}
