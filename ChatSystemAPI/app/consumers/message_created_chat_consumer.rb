class MessageCreatedChatConsumer
  def self.start
    exchange = RabbitMQ.exchange('message_created', :fanout)
    queue = RabbitMQ.queue('message_created_update_message_count_queue')
    queue.bind(exchange)

    queue.subscribe(manual_ack: true) do |delivery_info, properties, payload|
      begin
        message = JSON.parse(payload)
        Rails.logger.info "Message created event: #{message.inspect}"

        chat = Chat.find_by(id: message['chat_id'])
        if chat
          chat.update(messages_count: chat.messages_count + 1)
          Rails.logger.info "Chat message count updated for chat: #{chat.id}"
        end
        RabbitMQ.channel.ack(delivery_info.delivery_tag)
      rescue StandardError => e
        Rails.logger.error "Failed to update chat message count: #{e.message}"
      end
    end
  end
end
