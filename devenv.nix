{ pkgs, lib, config, inputs, ... }:

{
  dotenv.enable = true;

  # https://devenv.sh/basics/
  env.GREET = "devenv";

  # https://devenv.sh/packages/
  packages = [ pkgs.curl pkgs.git pkgs.buf pkgs.cfssl pkgs.sqlfluff ];

  # https://devenv.sh/scripts/
  scripts.hello.exec = "echo hello from $GREET";

  enterShell = ''
    hello
    git --version
  '';

  # https://devenv.sh/tests/
  enterTest = ''
    echo "Running tests"
    git --version | grep "2.42.0"
  '';

  # https://devenv.sh/services/
  services.redis.enable = true;
  services.postgres.enable = true;
  services.postgres.listen_addresses = "127.0.0.1";

  # https://devenv.sh/languages/
  # languages.nix.enable = true;
  languages.go.enable = true;

  # https://devenv.sh/pre-commit-hooks/
  # pre-commit.hooks.shellcheck.enable = true;

  # https://devenv.sh/processes/
  processes.dbrunner-service = {
    exec = "go run ./cmd/dbrunner-service";
    process-compose = {
      depends_on = {
        redis.condition = "process_healthy";
      };
      readiness_probe = {
        exec.command = "curl --cacert scripts/cert/ca-dev.pem --cert scripts/cert/client.pem --key scripts/cert/client-key.pem https://localhost:3000/healthz";
      };
      environment = [
        "PORT=3000"
      ];
    };
  };
  processes.question-manager-service = {
    exec = "go run ./cmd/question-manager-service";
    process-compose = {
      depends_on = {
        postgres.condition = "process_healthy";
      };
      readiness_probe = {
        exec.command = "curl --cacert scripts/cert/ca-dev.pem --cert scripts/cert/client.pem --key scripts/cert/client-key.pem https://localhost:3001/healthz";
      };
      environment = [
        "PORT=3001"
      ];
    };
  };
    processes.gateway-service = {
    exec = "go run ./cmd/gateway-service";
    process-compose = {
      depends_on = {
        # dbrunner-service.condition = "process_healthy";
        question-manager-service.condition = "process_healthy";
      };
      readiness_probe = {
        exec.command = "curl http://localhost:3100/healthz";
      };
      environment = [
        "PORT=3100"
      ];
    };
  };

  # See full reference at https://devenv.sh/reference/options/
}
