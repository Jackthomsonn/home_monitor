import { PropsWithChildren } from "react";
import { Area, AreaChart, CartesianGrid, ResponsiveContainer, XAxis, YAxis } from "recharts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";

export type EnergyConsumption = {
  name: string;
  total: number;
};

export type EnergyConsumptionCardProps = {
  data: EnergyConsumption[];
};

export const EnergyConsumptionCard = ({ data }: PropsWithChildren<EnergyConsumptionCardProps>) => (
  <Card>
    <CardHeader>
      <CardTitle>Energy consumption</CardTitle>
      <CardDescription>Your energy consumption over the last 3 days (kWh)</CardDescription>
    </CardHeader>
    <CardContent className="grid gap-4">
      <ResponsiveContainer width="100%" height={350}>
        <AreaChart data={data}>
          <CartesianGrid />
          <XAxis dataKey="name" className="text-sm" />
          <YAxis className="text-sm" />
          <Area type="monotone" dataKey="total" stroke="#8884d8" fill="#8884d8" />
        </AreaChart>
      </ResponsiveContainer>
    </CardContent>
  </Card>
);
