import { Switch } from "@radix-ui/react-switch";
import { Power } from "lucide-react";
import { PropsWithChildren } from "react";
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
};

export type DeviceCardProps = {
  devices: Device[];
};

export const DeviceCard = ({ devices }: PropsWithChildren<DeviceCardProps>) => (
  <Card>
    <CardHeader>
      <CardTitle>Devices</CardTitle>
      <CardDescription>Control devices in your home</CardDescription>
    </CardHeader>
    <CardContent className="grid gap-4">
      {devices.map((device) => {
        return (
          <div className=" flex items-center space-x-4 rounded-md border p-4">
            <Power />
            <div className="flex-1 space-y-1">
              <p className="text-sm font-medium leading-none">{device.alias}</p>
              <p className="text-sm text-muted-foreground">{device.feature}</p>
            </div>
            <Switch checked={true} />
          </div>
        );
      })}
    </CardContent>
    <CardFooter className="w-full">
      <Button className="w-full">Refresh list</Button>
    </CardFooter>
  </Card>
);
