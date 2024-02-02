# This file generated by `mix dagger.gen`. Please DO NOT EDIT.
defmodule Dagger.InterfaceTypeDef do
  @moduledoc "A definition of a custom interface defined in a Module."
  use Dagger.QueryBuilder
  @type t() :: %__MODULE__{}
  defstruct [:selection, :client]

  (
    @doc "The doc string for the interface, if any."
    @spec description(t()) :: {:ok, Dagger.String.t()} | {:error, term()}
    def description(%__MODULE__{} = interface_type_def) do
      selection = select(interface_type_def.selection, "description")
      execute(selection, interface_type_def.client)
    end
  )

  (
    @doc "Functions defined on this interface, if any."
    @spec functions(t()) :: {:ok, [Dagger.Function.t()]} | {:error, term()}
    def functions(%__MODULE__{} = interface_type_def) do
      selection = select(interface_type_def.selection, "functions")
      selection = select(selection, "args description id name returnType withArg withDescription")

      with {:ok, data} <- execute(selection, interface_type_def.client) do
        {:ok,
         data
         |> Enum.map(fn value ->
           elem_selection = Dagger.QueryBuilder.Selection.query()
           elem_selection = select(elem_selection, "loadFunctionFromID")
           elem_selection = arg(elem_selection, "id", value["id"])
           %Dagger.Function{selection: elem_selection, client: interface_type_def.client}
         end)}
      end
    end
  )

  (
    @doc "A unique identifier for this InterfaceTypeDef."
    @spec id(t()) :: {:ok, Dagger.InterfaceTypeDefID.t()} | {:error, term()}
    def id(%__MODULE__{} = interface_type_def) do
      selection = select(interface_type_def.selection, "id")
      execute(selection, interface_type_def.client)
    end
  )

  (
    @doc "The name of the interface."
    @spec name(t()) :: {:ok, Dagger.String.t()} | {:error, term()}
    def name(%__MODULE__{} = interface_type_def) do
      selection = select(interface_type_def.selection, "name")
      execute(selection, interface_type_def.client)
    end
  )

  (
    @doc "If this InterfaceTypeDef is associated with a Module, the name of the module. Unset otherwise."
    @spec source_module_name(t()) :: {:ok, Dagger.String.t()} | {:error, term()}
    def source_module_name(%__MODULE__{} = interface_type_def) do
      selection = select(interface_type_def.selection, "sourceModuleName")
      execute(selection, interface_type_def.client)
    end
  )
end