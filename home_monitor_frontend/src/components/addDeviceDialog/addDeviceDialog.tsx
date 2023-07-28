import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { zodResolver } from "@hookform/resolvers/zod";
import { FormProvider, useForm } from "react-hook-form";
import z from "zod";
import { FormField, FormItem, FormLabel, FormControl, FormDescription, FormMessage } from "../ui/form";

export function AddDeviceDialog() {
  const formSchema = z.object({
    device_name: z.string(),
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      device_name: "",
    },
  });

  const onChange = () => {};
  const onSubmit = () => {};

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button className="w-full">Add device</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add device</DialogTitle>
          <DialogDescription>Add a new device to your home</DialogDescription>
        </DialogHeader>
        <div>
          <FormProvider {...form}>
            <form name="add_device_form" onSubmit={form.handleSubmit(onSubmit)}>
              <FormField
                control={form.control}
                name="device_name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Device name</FormLabel>
                    <FormControl>
                      <Input {...field} type="text" onChange={onChange} />
                    </FormControl>
                    <FormDescription>Enter a name for your device</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </form>
          </FormProvider>
        </div>
        <DialogFooter>
          <Button type="submit">Add device</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
