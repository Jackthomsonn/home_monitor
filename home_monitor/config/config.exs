# This file is responsible for configuring your application and its
# dependencies.
#
# This configuration file is loaded before any dependency and is restricted to
# this project.
import Config

# Enable the Nerves integration with Mix
Application.start(:nerves_bootstrap)

config :home_monitor, :interval, 5000

config :nerves, :firmware, rootfs_overlay: "rootfs_overlay"

config :nerves, source_date_epoch: "1671808184"

config :elixir, :time_zone_database, Tzdata.TimeZoneDatabase

if Mix.target() == :host do
  import_config "host.exs"
else
  import_config "target.exs"
end
