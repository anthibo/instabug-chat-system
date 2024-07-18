class CreateMessages < ActiveRecord::Migration[6.0]
  def change
    create_table :messages do |t|
      t.references :chat, null: false, foreign_key: true
      t.text :body, null: false
      t.integer :number, null: false
      
      t.timestamps
    end

    add_index :messages, [:chat_id, :number], unique: true
  end
end
