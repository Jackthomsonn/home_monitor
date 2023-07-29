defmodule HomeMonitor.Target.Host do
  def get_temperature() do
    10.0 + 2.0 * :rand.normal()
  end

  def get_board_id() do
    "host"
  end

  def send_command(action, ip, device_type) do
    HomeMonitor.Tp.TpProc.send_command(action, ip, device_type)
  end

  def send_command(action) do
    HomeMonitor.Tp.TpProc.send_command(action)
  end
end
