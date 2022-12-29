defmodule HomeMonitor.Target.Host do
  def get_temperature() do
    10.0 + 2.0 * :rand.normal()
  end

  def get_board_id() do
    "host"
  end
end
