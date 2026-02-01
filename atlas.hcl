data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./cmd/loader",
  ]
}

env "local" {
  src = data.external_schema.gorm.url
  
  # Use a temporary Docker container for the dev database to ensure a clean state
  dev = "docker://postgres/16/dev?search_path=public"
  
  migration {
    dir = "file://migrations"
  }
}
