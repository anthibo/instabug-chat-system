class MessageCreatedEventConsumer
  def self.start
    RabbitMQ.exchange('', :direct)
    queue = RabbitMQ.queue('message_created_queue')

    queue.subscribe(manual_ack: true) do |delivery_info, properties, payload|
      begin
        message_data = JSON.parse(payload)
        Rails.logger.info "Message created event: #{message_data.inspect}"

        MessageSearch.index_message(message_data)

        Rails.logger.info "Message search index updated for message_number: #{message_data['message_number'] }"
        RabbitMQ.channel.ack(delivery_info.delivery_tag)
      rescue StandardError => e
        Rails.logger.error "Failed to update Elasticsearch index: #{e.message}"
      end
    end
  end
end
