#!/usr/bin/env -S deno run -A
import type * as sunbeam from "https://deno.land/x/sunbeam/mod.ts";

const manifest = {
  title: "Tailscale",
  description: "Manage your tailscale devices",
  commands: [
    {
      name: "list-devices",
      title: "Search My Devices",
      mode: "filter",
    },
    {
      name: "ssh-to-device",
      title: "SSH to Device",
      mode: "tty",
      params: [
        {
          name: "ip",
          title: "Device IP",
          type: "string",
        },
      ],
    },
  ],
} as const satisfies sunbeam.Manifest;

if (Deno.args.length == 0) {
  console.log(JSON.stringify(manifest));

  Deno.exit(0);
}

type Device = {
  TailscaleIPs: string[];
  DNSName: string;
  OS: string;
  Online: boolean;
};

const payload: sunbeam.Payload<typeof manifest> = JSON.parse(Deno.args[0]);

if (payload.command == "list-devices") {
  const command = new Deno.Command("tailscale", { args: ["status", "--json"] });
  const { stdout } = await command.output();
  const status = JSON.parse(new TextDecoder().decode(stdout));
  const devices: Device[] = Object.values(status.Peer);
  const items: sunbeam.ListItem[] = devices.map((device) => ({
    title: device.DNSName.split(".")[0],
    subtitle: device.TailscaleIPs[0],
    accessories: [device.OS, device.Online ? "online" : "offline"],
    actions: [
      {
        title: "SSH to Device",
        command: "ssh-to-device",
        params: {
          ip: device.TailscaleIPs[0],
        },
      },
      {
        title: "Copy SSH Command",
        extension: "std",
        command: "copy",
        params: {
          text: `ssh ${device.TailscaleIPs[0]}`,
        },
      },
      {
        title: "Copy IP",
        extension: "std",
        command: "copy",
        params: {
          text: device.TailscaleIPs[0],
        },
      },
    ],
  }));

  const list: sunbeam.List = { items };

  console.log(JSON.stringify(list));
} else if (payload.command == "ssh-to-device") {
  const command = new Deno.Command("ssh", { args: [payload.params.ip] });
  const ps = command.spawn();
  await ps.status;
}
