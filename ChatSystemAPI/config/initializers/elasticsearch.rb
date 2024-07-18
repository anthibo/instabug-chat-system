require "elasticsearch"
Elasticsearch::Model.client = Elasticsearch::Client.new(
  url: "http://localhost:9200",
  retry_on_failure: 5,
  request_timeout: 30,
  log: Rails.env.development?
)


