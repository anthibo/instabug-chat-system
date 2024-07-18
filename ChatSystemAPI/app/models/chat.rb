class Chat < ApplicationRecord
  belongs_to :application, counter_cache: true
  has_many :messages, dependent: :destroy

  validates :number, presence: true, uniqueness: { scope: :application_id }

  before_validation :assign_number

  private

  def assign_number
    if self.number.nil?
      application.with_lock do
        max_number = application.chats.maximum(:number) || 0
        self.number = max_number + 1
      end
    end
  end
end
