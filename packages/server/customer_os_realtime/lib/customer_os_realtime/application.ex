defmodule CustomerOsRealtime.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  @impl true
  def start(_type, _args) do
    children = [
      CustomerOsRealtimeWeb.Telemetry,
      {DNSCluster,
       query: Application.get_env(:customer_os_realtime, :dns_cluster_query) || :ignore},
      {Phoenix.PubSub, name: CustomerOsRealtime.PubSub},
      # Start a worker by calling: CustomerOsRealtime.Worker.start_link(arg)
      # {CustomerOsRealtime.Worker, arg},
      # Start to serve requests, typically the last entry
      CustomerOsRealtimeWeb.Presence,
      CustomerOsRealtimeWeb.Endpoint,
      CustomerOsRealtime.ColorManager,
      CustomerOsRealtime.DeltaManager,
      CustomerOsRealtime.StoreManager
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: CustomerOsRealtime.Supervisor]
    Supervisor.start_link(children, opts)
  end

  # Tell Phoenix to update the endpoint configuration
  # whenever the application is updated.
  @impl true
  def config_change(changed, _new, removed) do
    CustomerOsRealtimeWeb.Endpoint.config_change(changed, removed)
    :ok
  end
end
