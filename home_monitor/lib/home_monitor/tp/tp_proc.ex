defmodule HomeMonitor.Tp.TpProc do
  use GenServer

  require Logger

  import HomeMonitor.Helpers.DefDelegate

  @hal_system Application.compile_env!(:home_monitor, :hal_system)

  def_get_impl(:hal_system, impl: @hal_system)

  def start_link([]) do
    GenServer.start_link(__MODULE__, [])
  end

  def init([]) do
    Process.send_after(self(), :monitor_energy_consumption, 60_000)
    {:ok, []}
  end

  def send_command(action, device_ip, device_type) do
    case device_type do
      "plug" ->
        handle_plug_action(action, device_ip)

      _ ->
        Logger.error("TpProc: Unknown device type: #{inspect(device_type)}")
    end
  end

  def handle_plug_action("turn_on", device_ip) do
    TpLink.local_device(device_ip)
    |> TpLink.Type.Plug.set_relay_state(true)
  end

  def handle_plug_action("turn_off", device_ip) do
    TpLink.local_device(device_ip)
    |> TpLink.Type.Plug.set_relay_state(false)
  end

  def handle_plug_action(_, _device_ip) do
    Logger.error("TpProc: Unknown plug action")
  end

  def monitor_energy_consumption() do
    case TpLink.Local.list_devices() do
      {:ok, devices} ->
        devices
        |> Enum.filter(fn device -> Map.get(device.system_info, "feature") == "TIM:ENE" end)
        |> Enum.map(fn device -> monitor_plug(device) end)

      {:error, reason} ->
        Logger.error("TpProc: Failed to list devices: #{inspect(reason)}")
    end

    Process.send_after(self(), :monitor_energy_consumption, 60_000)
  end

  def monitor_plug(device_details) do
    IO.inspect(device_details, label: "Device details")
    device = TpLink.local_device(device_details.ip)

    with {:ok, info} <- TpLink.Type.Plug.get_energy_meter_information(device) do
      packet = %{
        "voltage_mv" => Map.get(info, "voltage_mv"),
        "current_ma" => Map.get(info, "current_ma"),
        "power_mw" => Map.get(info, "power_mw"),
        "total_wh" => Map.get(info, "total_wh"),
        "err_code" => Map.get(info, "err_code"),
        "err_msg" => Map.get(info, "err_msg"),
        "ip" => Tuple.to_list(device_details.ip) |> Enum.join("."),
        "alias" => Map.get(device_details.system_info, "alias")
      }

      HomeMonitor.Mqtt.MqttProc.publish("energy", packet)
    else
      {:error, reason} ->
        Logger.error("TpProc: Failed to get energy data: #{inspect(reason)}")
    end
  end

  def handle_info(:monitor_energy_consumption, state) do
    monitor_energy_consumption()

    {:noreply, state}
  end
end
