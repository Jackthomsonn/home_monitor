import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import z from "zod";
import "./App.css";
import { Device, DeviceCard } from "./components/deviceCard/DeviceCard";
import { EnergyConsumptionCard } from "./components/energyConsumptionCard/EnergyConsumptionCard";

function App() {
	const formSchema = z.object({
		device_name: z.string(),
	});

	useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			device_name: "",
		},
	});

	const devices: Device[] = [
		{
			name: "Living room light",
			description: "Lightbulb ",
			status: false,
		},
		{
			name: "Bedroom light",
			description: "Lightbulb ",
			status: true,
		},
		{
			name: "Kitchen light",
			description: "Lightbulb ",
			status: true,
		},
	];

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
