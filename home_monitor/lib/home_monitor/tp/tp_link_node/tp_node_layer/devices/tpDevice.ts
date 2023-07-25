import { Device } from "tplink-smarthome-api";

export abstract class TpDevice {
	powerOn(device: Device): Promise<boolean> {
		throw new Error("Method not implemented.");
	}

	powerOff(device: Device): Promise<boolean> {
		throw new Error("Method not implemented.");
	}
}
