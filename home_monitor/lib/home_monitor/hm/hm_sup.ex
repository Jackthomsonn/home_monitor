defmodule HomeMonitor.Hm.HmSup do
  use Supervisor

  require Logger

  def start_link(_opts) do
    Logger.info("Starting HomeMonitor.Hm.HmSup")
    Supervisor.start_link(__MODULE__, [], name: __MODULE__)
  end

  def init(_opts) do
    children = [
      %{
        id: HomeMonitor.Hm.HmProc,
        start: {HomeMonitor.Hm.HmProc, :start_link, [[]]},
        type: :worker,
        restart: :permanent
      }
    ]

    Supervisor.init(children, strategy: :one_for_one)
  end
end
