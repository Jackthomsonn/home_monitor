defmodule HomeMonitor.Hm.HmProc do
  use GenServer
  require Logger

  import HomeMonitor.Helpers.DefDelegate

  @hal_system Application.compile_env!(:home_monitor, :hal_system)

  def_get_impl(:hal_system, impl: @hal_system)

  def start_link([]) do
    GenServer.start_link(__MODULE__, [])
  end

  def get_clientid() do
    hal_system().get_board_id()
  end

  def init([]) do
    interval = Application.fetch_env!(:home_monitor, :interval)

    emqtt_opts = Application.fetch_env!(:home_monitor, :emqtt)

    emqtt_opts = Keyword.put(emqtt_opts, :clientid, get_clientid())

    report_topic = "reports/#{emqtt_opts[:clientid]}/temperature"

    Process.sleep(5000)

    Logger.info("EMQTT: Starting #{inspect(emqtt_opts)}")

    {:ok, pid} = :emqtt.start_link(emqtt_opts)

    st = %{
      interval: interval,
      timer: nil,
      report_topic: report_topic,
      pid: pid
    }

    {:ok, set_timer(st), {:continue, :start_emqtt}}
  end

  def handle_continue(:start_emqtt, %{pid: pid} = st) do
    clientid = get_clientid()

    case :emqtt.connect(pid) do
      {:ok, _prop} ->
        Logger.info("EMQTT: Connected")

      {:error, reason} ->
        Logger.error("EMQTT: Failed to connect: #{inspect(reason)}")
    end

    {:ok, _, _} = :emqtt.subscribe(pid, {"commands/#{clientid}/+", 1})

    {:noreply, st}
  end

  def handle_info(:tick, %{report_topic: topic, pid: pid} = st) do
    report_temperature(pid, topic)

    {:noreply, set_timer(st)}
  end

  def handle_info({:publish, publish}, st) do
    handle_publish(parse_topic(publish), publish, st)
  end

  def handle_info(_, st) do
    {:noreply, st}
  end

  defp handle_publish(["commands", _, "test"], %{payload: payload}, st) do
    case JSON.decode(payload) do
      {:ok, %{"device_id" => device_id, "action" => "turn_on"}} ->
        Logger.info("HmProc: Received turn on command")
        HomeMonitor.Tp.TpProc.turn_on(device_id)

      {:ok, %{"device_id" => device_id, "action" => "turn_off"}} ->
        Logger.info("HmProc: Received turn off command")
        HomeMonitor.Tp.TpProc.turn_off(device_id)

      {:error, reason} ->
        Logger.error("HmProc: Failed to decode test command: #{inspect(reason)}")
    end

    {:noreply, st}
  end

  defp handle_publish(_, _, st) do
    {:noreply, st}
  end

  defp parse_topic(%{topic: topic}) do
    String.split(topic, "/", trim: true)
  end

  defp set_timer(st) do
    if st.timer do
      Process.cancel_timer(st.timer)
    end

    timer = Process.send_after(self(), :tick, st.interval)
    %{st | timer: timer}
  end

  defp report_temperature(pid, topic) do
    temp = hal_system().get_temperature()

    now =
      DateTime.utc_now()
      |> DateTime.to_iso8601()

    case JSON.encode(%{temperature: temp, timestamp: now}) do
      {:ok, payload} ->
        :emqtt.publish(pid, topic, payload)

      {:error, reason} ->
        Logger.error("HmProc: Failed to encode temperature: #{inspect(reason)}")
    end
  end
end
