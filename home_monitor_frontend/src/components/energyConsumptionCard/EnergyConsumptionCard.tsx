import { PropsWithChildren } from "react";
import { Area, AreaChart, CartesianGrid, ResponsiveContainer, XAxis, YAxis } from "recharts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
import { DateTime } from "luxon";

export type EnergyConsumption = {
  name: string;
  total: number;
};

export type EnergyConsumptionCardProps = {};

export const EnergyConsumptionCard = (_props: PropsWithChildren<EnergyConsumptionCardProps>) => {
  const data = [
    {
      total: 0,
      timestamp: DateTime.fromISO("2023-07-30T18:03:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 0,
      timestamp: DateTime.fromISO("2023-07-30T18:04:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 0,
      timestamp: DateTime.fromISO("2023-07-30T18:05:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 0,
      timestamp: DateTime.fromISO("2023-07-30T18:06:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 0,
      timestamp: DateTime.fromISO("2023-07-30T18:07:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 0,
      timestamp: DateTime.fromISO("2023-07-30T18:08:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 0,
      timestamp: DateTime.fromISO("2023-07-30T18:09:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 0,
      timestamp: DateTime.fromISO("2023-07-30T18:10:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 0,
      timestamp: DateTime.fromISO("2023-07-30T18:11:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 0,
      timestamp: DateTime.fromISO("2023-07-30T18:12:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 5354,
      timestamp: DateTime.fromISO("2023-07-30T18:13:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 5428,
      timestamp: DateTime.fromISO("2023-07-30T18:14:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 5492,
      timestamp: DateTime.fromISO("2023-07-30T18:15:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 5440,
      timestamp: DateTime.fromISO("2023-07-30T18:16:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 5480,
      timestamp: DateTime.fromISO("2023-07-30T18:17:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 5505,
      timestamp: DateTime.fromISO("2023-07-30T18:18:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 5636,
      timestamp: DateTime.fromISO("2023-07-30T18:19:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 5552,
      timestamp: DateTime.fromISO("2023-07-30T18:20:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 5641,
      timestamp: DateTime.fromISO("2023-07-30T18:21:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 5721,
      timestamp: DateTime.fromISO("2023-07-30T18:22:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 5777,
      timestamp: DateTime.fromISO("2023-07-30T18:23:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
    {
      total: 5813,
      timestamp: DateTime.fromISO("2023-07-30T18:24:04.733139", { zone: "utc" }).toLocal().toFormat("HH:mm"),
    },
  ];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Energy consumption</CardTitle>
        <CardDescription>Your energy consumption over the last 3 days (kWh)</CardDescription>
      </CardHeader>
      <CardContent className="grid gap-4">
        <ResponsiveContainer width="100%" height={350}>
          <AreaChart data={data}>
            <CartesianGrid />
            <XAxis dataKey="timestamp" className="text-sm" />
            <YAxis className="text-sm" />
            <Area type="monotone" dataKey="total" stroke="#8884d8" fill="#8884d8" />
          </AreaChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
};
