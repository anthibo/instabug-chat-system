Rails.application.config.after_initialize do
  Thread.new { MessageCreatedEventConsumer.start }
end
