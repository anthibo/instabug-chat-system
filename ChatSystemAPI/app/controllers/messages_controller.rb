class MessagesController < ApplicationController
  def initialize(chat_service = ChatService.new(chat_model: Chat, application_model: Application), message_service = MessageService.new)
    @message_service = message_service
  end

  def search
    search_term = params[:search]
    chat_number = params[:chat_number]
    page = params[:page] || 1
    per_page = params[:per_page] || 10

    chat = Chat.find_by!(number: chat_number)
    results = @message_service.search_messages(search_term, chat.id, page, per_page)
    render json: results
  end

end
