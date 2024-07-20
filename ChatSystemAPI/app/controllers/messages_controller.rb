class MessagesController < ApplicationController
  def initialize(chat_service = ChatService.new(chat_model: Chat, application_model: Application), message_service = MessageService.new)
    @message_service = message_service
  end

  def search
    search_term = params[:body]
    chat_number = (params[:chat_number]).to_i
    application_token = (params[:application_token])

    if chat_number.blank? || application_token.blank?
      render json: { errors: ['chat_number and application_token query params are required'] }, status: :bad_request
      return
    end

    page = (params[:page] || 1).to_i
    per_page = (params[:per_page] || 10).to_i

    results = @message_service.search_messages(
      search: search_term,
      chat_number: chat_number,
      application_token: application_token,
      page: page,
      per_page: per_page
    )

    render json: { items: results }
  end

end
