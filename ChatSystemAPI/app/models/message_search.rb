class MessageSearch
  include Elasticsearch::Model
  include Elasticsearch::Model::Callbacks

  index_name "message_body_search_#{Rails.env}"

  settings index: { number_of_shards: 1 } do
    mappings dynamic: 'false' do
      indexes :application_token, type: 'text'
      indexes :chat_number, type: 'integer'
      indexes :message_number, type: 'integer'
      indexes :body, type: 'text'
    end
  end

  def self.index_message(message_data)
    self.__elasticsearch__.client.index(
      index: index_name,
      body: {
        application_token: message_data['application_token'],
        message_number: message_data['message_number'],
        chat_number: message_data['chat_number'],
        body: message_data['body']
      }
    )
  end
end

# Ensure index is created
MessageSearch.__elasticsearch__.create_index!
