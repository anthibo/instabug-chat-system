class MessageCountWorker
  include Sidekiq::Worker

  def perform(arg1, arg2)
    # Your job logic here
    puts "Doing hard work"
  rescue => e
    Sidekiq.logger.error("Failed job with exception: #{e.message}")
  end
end
