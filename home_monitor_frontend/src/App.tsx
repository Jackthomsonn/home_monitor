import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import z from "zod";
import "./App.css";
import { Device, DeviceCard } from "./components/deviceCard/DeviceCard";
import { EnergyConsumptionCard } from "./components/energyConsumptionCard/EnergyConsumptionCard";
import { useEffect, useState } from "react";

function App() {
  const formSchema = z.object({
    device_name: z.string(),
  });

  const [devices, setDevices] = useState<Device[]>([]);

  useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      device_name: "",
    },
  });

  useEffect(() => {
    const getDevices = async () => {
      const response = await fetch("http://localhost:8080/getDevices", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          api_key: "0267dddec74ba4a3819ab89342feb108507cb8a67a3e8dc99c992c1058cec74d",
        },
      });
      const data = await response.json();
      setDevices(data);
    };
    getDevices();
  }, []);

  const data = [
    {
      name: "Mon",
      total: Math.floor(Math.random() * 5000) + 1000,
    },
    {
      name: "Tues",
      total: Math.floor(Math.random() * 5000) + 1000,
    },
    {
      name: "Wed",
      total: Math.floor(Math.random() * 5000) + 1000,
    },
  ];

  return (
    <>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 m-4">
        <DeviceCard devices={devices} />
        <EnergyConsumptionCard data={data} />
      </div>
    </>
  );
}

export default App;
