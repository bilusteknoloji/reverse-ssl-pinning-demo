task :command_exists, [:command] do |_, args|
  abort "#{args.command} doesn't exists" if `command -v #{args.command} > /dev/null 2>&1 && echo $?`.chomp.empty?
end
task :is_repo_clean do
  abort 'please commit your changes first!' unless `git status -s | wc -l`.strip.to_i.zero?
end
task :has_bumpversion do
  Rake::Task['command_exists'].invoke('bumpversion')
end
task :has_curl do
  Rake::Task['command_exists'].invoke('curl')
end
task :has_python do
  Rake::Task['command_exists'].invoke('python')
end

AVAILABLE_REVISIONS = %w[major minor patch].freeze
task :bump, [:revision] => [:has_bumpversion] do |_, args|
  args.with_defaults(revision: 'patch')
  unless AVAILABLE_REVISIONS.include?(args.revision)
    abort "Please provide valid revision: #{AVAILABLE_REVISIONS.join(',')}"
  end

  system "bumpversion #{args.revision}"
  exit $?.exitstatus
end

desc "release new version #{AVAILABLE_REVISIONS.join(',')}, default: patch"
task :release, [:revision] => [:is_repo_clean] do |_, args|
  args.with_defaults(revision: 'patch')
  Rake::Task['bump'].invoke(args.revision)
end

namespace :run do
  desc "run server"
  task :server do
    system %{
      go run cmd/server/main.go
    }
  end

  desc "run go client"
  task :go_client do
    system %{
      go run cmd/client/main.go
    }
  end

  desc "run curl client"
  task :curl_client => [:has_curl] do
    system %{
      $(command -v curl) --cacert certs/server/server-cert.pem \
        --cert certs/client/client-cert-signed.pem \
        --key certs/client/client-key.pem \
        https://localhost:8443
    }
  end

  desc "run python client"
  task :python_client => [:has_python] do
    system %{
      if [[ -d "venv" ]]; then
          if [ -z "${VIRTUAL_ENV}" ]; then
            source venv/bin/activate
          fi
          if ! pip freeze | grep -q 'requests'; then
            pip install requests
          fi
          python client.py
      else
          echo "you need to create virtual environment"
          echo "python -m venv venv"
          exit 1
      fi
      echo $VIRTUAL_ENV
    }
  end
end
