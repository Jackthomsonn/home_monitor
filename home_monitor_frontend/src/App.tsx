import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2, LogOut } from "lucide-react";
import { Magic } from "magic-sdk";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import z from "zod";
import "./App.css";
import { DeviceCard } from "./components/deviceCard/DeviceCard";
import { EnergyConsumptionCard } from "./components/energyConsumptionCard/EnergyConsumptionCard";
import { Button } from "./components/ui/button";
import { Input } from "./components/ui/input";

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [email, setEmail] = useState("");

  const formSchema = z.object({
    device_name: z.string(),
  });

  useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      device_name: "",
    },
  });

  console.log(import.meta.env);

  const magic = new Magic(import.meta.env.VITE_MAGIC_API_KEY, {
    network: "mainnet",
  });

  const login = async () => {
    const success = await magic.auth.loginWithMagicLink({ email });
    if (success) setIsLoggedIn(true);
  };

  useEffect(() => {
    const checkLogin = async () => {
      setIsLoading(true);
      const isLoggedIn = await magic.user.isLoggedIn();
      setIsLoggedIn(isLoggedIn);
      setIsLoading(false);
    };
    checkLogin();
  }, []);

  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center h-screen">
        <Loader2 className="animate-spin w-8 h-8" />
        <p className="font-bold">Loading...</p>
      </div>
    );
  }

  return (
    <>
      {isLoggedIn ? (
        <>
          <div className="pb-4 pt-4 bg-white border-gray-100 border-b-2 flex justify-between items-center">
            {/* rome-ignore lint/a11y/noSvgWithoutTitle: <explanation> */}
            <svg
              className="ml-4"
              width="52"
              height="52"
              viewBox="0 0 113 113"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <rect width="113" height="113" rx="20" fill="#33C676" />
              <path
                d="M50.3594 46H47.9219V54.2031H38.1719V46H35.6875V65.3828H38.1719V56.6641H47.9219V65.3828H50.3594V46ZM45.4844 43.8906C47.0625 43.8906 48.4922 44.125 49.7734 44.5938C51.0547 45.0625 52.1406 45.7031 53.0312 46.5156C53.9375 47.3281 54.6328 48.2812 55.1172 49.375C55.6016 50.4531 55.8438 51.6016 55.8438 52.8203V67.0703C55.8438 68.2891 55.6016 69.4453 55.1172 70.5391C54.6328 71.6172 53.9375 72.5625 53.0312 73.375C52.1406 74.1875 51.0547 74.8281 49.7734 75.2969C48.4922 75.7656 47.0625 76 45.4844 76H38.2656C36.6875 76 35.2578 75.7656 33.9766 75.2969C32.6953 74.8281 31.6016 74.1875 30.6953 73.375C29.8047 72.5625 29.1172 71.6172 28.6328 70.5391C28.1484 69.4453 27.9062 68.2891 27.9062 67.0703V52.8203C27.9062 51.6016 28.1484 50.4531 28.6328 49.375C29.1172 48.2812 29.8047 47.3281 30.6953 46.5156C31.6016 45.7031 32.6953 45.0625 33.9766 44.5938C35.2578 44.125 36.6875 43.8906 38.2656 43.8906H45.4844ZM73.2344 43.8906C74.8125 43.8906 76.2422 44.125 77.5234 44.5938C78.8047 45.0625 79.8906 45.7031 80.7812 46.5156C81.6875 47.3281 82.3828 48.2812 82.8672 49.375C83.3516 50.4531 83.5938 51.6016 83.5938 52.8203V67.0703C83.5938 68.2891 83.3516 69.4453 82.8672 70.5391C82.3828 71.6172 81.6875 72.5625 80.7812 73.375C79.8906 74.1875 78.8047 74.8281 77.5234 75.2969C76.2422 75.7656 74.8125 76 73.2344 76H66.0156C64.4375 76 63.0078 75.7656 61.7266 75.2969C60.4453 74.8281 59.3516 74.1875 58.4453 73.375C57.5547 72.5625 56.8672 71.6172 56.3828 70.5391C55.8984 69.4453 55.6562 68.2891 55.6562 67.0703V52.8203C55.6562 51.6016 55.8984 50.4531 56.3828 49.375C56.8672 48.2812 57.5547 47.3281 58.4453 46.5156C59.3516 45.7031 60.4453 45.0625 61.7266 44.5938C63.0078 44.125 64.4375 43.8906 66.0156 43.8906H73.2344ZM69.5781 46.375L68.2891 46H65.2188L63.25 46.3281L61.75 47.8281V65.3828H64.2578V48.6953C64.4766 48.6484 64.6719 48.6172 64.8438 48.6016C65.0156 48.5859 65.2031 48.5469 65.4062 48.4844H68.0781C68.25 48.5469 68.4141 48.5859 68.5703 48.6016C68.7422 48.6172 68.9141 48.6484 69.0859 48.6953V58.6328H71.5938V48.6953L72.6016 48.4844H75.2969C75.4844 48.5469 75.6719 48.5859 75.8594 48.6016C76.0469 48.6172 76.2422 48.6484 76.4453 48.6953V65.3828H78.9297V47.8281L77.4531 46.3281L75.4844 46H72.4141L71.1484 46.375L70.3516 47.9922L69.5781 46.375Z"
                fill="white"
              />
            </svg>
            <Button className="mr-4">
              <LogOut
                onClick={async () => {
                  await magic.user.logout();
                  setIsLoggedIn(false);
                }}
              />
              <span className="ml-2">Logout</span>
            </Button>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 m-4">
            <DeviceCard />
            <EnergyConsumptionCard />
          </div>
        </>
      ) : (
        <div className="flex flex-col items-center justify-center h-screen">
          <div className="mb-4">
            {/* rome-ignore lint/a11y/noSvgWithoutTitle: Suppress SVG */}
            <svg width="96" height="96" viewBox="0 0 113 113" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect width="113" height="113" rx="20" fill="#33C676" />
              <path
                d="M50.3594 46H47.9219V54.2031H38.1719V46H35.6875V65.3828H38.1719V56.6641H47.9219V65.3828H50.3594V46ZM45.4844 43.8906C47.0625 43.8906 48.4922 44.125 49.7734 44.5938C51.0547 45.0625 52.1406 45.7031 53.0312 46.5156C53.9375 47.3281 54.6328 48.2812 55.1172 49.375C55.6016 50.4531 55.8438 51.6016 55.8438 52.8203V67.0703C55.8438 68.2891 55.6016 69.4453 55.1172 70.5391C54.6328 71.6172 53.9375 72.5625 53.0312 73.375C52.1406 74.1875 51.0547 74.8281 49.7734 75.2969C48.4922 75.7656 47.0625 76 45.4844 76H38.2656C36.6875 76 35.2578 75.7656 33.9766 75.2969C32.6953 74.8281 31.6016 74.1875 30.6953 73.375C29.8047 72.5625 29.1172 71.6172 28.6328 70.5391C28.1484 69.4453 27.9062 68.2891 27.9062 67.0703V52.8203C27.9062 51.6016 28.1484 50.4531 28.6328 49.375C29.1172 48.2812 29.8047 47.3281 30.6953 46.5156C31.6016 45.7031 32.6953 45.0625 33.9766 44.5938C35.2578 44.125 36.6875 43.8906 38.2656 43.8906H45.4844ZM73.2344 43.8906C74.8125 43.8906 76.2422 44.125 77.5234 44.5938C78.8047 45.0625 79.8906 45.7031 80.7812 46.5156C81.6875 47.3281 82.3828 48.2812 82.8672 49.375C83.3516 50.4531 83.5938 51.6016 83.5938 52.8203V67.0703C83.5938 68.2891 83.3516 69.4453 82.8672 70.5391C82.3828 71.6172 81.6875 72.5625 80.7812 73.375C79.8906 74.1875 78.8047 74.8281 77.5234 75.2969C76.2422 75.7656 74.8125 76 73.2344 76H66.0156C64.4375 76 63.0078 75.7656 61.7266 75.2969C60.4453 74.8281 59.3516 74.1875 58.4453 73.375C57.5547 72.5625 56.8672 71.6172 56.3828 70.5391C55.8984 69.4453 55.6562 68.2891 55.6562 67.0703V52.8203C55.6562 51.6016 55.8984 50.4531 56.3828 49.375C56.8672 48.2812 57.5547 47.3281 58.4453 46.5156C59.3516 45.7031 60.4453 45.0625 61.7266 44.5938C63.0078 44.125 64.4375 43.8906 66.0156 43.8906H73.2344ZM69.5781 46.375L68.2891 46H65.2188L63.25 46.3281L61.75 47.8281V65.3828H64.2578V48.6953C64.4766 48.6484 64.6719 48.6172 64.8438 48.6016C65.0156 48.5859 65.2031 48.5469 65.4062 48.4844H68.0781C68.25 48.5469 68.4141 48.5859 68.5703 48.6016C68.7422 48.6172 68.9141 48.6484 69.0859 48.6953V58.6328H71.5938V48.6953L72.6016 48.4844H75.2969C75.4844 48.5469 75.6719 48.5859 75.8594 48.6016C76.0469 48.6172 76.2422 48.6484 76.4453 48.6953V65.3828H78.9297V47.8281L77.4531 46.3281L75.4844 46H72.4141L71.1484 46.375L70.3516 47.9922L69.5781 46.375Z"
                fill="white"
              />
            </svg>
          </div>
          <span className="mb-4 font-bold text-lg">Home Monitor</span>
          <Input className="w-1/2 mb-4" placeholder="Enter your email" onChange={(e) => setEmail(e.target.value)} />
          <Button className="w-1/2" onClick={login}>
            Login
          </Button>
        </div>
      )}
    </>
  );
}

export default App;
