class MessageSearch
  include Elasticsearch::Model
  include Elasticsearch::Model::Callbacks

  index_name "message_body_search_#{Rails.env}"

  settings index: { number_of_shards: 1 } do
    mappings dynamic: 'false' do
      indexes :id, type: 'integer'
      indexes :chat_id, type: 'integer'
      indexes :body, type: 'text'
    end
  end

  def self.index_message(message_data)
    self.__elasticsearch__.client.index(
      index: index_name,
      id: message_data['id'],
      body: {
        id: message_data['id'],
        chat_id: message_data['chat_id'],
        body: message_data['body']
      }
    )
  end
end

# Ensure index is created
MessageSearch.__elasticsearch__.create_index!
