import { Switch } from "@radix-ui/react-switch";
import { Power } from "lucide-react";
import { PropsWithChildren } from "react";
import { Button } from "../ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../ui/card";
import useSWRMutation from "swr/mutation";

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
};

export type DeviceCardProps = {
  devices?: Device[] | undefined;
};

async function sendCommand(url: string, { arg }: { arg: Device }) {
  await fetch(url, {
    method: "POST",
    body: JSON.stringify({
      action: "turn_on",
      device_ip: arg.ip.join(", "),
      device_id: arg.client_id,
      device_type: "plug",
    }),
    headers: {
      api_key: import.meta.env.VITE_API_KEY,
    },
  });
}

export const DeviceCard = ({ devices }: PropsWithChildren<DeviceCardProps>) => {
  const { trigger } = useSWRMutation(
    "https://europe-west1-home-monitor-373013.cloudfunctions.net/SendCommand",
    sendCommand,
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
              <Button onClick={() => trigger(device)}>Toggle</Button>
              <Switch checked={Boolean(device.relay_state)} />
            </div>
          );
        })}
      </CardContent>
      <CardFooter className="w-full">
        <Button className="w-full">Refresh list</Button>
      </CardFooter>
    </Card>
  );
};
