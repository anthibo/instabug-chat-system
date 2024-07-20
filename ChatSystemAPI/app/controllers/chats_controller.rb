class ChatsController < ApplicationController
  before_action :find_chat_by_number, only: [:show]

  def initialize(chat_service = ChatService.new(chat_model: Chat, application_model: Application), message_service = MessageService.new)
    @chat_service = chat_service
    @message_service = message_service
  end

  def index
    page = (params[:page].presence || 1).to_i
    per_page = (params[:per_page].presence || 10).to_i

    @chats = @chat_service.all_chat_items(page: page, per_page: per_page)

    render json: { items: @chats }
  end

  def show
    render json: @chat, status: :ok
  end

  def create
    app_token = params[:application_token]
    response = @chat_service.create_chat(app_token)
    if response[:errors].present?
      render json: response, status: :not_found
      return
    end

    render json: response, status: :created
  end

  private

  def find_chat_by_number
    number = params[:number]
    chat = @chat_service.find_chat_by_number(number)
    if chat.nil?
      render json: { errors: ['Chat not found'] }, status: :not_found
    else
      @chat = chat
    end
  end
end
