"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DevicePlug = void 0;
var DevicePlug = (function () {
    function DevicePlug() {
    }
    DevicePlug.prototype.powerOn = function (device) {
        return device.setPowerState(true);
    };
    DevicePlug.prototype.powerOff = function (device) {
        return device.setPowerState(false);
    };
    return DevicePlug;
}());
exports.DevicePlug = DevicePlug;
