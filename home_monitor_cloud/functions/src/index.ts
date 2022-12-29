import { HttpFunction } from "@google-cloud/functions-framework";
import { DateTime } from "luxon";
import fetch from "node-fetch";

const suppressOpts = {
  suppressMilliseconds: true,
};

type Carbonintensity = {
  index: "very low" | "low" | "moderate" | "high" | "very high";
  forecast: number;
  actual: number;
};

type CarbonintensityData = {
  from: string;
  to: string;
  intensity: Carbonintensity;
};

type CarbonintensityResponse = {
  data: CarbonintensityData[];
};

enum Action {
  TURN_ON = "TURN_ON",
  TURN_OFF = "TURN_OFF",
  MAYBE_TURN_ON = "MAYBE_TURN_ON",
}

export const performCheck: HttpFunction = async (req: any, res: any) => {
  const utc = DateTime.utc().set({ millisecond: 0 });

  const now = utc.toISO(suppressOpts);
  const nowPlus30Minutes = utc.plus({ minutes: 30 }).toISO(suppressOpts);

  // Check the current grid energy (dirty, clean)
  const result = await fetch(
    `https://api.carbonintensity.org.uk/intensity/${now}/${nowPlus30Minutes}`
  );

  const data = (await result.json()) as CarbonintensityResponse;

  const [latestData] = data.data;

  const { intensity } = latestData;

  const efficientIntensities = ["very low", "low"];

  const midEfficientIntensities = ["moderate"];

  if (efficientIntensities.includes(intensity.index)) {
    return res.send({
      action: Action.TURN_ON,
      index: intensity.index,
      forecast: intensity.forecast,
      unit: "gCO2/kWh",
    });
  }

  if (midEfficientIntensities.includes(intensity.index)) {
    return res.send({
      action: Action.MAYBE_TURN_ON,
      index: intensity.index,
      forecast: intensity.forecast,
      unit: "gCO2/kWh",
    });
  }

  return res.send({
    action: Action.TURN_OFF,
  });
};
