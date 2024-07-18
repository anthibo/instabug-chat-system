# app/services/application_service.rb
class ApplicationService
  include Pagination

  def initialize( application_model)
    @application_model = application_model
  end

  def create_application(body)
      application = Application.create(body)

      if application.save
        Rails.logger.info("Application created successfully: #{application.inspect}")
        application
      else
        Rails.logger.error("Failed to save application: #{application.errors.full_messages.join(', ')}")
        { errors: application.errors.full_messages }
      end
  end

  def update_application(application, params)
    application.update(params)
  end

  def find_application_by_token(token)
    @application_model.find_by(token: token)
  end

  def all_applications(page: 1, per_page: 10)
    paginate(@application_model, page: page, per_page: per_page)
  end
end
