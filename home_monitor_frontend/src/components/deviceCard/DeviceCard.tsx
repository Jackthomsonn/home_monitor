import { cn } from "@/lib/utils";
import { Switch } from "@radix-ui/react-switch";
import { Power } from "lucide-react";
import { PropsWithChildren } from "react";
import { AddDeviceDialog } from "../addDeviceDialog/addDeviceDialog";
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "../ui/card";

export type Device = {
	name: string;
	description: string;
	status: boolean;
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
							<p className="text-sm font-medium leading-none">{device.name}</p>
							<p className="text-sm text-muted-foreground">
								{device.description}
							</p>
						</div>
						<Switch checked={device.status} />
					</div>
				);
			})}
		</CardContent>
		<CardFooter className="w-full">
			<AddDeviceDialog />
		</CardFooter>
	</Card>
);
