# app/services/application_service.rb
class ChatService
  include Pagination

  def initialize(chat_model: Chat, application_model: Application)
    @chat_model = chat_model
    @application_model = Application
  end

  def create_chat(app_token)
    Rails.logger.info("Chat application token params: #{app_token}")
    application = @application_model.find_by(token: app_token)
    if application.nil?
      Rails.logger.error("Application not found for token: #{app_token}")
      return { errors: ['Application not found'] }
    end

    chat = @chat_model.new(application: application)
    if chat.save
      Rails.logger.info("Chat created successfully: #{chat.inspect}")
      chat
    else
      Rails.logger.error("Failed to save chat: #{chat.errors.full_messages.join(', ')}")
      { errors: chat.errors.full_messages }
    end
  end

  def update_chat(chat, params)
    chat.update(params)
  end

  def find_chat_by_number(number)
    Rails.logger.info("Chat number: #{number}")
    @chat_model.find_by(number: number)
  end

  def all_chat_items(page: 1, per_page: 10)
    paginate(@chat_model, page: page, per_page: per_page)
  end
end
