defmodule HomeMonitor.Mqtt.MqttProc do
  use GenServer

  require Logger

  import HomeMonitor.Helpers.DefDelegate

  alias HomeMonitor.Helpers.DefDelegate, as: Delegate

  @hal_system Application.compile_env!(:home_monitor, :hal_system)

  def_get_impl(:hal_system, impl: @hal_system)

  def start_link([]) do
    GenServer.start_link(__MODULE__, [], name: __MODULE__)
  end

  def init([]) do
    emqtt_opts =
      Application.fetch_env!(:home_monitor, :emqtt)
      |> Keyword.put(:clientid, hal_system().get_board_id())

    Process.sleep(5_000)

    Logger.info("EMQTT: Starting #{inspect(emqtt_opts)}")

    {:ok, pid} = :emqtt.start_link(emqtt_opts)

    state = %{
      pid: pid,
      emqtt_opts: emqtt_opts
    }

    {:ok, state, {:continue, :start_emqtt}}
  end

  def handle_continue(:start_emqtt, %{pid: pid, emqtt_opts: emqtt_opts} = state) do
    case :emqtt.connect(pid) do
      {:ok, _prop} ->
        Logger.info("EMQTT: Connected")

      {:error, reason} ->
        Logger.error("EMQTT: Failed to connect: #{inspect(reason)}")
    end

    {:ok, _, _} = :emqtt.subscribe(pid, {"commands/#{emqtt_opts[:clientid]}/+", 1})

    {:noreply, state}
  end

  def publish(topic, payload) do
    GenServer.cast(__MODULE__, {:publish, topic, payload})
  end

  # Handle receiving data
  def handle_info({:publish, packet}, state) do
    with {:ok, %{"action" => action, "device_ip" => device_ip, "device_type" => device_type}} <-
           JSON.decode(packet.payload) do
      hal_system().send_command(action, Delegate.convert_ip(device_ip), device_type)
    else
      {:error, reason} ->
        Logger.error("MqttProc: Failed to decode payload: #{inspect(reason)}")
    end

    {:noreply, state}
  end

  # Handle publishing data
  def handle_cast({:publish, topic, payload}, state) do
    payload =
      Map.put(
        payload,
        "timestamp",
        DateTime.utc_now()
        |> DateTime.to_iso8601()
      )

    payload = JSON.encode!(payload)

    case :emqtt.publish(state.pid, "reports/#{state.emqtt_opts[:clientid]}/#{topic}", payload) do
      :ok ->
        Logger.info("EMQTT: Published")

      _ ->
        Logger.error("EMQTT: Failed to publish")
    end

    {:noreply, state}
  end
end
