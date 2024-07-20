class CreateChatJob
  include Sidekiq::Job

  sidekiq_options queue: 'create_chat_queue'

  def perform(application_id, chat_number)
    Rails.logger.info("Saving Chat Job: application_id: #{application_id}, chat number: #{chat_number}")
    chat = Chat.new(application_id: application_id, number: chat_number)
    if chat.save
      Rails.logger.info("Chat created successfully: #{chat.inspect}")
    else
      Rails.logger.error("Failed to save chat: #{chat.errors.full_messages.join(', ')}")
    end
  end
end
