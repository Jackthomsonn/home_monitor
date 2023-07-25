defmodule HomeMonitor.Target.Host do
  def get_temperature() do
    10.0 + 2.0 * :rand.normal()
  end

  def get_board_id() do
    "host"
  end

  def turn_on(device_id) do
    HomeMonitor.Tp.TpProc.turn_on(device_id)
  end

  def turn_off(device_id) do
    HomeMonitor.Tp.TpProc.turn_off(device_id)
  end
end
