defmodule HomeMonitor.Tp.TpProc do
  use GenServer

  require Logger

  def start_link([]) do
    path =
      __ENV__.file
      |> Path.dirname()
      |> Path.join("tp_link_node")
      |> Path.join("tp_node_layer")

    GenServer.start_link(__MODULE__, [path])
  end

  def init([path]) do
    NodeJS.start_link(path: path)

    start_discovery()

    {:ok, []}
  end

  def start_discovery() do
    Logger.info("Starting discovery")

    NodeJS.call({"index", :startDiscovery}, [])
    |> IO.inspect()
  end

  def turn_on(device_id) do
    Task.async(fn ->
      NodeJS.call({"index", :turnOn}, [device_id])
      {:ok, :device_turned_on}
    end)
  end

  def turn_off(device_id) do
    Task.async(fn ->
      NodeJS.call({"index", :turnOff}, [device_id])
      {:ok, :device_turned_off}
    end)
  end
end
