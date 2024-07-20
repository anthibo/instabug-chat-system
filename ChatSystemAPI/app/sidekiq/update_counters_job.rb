class UpdateCountersJob
  include Sidekiq::Job

  def perform(*args)
    update_application_counters
    update_chat_counters
  end

  def update_application_counters
    logger.info("Updating application counters")
    Application.find_in_batches(batch_size: 100) do |applications|
      applications.each do |application|
        logger.info("Updating application counters for application: #{application.id}, count: #{application.chats.count}")
        chats_count = application.chats.count
        ActiveRecord::Base.connection.execute(
          "UPDATE applications SET chats_count = #{chats_count} WHERE id = #{application.id}"
        )
        logger.info("Updated chats_count for application #{application.id} to #{chats_count}")
      end
    end
  end

  def update_chat_counters
    logger.info("Updating chat counters")
    Chat.find_in_batches(batch_size: 100) do |chats|
      chats.each do |chat|
        chat_id = chat.id
        messages_count = chat.messages.count
        logger.info("Chat messages count: #{messages_count} for chat: #{chat_id}")

        ActiveRecord::Base.connection.execute(
          "UPDATE chats SET messages_count = #{messages_count} WHERE id = #{chat_id}"
        )
        logger.info("Chat messages count updated to: #{messages_count} for chat: #{chat_id}")
      end
    end
  end
end
