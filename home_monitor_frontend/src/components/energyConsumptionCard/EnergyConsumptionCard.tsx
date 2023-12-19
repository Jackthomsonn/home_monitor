import { InfoIcon } from "lucide-react";
import { PropsWithChildren } from "react";
import useSWR from "swr";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";

export type EnergyConsumption = {
  power_mw: number;
  timestamp: number;
  alias: string;
};

type HomeTotals = {
  carbonTotal: number;
  consumptionTotal: number;
};

export type EnergyConsumptionCardProps = {};

const getEnergyConsumption = async () => {
  const data = await fetch("https://europe-west1-home-monitor-373013.cloudfunctions.net/GetEnergyConsumption", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  });

  return data.json();
};

const getHomeTotals = async () => {
  const data = await fetch("https://europe-west1-home-monitor-373013.cloudfunctions.net/GetTotalsForHome", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  });

  return data.json();
};

export const EnergyConsumptionCard = (_props: PropsWithChildren<EnergyConsumptionCardProps>) => {
  const { data: energyConsumption } = useSWR<EnergyConsumption[]>("energy_consumption", getEnergyConsumption, {
    refreshInterval: 60_000,
  });

  const { data: homeTotals } = useSWR<HomeTotals>("home_totals", getHomeTotals, {
    refreshInterval: 60_000,
  });

  // rome-ignore lint/suspicious/noExplicitAny: this is fine to do here
  const totals = energyConsumption?.reduce((a: any, b: any) => {
    const key = b.alias;

    if (a[key]) {
      a[key] += b.power_mw;
    } else {
      a[key] = b.power_mw;
    }

    return a;
  }, []);

  console.log(totals);

  return (
    <>
      <Card>
        <CardHeader>
          <CardTitle>Energy consumption</CardTitle>
          <CardDescription>Your energy consumption over the last 1 hour (kWh)</CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          {Object.keys(totals ?? []).map((key) => {
            return (
              <div key={key} className="bg-violet-50 p-4 rounded-lg flex items-start">
                <InfoIcon />
                <p className="pl-2 text-sm">
                  Your {key} has consumed <span className="font-bold text-green-500">{totals[key] / 100000}</span> kWh
                  of power in the last 1 hour
                </p>
              </div>
            );
          })}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Totals for the previous day</CardTitle>
          <CardDescription>
            Below is your total emitted carbon and total consumption for the previous day
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          <p className="text-sm">
            In the last 1 day, you have outputed{" "}
            <span className="font-bold text-green-500 bg-green-50 p-1 ml-1 mr-1 rounded-lg">
              {homeTotals?.carbonTotal} g/CO2
            </span>{" "}
            of carbon and consumed
            <span className="font-bold text-green-500 bg-green-50 p-1 ml-1 mr-1 rounded-lg">
              {" "}
              {homeTotals?.consumptionTotal} kWh of energy
            </span>
          </p>
        </CardContent>
      </Card>
    </>
  );
};
