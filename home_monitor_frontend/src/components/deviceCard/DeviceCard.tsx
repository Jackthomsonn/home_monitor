import { Loader2, Power } from "lucide-react";
import { PropsWithChildren } from "react";
import useSWR from "swr";
import useSWRMutation from "swr/mutation";
import { Button } from "../ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../ui/card";

export type Device = {
  A: number;
  K: string;
  ip: string[];
  alias: string;
  feature: string;
  on_time: number;
  device_id: string;
  relay_state: number;
  client_id: string;
  device_type: string;
} & { action: string };

export type DeviceCardProps = {};

const convertDeviceType = (device_type: string) => {
  if (device_type === "IOT.SMARTPLUGSWITCH") return "plug";
};

async function sendCommand(_key: string, { arg }: { arg: Device }) {
  await fetch("https://europe-west1-home-monitor-373013.cloudfunctions.net/SendCommand", {
    method: "POST",
    body: JSON.stringify({
      action: arg.relay_state === 1 ? "turn_off" : "turn_on",
      device_ip: arg.ip.join(", "),
      device_id: arg.client_id,
      device_type: convertDeviceType(arg.device_type),
    }),
    headers: {
      api_key: import.meta.env.VITE_API_KEY,
    },
  });
}

async function discoverDevices(_key: string, { arg }: { arg: Device }) {
  await fetch("https://europe-west1-home-monitor-373013.cloudfunctions.net/SendCommand", {
    method: "POST",
    body: JSON.stringify({
      action: "discover",
      device_id: arg.client_id,
    }),
    headers: {
      api_key: import.meta.env.VITE_API_KEY,
    },
  });
}

const getDevices = async () => {
  const data = await fetch("https://europe-west1-home-monitor-373013.cloudfunctions.net/GetDevices", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      api_key: import.meta.env.VITE_API_KEY,
    },
  });

  return data.json();
};

export const DeviceCard = (_props: PropsWithChildren<DeviceCardProps>) => {
  const { data: devices } = useSWR<Device[]>("devices", getDevices);

  const { trigger: sendCommandTrigger, isMutating: sendCommandIsMutating } = useSWRMutation("devices", sendCommand, {
    optimisticData: (arg: Device[]) => {
      return devices?.map((device, index) => {
        if (device.client_id === arg[index].client_id) {
          return {
            ...device,
            relay_state: device.relay_state === 1 ? 0 : 1,
          };
        }
      });
    },
    revalidate: false,
  });

  const { trigger: discoverDevicesTrigger, isMutating: discoverDevicesIsMutating } = useSWRMutation(
    "devices",
    discoverDevices,
    { revalidate: false },
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle>Devices</CardTitle>
        <CardDescription>Control devices in your home</CardDescription>
      </CardHeader>
      <CardContent className="grid gap-4">
        {devices?.map((device) => {
          const className = device.relay_state ? "text-green-500" : "text-gray-500";
          return (
            <div key={device.device_id} className=" flex items-center space-x-4 rounded-md border p-4">
              <Power className={className} />
              <div className="flex-1 space-y-1">
                <p className="text-sm font-medium leading-none">{device.alias}</p>
                <p className="text-sm text-muted-foreground">{device.feature}</p>
              </div>
              <Button onClick={() => sendCommandTrigger(device)}>
                {sendCommandIsMutating && <Loader2 className="animate-spin mr-2" />}
                {device.relay_state === 1 ? "Turn off" : "Turn on"}
              </Button>
            </div>
          );
        })}
      </CardContent>
      <CardFooter className="w-full">
        {devices && devices?.length > 0 && (
          <Button
            className="w-full"
            onClick={async () => {
              await discoverDevicesTrigger({ client_id: devices[0].client_id, action: "discover" } as Device);
              window.location.reload();
            }}
          >
            {discoverDevicesIsMutating && <Loader2 className="animate-spin mr-2" />}
            Refresh list
          </Button>
        )}
      </CardFooter>
    </Card>
  );
};
