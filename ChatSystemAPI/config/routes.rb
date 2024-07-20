Rails.application.routes.draw do
  resources :applications, param: :token, only: [:index, :show, :create, :update] do
    resources :chats, only: [:create]
  end

  resources :chats, param: :number, only: [:show, :index]

  get 'messages/', to: 'messages#search'

  get "up" => "rails/health#show", as: :rails_health_check
  # Defines the root path route ("/")
  # root "posts#index"
end
