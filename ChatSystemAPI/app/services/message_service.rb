class MessageService
  def search_messages(search, chat_id, page = 1, per_page = 10)
    query = if search.present?
              {
                bool: {
                  must: [
                    { match_phrase: { body: search } }
                  ],
                  filter: [
                    { term: { chat_id: chat_id } }
                  ]
                }
              }
            else
              {
                term: { chat_id: chat_id }
              }
            end

    search_definition = {
      from: (page - 1) * per_page,
      size: per_page,
      query: query
    }

    messages = MessageSearch.__elasticsearch__.search(search_definition).results

    messages.map { |message| message['_source'] }
  end

end