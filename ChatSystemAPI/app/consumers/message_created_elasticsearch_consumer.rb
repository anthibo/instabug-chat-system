class MessageCreatedElasticsearchConsumer
  def self.start
    exchange = RabbitMQ.exchange('message_created', :fanout)
    queue = RabbitMQ.queue('message_created_elasticsearch_queue')
    queue.bind(exchange)

    queue.subscribe(manual_ack: true) do |delivery_info, properties, payload|
      begin
        message_data = JSON.parse(payload)
        Rails.logger.info "Message created event: #{message_data.inspect}"

        MessageSearch.index_message(message_data)

        Rails.logger.info "Message search index updated for message: #{message_data['id']}"
        RabbitMQ.channel.ack(delivery_info.delivery_tag)
      rescue StandardError => e
        Rails.logger.error "Failed to update Elasticsearch index: #{e.message}"
      end
    end
  end
end
