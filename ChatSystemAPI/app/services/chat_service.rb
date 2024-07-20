# app/services/application_service.rb
class ChatService
  include Pagination

  def initialize(chat_model: Chat, application_model: Application)
    @chat_model = chat_model
    @application_model = Application
  end

  def create_chat(app_token)
    Rails.logger.info("Chat application token params: #{app_token}")
    @application = @application_model.find_by(token: app_token)
    if @application.nil?
      return { errors: ['Application not found'] }
    end

    application_key = "application_#{@application.id}_last_chat_number"
    Rails.logger.info("Application chats last number cache key: #{application_key}")

    chat_number = Rails.cache.fetch(application_key) do
      Rails.logger.info("Cache miss for key #{application_key}, initializing to 0")
      0
    end

    chat_number += 1
    Rails.logger.info("Chat number incremented to: #{chat_number}")

    Rails.logger.info("Enqueueing create chat job for chat number: #{chat_number} and application_id: #{@application.id}")
    CreateChatJob.perform_async(@application.id, chat_number)

    Rails.logger.info("Updating the cache with the new chat number: #{chat_number} for application: #{@application.id}")
    Rails.cache.write(application_key, chat_number)

    { chat_number: chat_number, application_token: app_token }
  rescue => e
    Rails.logger.error("Error in create_chat: #{e.message}")
    { errors: ['An error occurred while creating chat'] }
  end

  def update_chat(chat, params)
    chat.update(params)
  end

  def find_chat_by_number(number)
    Rails.logger.info("Chat number: #{number}")
    @chat_model.select("number", "messages_count", "created_at", "updated_at").find_by(number: number)
  end

  def all_chat_items(page: 1, per_page: 10)
    @chat_model.select("number", "messages_count", "created_at", "updated_at").paginate(page: page, per_page: per_page)
  end
end
