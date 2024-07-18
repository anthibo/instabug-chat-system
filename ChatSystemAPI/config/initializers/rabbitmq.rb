require 'bunny'

class RabbitMQ
  def self.connection
    @connection ||= Bunny.new(
                              host: 'localhost',
                              port: '5672',
                              username: 'user',
                              password: 'bitnami',
                              vhost: '/',
                              connection_timeout: 5,
                              heartbeat: 30,
                              automatically_recover: true)
    @connection.start
  rescue Bunny::TCPConnectionFailedForAllHosts => e
    Rails.logger.error("Failed to connect to RabbitMQ: #{e.message}")
    sleep(5)
    retry
  end

  def self.channel
    @channel ||= connection.create_channel
  end

  def self.close_channel
    Thread.current[:rabbitmq_channel].close if Thread.current[:rabbitmq_channel]
    @channels.delete(Thread.current.object_id)
    Thread.current[:rabbitmq_channel] = nil
  end

  def self.queue(name)
    channel.queue(name, durable: true)
  end

  def self.exchange(name, type = :fanout)
    channel.exchange(name, type: type, durable: true)
  end
end
