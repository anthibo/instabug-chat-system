class Message < ApplicationRecord
  belongs_to :chat, counter_cache: true

  validates :number, presence: true, uniqueness: { scope: :chat_id }
  validates :body, presence: true

  before_create :assign_number

  private

  def assign_number
    self.number = Message.where(chat_id: chat_id).maximum(:number).to_i + 1
  end
end
