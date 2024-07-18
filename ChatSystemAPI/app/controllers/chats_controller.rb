class ChatsController < ApplicationController
  before_action :find_chat_by_id, only: [:show]

  def initialize(chat_service = ChatService.new(chat_model: Chat, application_model: Application), message_service = MessageService.new)
    @chat_service = chat_service
    @message_service = message_service
  end

  def index
    page = (params[:page].presence || 1).to_i
    per_page = (params[:per_page].presence || 10).to_i
    Rails.logger.info("Page: #{page}, Per page: #{per_page}")
    @chats = @chat_service.all_chat_items(page: page, per_page: per_page)
    render json: @chats
  end

  def show
    render json: @chat, status: :ok
  end

  def create
    # TODO: add error handling for invalid chat number
    app_token = params[:application_token]
    Rails.logger.info("Create Chat Params: #{params.inspect}")
    Rails.logger.info("Chat application token: #{app_token.inspect}")
    result = @chat_service.create_chat(app_token)

    if result.is_a?(Chat)
      render json: result, status: :created
    else
      render json: { errors: result[:errors] }, status: :bad_request
    end
  end

  private
  def find_chat_by_id
    number = params[:number]
    chat = @chat_service.find_chat_by_number(number)
    if chat.nil?
      render json: { errors: ['Chat not found'] }, status: :not_found
    else
      @chat = chat
    end
  end
end
