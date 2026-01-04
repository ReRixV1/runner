# runner

`runner` is a small CLI tool to manage background services without tying up your terminal.
Start, stop, list, and inspect running commands in a clean way.

## Installation

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/ReRixV1/runner/refs/heads/main/install.sh)"
```

(or clone & build manually)

Either way you need to have Go and git installed!

## Usage

```bash
runner <command> [service]
```

## Commands

### `run` (aliases: `r`, `start`)

Start a service in the background.

```bash
runner run sketchybar
```

---

### `list` (aliases: `ls`, `l`)

List all running services.

```bash
runner list
```

---

### `stop` (aliases: `s`, `kill`)

Stop a running service.

```bash
runner stop sketchybar
```

---

### `log`

Print the full output of a service and exit.

```bash
runner log sketchybar
```

This does **not** stay open or update live.

---

### `view` (aliases: `v`)

Live view of a service’s output.

```bash
runner view sketchybar
```

This stays open and updates in real-time until you exit with **Ctrl + C**.

---

## `view` vs `log`

* **`view`**
  Live output, keeps running, exits manually (Ctrl + C)

    * At first only shows last 20 lines!

* **`log`**
  One-time dump of all output, exits immediately

---

Small tool, does one thing, stays out of your way.
