import Config

config :home_monitor, :target, :host

config :home_monitor,
  hal_system: HomeMonitor.Target.Host
