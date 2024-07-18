Rails.application.config.after_initialize do
  Thread.new { MessageCreatedElasticsearchConsumer.start }
  # Thread.new { MessageCreatedChatConsumer.start }
end
