class ApplicationsController < ApplicationController
  include Pagination

  before_action :find_application_by_token, only: [:show, :update]

  def initialize(application_service = ApplicationService.new(Application))
    @application_service = application_service
  end

  def index
    page = (params[:page].presence || 1).to_i
    per_page = (params[:per_page].presence || 10).to_i
    Rails.logger.info("Page: #{page}, Per page: #{per_page}")
    @applications = @application_service.all_applications(page: page, per_page: per_page)
    render json: { items: @applications }
  end

  def show
    render json: @application
  end

  def create
    Rails.logger.info("Application params: #{application_params.inspect}")
    result = @application_service.create_application(application_params)

    if result.is_a?(Application)
      render json: result, status: :created
    else
      render json: { errors: result[:errors] }, status: :internal_server_error
    end
  end

  def update
    if @application_service.update_application(@application, application_params)
      render json: @application
    else
      render json: @application.errors, status: :internal_server_error
    end
  end

  private

  def find_application_by_token
    token = params[:token]
    @application = @application_service.find_application_by_token(token)
  end

  def application_params
    params.require(:application).permit(:name)
  end
end
