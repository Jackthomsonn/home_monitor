defmodule HomeMonitor.Target.Rpi4 do
  def get_temperature() do
    case NervesDHT.read(:dht22, 2) do
      {:ok, _, temp} ->
        temp

      {:error, :timeout} ->
        {:error, "the DHT22 sensor timed out while reading the temperature and humidity values"}
    end
  end

  def get_board_id() do
    case System.cmd("/usr/bin/boardid", ["-b", "uboot_env", "-u", "serial_number"])
         |> elem(0)
         |> String.trim() do
      "" -> "unknown"
      id -> id
    end
  end

  def send_command(action, ip, device_type) do
    HomeMonitor.Tp.TpProc.send_command(action, ip, device_type)
  end

  def send_command(action) do
    HomeMonitor.Tp.TpProc.send_command(action)
  end
end
