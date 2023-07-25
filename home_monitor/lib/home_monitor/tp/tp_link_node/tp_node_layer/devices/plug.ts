import { Plug } from "tplink-smarthome-api";
import { TpDevice } from "./tpDevice";

export class DevicePlug implements TpDevice {
	powerOn(device: Plug) {
		return device.setPowerState(true);
	}

	powerOff(device: Plug) {
		return device.setPowerState(false);
	}
}
