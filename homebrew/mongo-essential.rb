class MongoEssential < Formula
  desc "Essential MongoDB toolkit with migrations and AI-powered database analysis"
  homepage "https://github.com/jocham/mongo-essential"
  url "https://github.com/jocham/mongo-essential/archive/v1.0.0.tar.gz"
  sha256 "SHA256_PLACEHOLDER" # Will be updated when creating the actual release
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=#{version}"), "."
    
    # Install shell completions
    generate_completions_from_executable(bin/"mongo-essential", "completion")
    
    # Install documentation
    doc.install "README.md", "AI_ANALYSIS.md", "CHANGELOG.md"
    
    # Install example configuration
    pkgshare.install ".env.example"
  end

  service do
    run [opt_bin/"mongo-essential", "up"]
    environment_variables MONGO_URL: "mongodb://localhost:27017", MONGO_DATABASE: "your_database"
    error_log_path var/"log/mongo-essential.log"
    log_path var/"log/mongo-essential.log"
    working_dir var
  end

  test do
    # Test version command
    assert_match version.to_s, shell_output("#{bin}/mongo-essential version")
    
    # Test help command
    assert_match "mongo-essential", shell_output("#{bin}/mongo-essential --help")
    
    # Test certificate diagnostic
    assert_match "Certificate Verification Diagnosis", shell_output("#{bin}/mongo-essential cert diagnose 2>&1", 1)
    
    # Test AI command help
    assert_match "AI-powered database analysis", shell_output("#{bin}/mongo-essential ai --help")
  end
end
