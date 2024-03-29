defmodule HomeMonitor.Helpers.DefDelegate do
  @moduledoc """
  This module defines a helper function to return the implementation
  module for a function.
  If test in mode, the implementation is found via
    Application.get_env(application, env, default)
  and in dev and prod modes the default implementation is returned.
  ## Example
  ```
    defmodule SomeModule do
      use Helpers.DefDelegate
      def_get_impl(:register_impl, impl: RegisterImpl)
      def register(data) do
        register_impl().register(data)
      end
    end
  ```
  (c) 2017 Copyright Component X Software Limited / Antony Pinchbeck
  MIT Licenced
  """
  defmacro def_get_impl(func, opts) do
    name = get_env(opts, :impl)

    quote location: :keep do
      source_line = __ENV__.line

      env = Mix.env()
      validate_options(unquote(opts), source_line)

      if :test == env do
        def unquote(func)() do
          {:ok, app} = :application.get_application(__MODULE__)
          env_name = unquote(func)
          Application.get_env(app, env_name, unquote(name))
        end
      else
        @doc """
        Returns `#{unquote(name)}`
        See `Helpers.DefDelegate` for more information
        """
        def unquote(func)() do
          unquote(name)
        end
      end
    end
  end

  @doc false
  def validate_options(opts, source_line) do
    if nil == Keyword.get(opts, :impl, nil) do
      raise """
      Missing impl: for defdelegate_env on line #{source_line}
      """
    end
  end

  @doc false
  def get_env(opts, env) do
    Keyword.get(opts, env, nil)
  end

  def convert_ip(payload) do
    string_list = String.split(payload, ", ")

    integer_list = Enum.map(string_list, &String.to_integer/1)

    tuple = List.to_tuple(integer_list)
  end
end
