require "elasticsearch"
# log url of elastic search
Elasticsearch::Model.client = Elasticsearch::Client.new(
  url: ENV["ELASTICSEARCH_URL"] || "http://localhost:9200",
  retry_on_failure: 10,
  request_timeout: 30,
  retry_delay: 15,
  log: Rails.env.development?
)
