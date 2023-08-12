import { PropsWithChildren } from "react";
import { Area, AreaChart, CartesianGrid, ResponsiveContainer, XAxis, YAxis } from "recharts";
import useSWR from "swr";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";

export type EnergyConsumption = {
  power_mw: number;
  timestamp: number;
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

export const EnergyConsumptionCard = (_props: PropsWithChildren<EnergyConsumptionCardProps>) => {
  const { data: energyConsumption } = useSWR<EnergyConsumption[]>("energy_consumption", getEnergyConsumption, {
    refreshInterval: 60_000,
  });

  return (
    <Card>
      <CardHeader>
        <CardTitle>Energy consumption</CardTitle>
        <CardDescription>Your energy consumption over the last 3 days (kWh)</CardDescription>
      </CardHeader>
      <CardContent className="grid gap-4">
        <ResponsiveContainer width="100%" height={350}>
          <AreaChart data={energyConsumption}>
            <CartesianGrid />
            <XAxis dataKey="timestamp" className="text-sm" />
            <YAxis className="text-sm" />
            <Area type="monotone" dataKey="power_mw" stroke="#8884d8" fill="#8884d8" />
          </AreaChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
};
