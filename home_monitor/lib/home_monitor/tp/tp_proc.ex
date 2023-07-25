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

    case start_discovery() do
      {:ok, _} ->
        {:ok, []}

      {:error, reason} ->
        {:error, reason}
    end
  end

  def start_discovery() do
    Logger.info("Starting discovery")

    case NodeJS.call({"index", :startDiscovery}, []) do
      {:ok, devices} ->
        Logger.info("Discovered devices: #{inspect(devices)}")

        {:ok, devices}

      {:error, reason} ->
        Logger.info("Failed to start discovery: #{inspect(reason)}")

        {:error, reason}
    end
  end

  def turn_on(device_id) do
    task =
      Task.async(fn ->
        case NodeJS.call({"index", :turnOn}, [device_id]) do
          {:ok, _} ->
            Logger.info("Turned on device")
            {:ok, :device_turned_on}

          {:error, reason} ->
            Logger.info("Failed to turn on device: #{inspect(reason)}")

            {:error, reason}
        end
      end)

    Task.await(task)
  end

  def turn_off(device_id) do
    task =
      Task.async(fn ->
        case NodeJS.call({"index", :turnOff}, [device_id]) do
          {:ok, _} ->
            Logger.info("Turned off device")
            {:ok, :device_turned_off}

          {:error, reason} ->
            Logger.info("Failed to turn off device: #{inspect(reason)}")

            {:error, reason}
        end
      end)

    Task.await(task)
  end
end
