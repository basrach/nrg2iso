# **nrg2iso**

**nrg2iso** is a simple command-line tool that converts NERO image files (`.nrg`) to ISO format (`.iso`). ISO files are widely supported and offer greater compatibility across operating systems and disk mounting software. This utility allows you to easily convert proprietary NERO image files into standard ISO files.

---

## **Features**
- Converts `.nrg` files to `.iso`
- Fast and efficient
- Lightweight, minimal dependencies
- Cross-platform support (Linux, macOS, Windows)

---

## **Installation**

You can download the latest release from the [Releases](https://github.com/yourusername/nrg2iso/releases) page.

Clone the repository and compile:
```bash
git clone https://github.com/yourusername/nrg2iso.git
cd nrg2iso
go build .
```

---

## **Usage**

Convert a `.nrg` file to `.iso` using the following syntax:

```bash
nrg2iso <input.nrg> <output.iso>
```

### **Example:**
```bash
nrg2iso mydisk.nrg mydisk.iso
```

### **Options:**
- `-h`, `--help`    : Display the help message.

