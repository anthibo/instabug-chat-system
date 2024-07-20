class MessageService
  def search_messages(options = {})
    search = options[:search]
    chat_number = options[:chat_number]
    application_token = options[:application_token]
    page = options.fetch(:page, 1)
    per_page = options.fetch(:per_page, 10)

    query = {
      bool: {
        must: search.present? ? [
          { match_phrase: { body: search } }
        ] : [],
        filter: [
          { term: { chat_number: chat_number } },
          { term: { application_token: application_token } }
        ]
      }
    }

    search_definition = {
      from: (page - 1) * per_page,
      size: per_page,
      query: query
    }

    messages = MessageSearch.__elasticsearch__.search(search_definition).results

    messages.map { |message| message['_source'] }
  end

end