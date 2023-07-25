import { Device, Plug, Client, Bulb } from "tplink-smarthome-api";
import { DevicePlug } from "./devices/plug";
import { TpDevice } from "./devices/tpDevice";

const client = new Client();

const deviceFactory: Record<string, TpDevice> = {
	plug: new DevicePlug(),
};

const startDiscovery = () => {
	return new Promise((resolve) => {
		client.startDiscovery().on("device-new", async (device: Device) => {
			const info = (await device.getInfo()) as unknown as Bulb | Plug;

			resolve({
				deviceId: info.sysInfo.deviceId,
				deviceName: info.sysInfo.alias,
				emeter: info.emeter,
			});
		});
	});
};

const turnOn = async (deviceId: string) => {
	return new Promise((resolve) => {
		client.startDiscovery().on("device-new", async (device: Device) => {
			if (device.deviceId === deviceId) {
				await deviceFactory[device.deviceType].powerOn(device);

				resolve({});
			}
		});
	});
};

const turnOff = async (deviceId: string) => {
	return new Promise((resolve) => {
		client.startDiscovery().on("device-new", async (device: Device) => {
			if (device.deviceId === deviceId) {
				await deviceFactory[device.deviceType].powerOff(device);

				resolve({});
			}
		});
	});
};

module.exports = {
	startDiscovery,
	turnOn,
	turnOff,
};
