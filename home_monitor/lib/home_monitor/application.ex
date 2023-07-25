defmodule HomeMonitor.Application do
  use Application

  require Logger

  @impl true
  def start(_type, _args) do
    children =
      [
        HomeMonitor.Hm.HmSup,
        HomeMonitor.Tp.TpSup
      ] ++ children(target())

    opts = [strategy: :one_for_one, name: HomeMonitor.Supervisor]

    Logger.info("Starting HomeMonitor with target: #{target()}")

    Supervisor.start_link(children, opts)
  end

  def children(:host) do
    []
  end

  def children(_target) do
    []
  end

  def target() do
    Application.get_env(:home_monitor, :target)
  end
end
