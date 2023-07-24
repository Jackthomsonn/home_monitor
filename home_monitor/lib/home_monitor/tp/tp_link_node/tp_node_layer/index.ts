import { Device, Plug, Client } from "tplink-smarthome-api";
import * as mqtt from "mqtt";
const mqttClient = mqtt.connect("mqtt://35.187.59.21", {
	username: "admin",
	clientId: "host",
});

const client = new Client();

const startDiscovery = () => {
	return new Promise((resolve) => {
		client.startDiscovery().on("device-new", async (device: Device) => {
			if (device.deviceType === "plug") {
				const info = (await device.getInfo()) as unknown as Plug;
				resolve({
					deviceId: info.sysInfo.deviceId,
					deviceName: info.sysInfo.alias,
					emeter: info.emeter,
				});
			}
		});
	});
};

const turnOn = async (deviceId: string) => {
	return new Promise((resolve) => {
		client.startDiscovery().on("device-new", async (device: Device) => {
			if (device.deviceType === "plug") {
				if (device.deviceId === deviceId) {
					(device as Plug).setPowerState(true);
					resolve({});
				}
			}
		});
	});
};

const turnOff = async (deviceId: string) => {
	return new Promise((resolve) => {
		client.startDiscovery().on("device-new", async (device: Device) => {
			if (device.deviceId === deviceId) {
				(device as Plug).setPowerState(false);
				resolve({});
			}
		});
	});
};

const test = () => {
	mqttClient.on("connect", () => {
		console.log("connected");
		mqttClient.publish("commands/b766/test", "Hello mqtt");
	});

	mqttClient.on("message", (topic, message) => {
		console.log(message.toString());
		mqttClient.end();
	});
};

module.exports = {
	startDiscovery,
	turnOn,
	turnOff,
};

test();
